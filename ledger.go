package main

import (
	"database/sql"
	"flag"
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
		d := parseDate(*happenedAt)
		e := entry{
			source:      *source,
			destination: *destination,
			happenedAt:  d,
			amount:      *amount,
		}
		insertRepeating(db, e, *repeat)
	} else if *insertMode {
		// insert a transaction to the db
		d := parseDate(*happenedAt)
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
		td := parseDate(*through)
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

// insert a single transaction
func insert(db *sql.DB, entries []entry) error {
	q := `INSERT INTO transactions
		(source, destination, happened_at, amount)
		VALUES ($1, $2, $3, $4);`
	tx := beginTx(db)
	for _, e := range entries {
		log.Print("entry:", e)
		_, err := tx.Exec(q, e.source, e.destination, e.happenedAt, e.amount)
		if err != nil {
			log.Fatalf("executing the insert")
		}
	}
	if err := tx.Commit(); err != nil {
		log.Fatalf("committing the transaction")
	}

	return nil
}

// get net amount of a single bucket through a given date
func summarizeBucket(db *sql.DB, bucket string, through time.Time) (int, error) {
	tx := beginTx(db)
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
		return -1, err
	}
	if err := tx.Commit(); err != nil {
		log.Fatalf("committing the transaction")
	}
	return sum, nil
}

// get net amounts of all buckets through a given date
func summarizeAllThroughDate(db *sql.DB, through time.Time) (map[string]int, error) {
	tx := beginTx(db)
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
		log.Fatalf("summarizing transactions: %v", err)
	}

	result := make(map[string]int)

	for rows.Next() {
		var account string
		var total int
		if err := rows.Scan(&account, &total); err != nil {
			log.Fatal(err)
		}
		result[account] = total
		log.Printf("%s: %d \n", account, total)
	}
	if err := tx.Commit(); err != nil {
		log.Fatalf("committing the transaction")
	}
	return result, nil
}

// insert a transaction that repeats weekly or monthly
func insertRepeating(db *sql.DB, e entry, freq string) error {
	tx := beginTx(db)
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
			log.Fatalf("inserting the transaction: %v", err)
		}
		e.happenedAt = e.happenedAt.AddDate(0, freqMonth, freqDay)
	}
	// commit the transaction
	if err := tx.Commit(); err != nil {
		log.Fatalf("committing the transaction: %v", err)
	}
	return nil
}

func parseDate(s string) time.Time {
	d, err := time.Parse("2006-01-02", s)
	if err != nil {
		log.Fatalf("parsing time: %v", err)
	}
	return d
}

func beginTx(db *sql.DB) *sql.Tx {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("beginning the sql transaction")
	}
	return tx
}

// func classifyBuckets(db *sql.DB, buckets []bucket) error {
//
// }
