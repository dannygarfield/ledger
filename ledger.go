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
type transaction struct {
	source		string
	destination	string
	happened_at	time.Time
	amount		int
}

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	// define flags for input from the command line
	insertMode := flag.Bool("insert", false, "insert a transaction")
	summaryMode := flag.Bool("summary", false, "get balances of all buckets")
	through_date := flag.String("through_date", "", "date through which to summarize")
	source := flag.String("source", "", "bucket from which the amount is taken")
	destination := flag.String("destination", "", "bucket into which the amount is deposited")
	happened_at := flag.String("happened_at", "", "date of transaction")
	amount := flag.Int("amount", 0, "amount in cents of the transaction")
	bucket := flag.String("bucket", "", "bucket to categorize, used with summary")
	repeat := flag.Bool("repeat", false, "make a transaction repeating")
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
	} else if *insertMode && *repeat {
		// insert repeating 1x/month through 24 months from today
		transaction, err := db.Begin()
		if err != nil {
			log.Fatalf("beginning the transaction: %v", err)
		}
		d, err := time.Parse("2006-01-02", *happened_at)
		if err != nil {
			log.Fatalf("parsing time: %v", err)
		}
		q := "INSERT INTO transactions (source, destination, happened_at, amount) VALUES ($1, $2, $3, $4);"
		// add transaction once per month for two years
		end_date := time.Now().AddDate(2, 0, 0)
		for tx_date := d; tx_date.Before(end_date); tx_date = tx_date.AddDate(0, 1, 0) {
			if _, err := transaction.Exec(q, *source, *destination, tx_date, *amount); err != nil {
				log.Fatalf("inserting the transaction: %v", err)
			}
		}
		// commit the transaction
		if err := transaction.Commit(); err != nil {
			log.Fatalf("committing the transaction: %v", err)
		}
	} else if *insertMode {
		// insert a transaction to the db
		d, err := time.Parse("2006-01-02", *happened_at)
		if err != nil {
			log.Fatalf("parsing time: %v", err)
		}
		if err := insert(db, *source, *destination, d, *amount); err != nil {
			log.Fatalf("inserting the transaction: %v", err)
		}
	} else if *summaryMode && *through_date != "" {
		// summarize all buckets through a given date
		td, err := time.Parse("2006-01-02", *through_date)
		if err != nil {
			log.Fatalf("parsing time: %v", err)
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
	} else if *summaryMode && *bucket != "" {
		// summarize a single given bucket through today
		q := `SELECT sum(amount) FROM (
		    SELECT amount, happened_at FROM transactions WHERE destination = $1
		    UNION ALL
		    SELECT -amount, happened_at from transactions where source = $1
			)
		    WHERE date(happened_at) <= date("now")
			ORDER BY sum(amount) DESC;`
		row := db.QueryRow(q, bucket)
		var sum int
		if err := row.Scan(&sum); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("total amount, all time, for '%v': %d\n", *bucket, sum)
	} else if *summaryMode {
		// output summary of all buckets through today
		fmt.Println("Summary of all accounts, all time...")
		q := `SELECT account, sum(amount) FROM (
		    SELECT amount, happened_at, destination AS account FROM transactions
		    UNION ALL
		    SELECT -amount, happened_at, source AS account FROM transactions
		    )
		    WHERE date(happened_at) <= date("now")
		    GROUP BY account
			ORDER BY sum(amount) DESC;`
		rows, err := db.Query(q)
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
	}
}

func insert(db *sql.DB, source string, destination string, happenedAt time.Time, amount int) error {
	q := "INSERT INTO transactions (source, destination, happened_at, amount) VALUES ($1, $2, $3, $4);"
	_, err := db.Exec(q, source, destination, happenedAt, amount)
	return err
}

func summary(db *sql.DB, bucket string, through time.Time) (int, error) {
	q := `
SELECT sum(amount) FROM (
SELECT amount, happened_at FROM transactions WHERE destination = $1
UNION ALL
SELECT -amount, happened_at from transactions where source = $1
)
WHERE date(happened_at) <= date($2)
ORDER BY sum(amount) DESC;`
	row := db.QueryRow(q, bucket, through)
	var sum int
	if err := row.Scan(&sum); err != nil {
		return -1, err
	}
	return sum, nil
}
