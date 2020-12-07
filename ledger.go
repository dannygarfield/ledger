package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// transaction represents a double-entry accounting item in the ledger.
type entry struct {
	source      string
	destination string
	happenedAt  time.Time
	amount      int
}

// a bucket describes ownership and accessibility of money
type bucket struct {
	name      string
	asset     bool
	liquidity string
}

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	// define flags for input from the command line
	insertMode := flag.Bool("insert", false, "insert a transaction")
	summaryMode := flag.Bool("summary", false, "get balances of all buckets")
	through := flag.String("through", "", "date through which to summarize")
	source := flag.String("source", "", "bucket from which the amount is taken")
	destination := flag.String("destination", "", "bucket into which the amount is deposited")
	happenedAt := flag.String("happenedAt", "", "date of transaction")
	amount := flag.Int("amount", 0, "amount in cents of the transaction")
	repeat := flag.String("repeat", "", "how often an entry repeats: 'weekly' or 'monthly'")
	assets := flag.Bool("assets", false, "include only money in your posession")
	zeroMode := flag.Bool("zero", false, "find when a bucket zeroes out")
	flag.Parse()

	// open connection to the db
	db, err := sql.Open("sqlite3", "./db.sqlite3")
	if err != nil {
		log.Fatalf("opening database: %v", err)
	}
	defer db.Close()

	if *insertMode && *summaryMode {
		// instruct user to pick only one mode
		log.Printf("only use one of -insert or -summary")
		return
	} else if !*insertMode && !*summaryMode && !*zeroMode {
		// instruct user to pick a mode
		log.Printf("specify one of -insert or -summary or --zero")
		return
	} else if *insertMode && *repeat != "" {
		// insert entry that repeats through 2 years from today
		d, err := parseDate(*happenedAt)
		if err != nil {
			log.Print(err)
			return
		}
		e := entry{
			source:      *source,
			destination: *destination,
			happenedAt:  d,
			amount:      *amount,
		}
		tx, err := db.Begin()
		if err != nil {
			log.Fatalf("beginning sql transaction: %v", err)
		}
		if err := insertRepeating(tx, e, *repeat); err != nil {
			log.Fatalf("inserting a repeating entry: %v", err)
		}
		if err := tx.Commit(); err != nil {
			log.Fatalf("committing sql transaction: %v", err)
		}
	} else if *insertMode {
		// insert a transaction to the db
		d, err := parseDate(*happenedAt)
		if err != nil {
			log.Print(err)
			return
		}
		e := entry{
			source:      *source,
			destination: *destination,
			happenedAt:  d,
			amount:      *amount,
		}
		tx, err := db.Begin()
		if err != nil {
			log.Fatalf("beginning sql transaction: %v", err)
		}
		if err := insert(tx, e); err != nil {
			log.Fatalf("inserting single entry")
		}
		if err := tx.Commit(); err != nil {
			log.Fatalf("committing sql transaction: %v", err)
		}
	} else if *summaryMode && *through != "" && *assets {
		// summarize all assets through a given date
		td, err := parseDate(*through)
		if err != nil {
			log.Print(err)
			return
		}
		tx, err := db.Begin()
		if err != nil {
			log.Fatalf("beginning sql transaction: %v", err)
		}
		sum, err := sumAssets(tx, td)
		if err != nil {
			log.Fatalf("summing assets: %v", err)
		}
		if err := tx.Commit(); err != nil {
			log.Fatalf("committing sql transaction: %v", err)
		}
		log.Printf("All assets as of %v: %d", *through, sum)
	} else if *summaryMode && *through != "" {
		// summarize all buckets through a given date
		td, err := parseDate(*through)
		if err != nil {
			log.Print(err)
			return
		}
		tx, err := db.Begin()
		if err != nil {
			log.Fatalf("beginning sql transaction: %v", err)
		}
		result, err := summarizeAllThroughDate(tx, td)
		if err != nil {
			log.Fatalf("summarizing buckets: %v", err)
		}
		if err := tx.Commit(); err != nil {
			log.Fatalf("committing sql transaction: %v", err)
		}
		for b, v := range result {
			log.Printf("%s: %d", b, v)
		}
	} else if *summaryMode {
		// output summary of all buckets through today
		tx, err := db.Begin()
		if err != nil {
			log.Fatalf("beginning sql transaction: %v", err)
		}
		result, err := summarizeAllThroughDate(tx, time.Now())
		if err != nil {
			log.Fatalf("summarizing buckets: %v", err)
		}
		if err := tx.Commit(); err != nil {
			log.Fatalf("committing sql transaction: %v", err)
		}
		for b, v := range result {
			log.Printf("%s: %d", b, v)
		}
	} else if *zeroMode && *source != "" {
		// find when a bucket next has a zero balance
		tx, err := db.Begin()
		if err != nil {
			log.Fatalf("beginning sql transaction: %v", err)
		}
		d, err := findWhenZero(tx, *source)
		if err != nil {
			log.Fatalf("finding when a bucket becomes zero: %v", err)
		}
		log.Printf("Bucket %s next becomes zero on %v", *source, d)
	}
}

// END MAIN

// BEGIN FUNCTIONS
// insert a slice of entries
func insert(tx *sql.Tx, e entry) error {
	q := `INSERT INTO transactions
		(source, destination, happened_at, amount)
		VALUES ($1, $2, $3, $4);`
	_, err := tx.Exec(q, e.source, e.destination, e.happenedAt, e.amount)
	if err != nil {
		return fmt.Errorf("insert() - executing the insert: %w", err)
	}
	return nil
}

// get net amount of a single bucket through a given date
func summarizeBucket(tx *sql.Tx, bucket string, through time.Time) (int, error) {
	q := `SELECT sum(amount) FROM (
		SELECT amount, happened_at FROM transactions WHERE destination = $1
		UNION ALL
		SELECT -amount, happened_at from transactions where source = $1
		)
		WHERE date(happened_at) <= date($2)
		ORDER BY sum(amount) DESC;`
	row := tx.QueryRow(q, bucket, through)
	var sum int
	if err := row.Scan(&sum); err != nil {
		return -1, fmt.Errorf("summarizeBucket - querying rows: %w", err)
	}
	return sum, nil
}

// get net amounts of all buckets through a given date
func summarizeAllThroughDate(tx *sql.Tx, through time.Time) (map[string]int, error) {
	q := `SELECT account, sum(amount) FROM (
		SELECT amount, happened_at, destination AS account FROM transactions
		UNION ALL
		SELECT -amount, happened_at, source AS account FROM transactions
		)
		WHERE date(happened_at) <= date($1)
		GROUP BY account
		ORDER BY sum(amount) DESC;`
	rows, err := tx.Query(q, through)
	if err != nil {
		return nil, fmt.Errorf("summarizeAllThroughDate() - querying rows: %w", err)
	}
	result := make(map[string]int)
	for rows.Next() {
		var account string
		var total int
		if err := rows.Scan(&account, &total); err != nil {
			return nil, fmt.Errorf("summarizeAllThroughDate() - iterating through rows: %w", err)
		}
		result[account] = total
		// log.Printf("%s: %d \n", account, total)
	}
	return result, nil
}

// insert a transaction that repeats weekly or monthly
func insertRepeating(tx *sql.Tx, e entry, freq string) error {
	q := `INSERT INTO transactions
		(source, destination, happened_at, amount)
		VALUES ($1, $2, $3, $4);`
	var freqMonth int
	var freqDay int
	if freq == "monthly" {
		freqMonth = 1
		freqDay = 0
	} else if freq == "weekly" {
		freqMonth = 0
		freqDay = 7
	}
	endDate := time.Now().AddDate(2, 0, 0)
	for e.happenedAt.Before(endDate) {
		if _, err := tx.Exec(q, e.source, e.destination, e.happenedAt, e.amount); err != nil {
			return fmt.Errorf("insertRepeating() - inserting transactions: %w", err)
		}
		e.happenedAt = e.happenedAt.AddDate(0, freqMonth, freqDay)
	}
	return nil
}

// parse a date
func parseDate(s string) (time.Time, error) {
	d, err := time.Parse("2006-01-02", s)
	if err != nil {
		return time.Time{}, fmt.Errorf("parseDate() - parsing time: %w", err)
	}
	return d, nil
}

// add buckets to the db
func addBucket(tx *sql.Tx, bucket bucket) error {
	q := `INSERT INTO buckets
		(name, asset, liquidity)
		VALUES ($1, $2, $3)`
	var x int
	if bucket.asset == true {
		x = 1
	} else {
		x = 0
	}
	_, err := tx.Exec(q, bucket.name, x, bucket.liquidity)
	if err != nil {
		return fmt.Errorf("addBuckets() - executing query: %w", err)
	}
	return nil
}

// get total assets owned on a given date
func sumAssets(tx *sql.Tx, through time.Time) (int, error) {
	q := `SELECT sum(amount) FROM (
    	SELECT destination AS account, amount, happened_at
		FROM transactions
        UNION ALL
        SELECT source AS account, -amount, happened_at
		FROM transactions
        ) t
        LEFT JOIN buckets b
        ON t.account = b.name
        WHERE date(t.happened_at) < date($1) AND b.asset = 1;`
	row := tx.QueryRow(q, through)
	var sum int
	if err := row.Scan(&sum); err != nil {
		return -1, fmt.Errorf("sumAssets() - scanning rows: %w", err)
	}
	return sum, nil
}

// find the first date after today that a bucket becomes <= 0
func findWhenZero(tx *sql.Tx, bucket string) (time.Time, error) {
	todayBalance, err := summarizeBucket(tx, bucket, time.Now())
	if err != nil {
		return time.Now(), fmt.Errorf("findWhenZero() - summarizing balance today of bucket: %w", err)
	}
	if todayBalance <= 0 {
		fmt.Println("bucket is already below zero")
		return time.Now(), nil
	}

	today := convertToDate(time.Now())
	for t := today; t.Before(today.AddDate(2, 0, 0)); t = t.AddDate(0, 0, 1) {
		if balance, _ := summarizeBucket(tx, bucket, t); balance <= 0 {
			return t, nil
		}
	}
	return time.Now(), nil
}

// return a time with year, month, and day values; all other values equal 0
func convertToDate(t time.Time) time.Time {
	year := t.Year()
	month := t.Month()
	day := t.Day()
	return time.Date(year, month, day, 0, 0, 0, 0, time.Local)
}
