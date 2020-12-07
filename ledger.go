package main

import (
	"database/sql"
	"flag"
	"fmt"
	"ledger/sqlstatements"
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
	} else if !*insertMode && !*summaryMode {
		// instruct user to pick a mode
		log.Printf("specify one of -insert or -summary")
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
		insertRepeating(db, e, *repeat)
	} else if *insertMode {
		// insert a transaction to the db
		d, err := parseDate(*happenedAt)
		if err != nil {
			log.Print(err)
			return
		}
		es := []entry{
			{
				source:      *source,
				destination: *destination,
				happenedAt:  d,
				amount:      *amount,
			},
		}
		if err := insert(db, es); err != nil {
			log.Fatalf("inserting single entry")
		}
	} else if *summaryMode && *through != "" {
		// summarize all buckets through a given date
		td, err := parseDate(*through)
		if err != nil {
			log.Print(err)
			return
		}
		q := `SELECT account, sum(amount) FROM (
		    SELECT amount, happened_at, destination AS account FROM transactions
		    UNION ALL
		    SELECT -amount, happened_at, source AS account FROM transactions
		    )
		    WHERE date(happened_at) <= date($1)
		    GROUP BY account
			ORDER BY sum(amount) DESC;`
		rows, err := db.Query(q, td)
		if err != nil {
			log.Fatalf("summarizing transactions: %v", err)
		}
		for rows.Next() {
			var account string
			var total int
			if err := rows.Scan(&account, &total); err != nil {
				log.Fatal(err)
			}
			log.Printf("%s: %d \n", account, total)
		}
	} else if *summaryMode {
		// output summary of all buckets through today
		summarizeAllThroughDate(db, time.Now())
	}
}

// END MAIN

// BEGIN FUNCTIONS
// insert a slice of entries
func insert(tx *sql.Tx, e entry) error {
	q := `INSERT INTO transactions
		(source, destination, happened_at, amount)
		VALUES ($1, $2, $3, $4);`
	tx, err := sqlstatements.BeginTx(db)
	if err != nil {
		return fmt.Errorf("insert() - beginning the sql tx: %w", err)
	}
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
func summarizeAllThroughDate(db *sql.DB, through time.Time) (map[string]int, error) {
	tx, err := sqlstatements.BeginTx(db)
	if err != nil {
		return nil, fmt.Errorf("summarizeAllThroughDate() - beginning the sql tx: %w", err)
	}
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
		log.Printf("%s: %d \n", account, total)
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("summarizeAllThroughDate() - committing sql tx: %w", err)
	}
	return result, nil
}

// insert a transaction that repeats weekly or monthly
func insertRepeating(db *sql.DB, e entry, freq string) error {
	tx, err := sqlstatements.BeginTx(db)
	if err != nil {
		return fmt.Errorf("insertRepeating - beginning the sql tx: %w", err)
	}
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
	// commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("insertRepeating() - committing sql tx: %w", err)
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
func addBuckets(db *sql.DB, buckets []bucket) error {
	tx, err := sqlstatements.BeginTx(db)
	if err != nil {
		return fmt.Errorf("addBuckets() - beginning the sql tx: %w", err)
	}
	q := `INSERT INTO buckets
		(name, asset, liquidity)
		VALUES ($1, $2, $3)`
	var x int
	for _, b := range buckets {
		if b.asset == true {
			x = 1
		} else {
			x = 0
		}
		_, err := tx.Exec(q, b.name, x, b.liquidity)
		if err != nil {
			return fmt.Errorf("addBuckets() - executing query: %w", err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("addBuckets() - committing the transaction: %w", err)
	}
	return nil
}

// get total assets owned on a given date
func sumAssets(db *sql.DB, through time.Time) (int, error) {
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
	tx, err := sqlstatements.BeginTx(db)
	if err != nil {
		return -1, fmt.Errorf("sumAssets() - beginning sql tx: %w", err)
	}
	row := tx.QueryRow(q, through)
	var sum int
	if err := row.Scan(&sum); err != nil {
		return -1, fmt.Errorf("sumAssets() - scanning rows: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return -1, fmt.Errorf("sumAssets() - committing the transaction: %w", err)
	}
	return sum, nil
}

// find the first date after today that a bucket becomes <= 0
func findWhenZero(db *sql.DB, bucket string) (time.Time, error) {
	todayBalance, err := summarizeBucket(db, bucket, time.Now())
	if err != nil {
		return time.Now(), fmt.Errorf("findWhenZero() - summarizing balance today of bucket: %w", err)
	}
	if todayBalance <= 0 {
		fmt.Println("bucket is already below zero")
		return time.Now(), nil
	}

	today := convertToDate(time.Now())
	for t := today; t.Before(today.AddDate(2, 0, 0)); t = t.AddDate(0, 0, 1) {
		if balance, _ := summarizeBucket(db, bucket, t); balance <= 0 {
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
