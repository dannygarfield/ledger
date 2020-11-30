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
	source      string
	destination string
	happened_at        time.Time
	amount      uint
}

// tranaction implements fmt.Stringer
// It lets us "pretty-print" this structure more easily in fmt.Printf.
func (t transaction) String() string {
	return fmt.Sprintf("source=%s,destination=%s,happened_at=%s,amount=%d\n", t.source, t.destination, t.happened_at, t.amount)
}

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	// define flags for input from the command line
	insertMode := flag.Bool("insert", false, "insert a transaction")
	summaryMode := flag.Bool("summary", false, "get balances of all buckets")
	source := flag.String("source", "", "bucket from which the amount is taken")
	destination := flag.String("destination", "", "bucket into which the amount is deposited")
	happened_at := flag.String("happened_at", "", "date of transaction")
	amount := flag.Uint("amount", 0, "amount in cents of the transaction")
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
			log.Fatal("beginning the transaction: %v", err)
		}
		d, err := time.Parse("2006-01-02", *happened_at)
		if err != nil {
			log.Fatalf("parsing time: %v", err)
		}
		q := "INSERT INTO transactions (source, destination, happened_at, amount) VALUES ($1, $2, $3, $4);"
		// add transaction once per month for two years
		end_date := time.Now().AddDate(2, 0, 0)
		for tx_date := d; tx_date.Before(end_date); tx_date = tx_date.AddDate(0, 1, 0) {
			fmt.Println("end_date:", end_date)
			fmt.Println("tx_date:", tx_date)
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
		q := "INSERT INTO transactions (source, destination, happened_at, amount) VALUES ($1, $2, $3, $4);"
		if _, err := db.Exec(q, *source, *destination, d, *amount); err != nil {
			log.Fatalf("inserting the transaction: %v", err)
		}
	} else if *summaryMode && *bucket != "" {
		// summarize a single given bucket
		q := `
		SELECT sum(amount) FROM (
		   select amount from transactions where destination = $1
		   UNION
		   select -amount from transactions where source = $1
		   );`
		row := db.QueryRow(q, bucket)
		var sum int
		if err := row.Scan(&sum); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("total amount, all time, for '%v': %d\n", *bucket, sum)
	} else if *summaryMode {
		// output summary of all buckets
		fmt.Println("Summary of all accounts, all time")
		q := `SELECT sum(amount), account FROM (
			SELECT amount, destination AS account FROM transactions
			UNION
			SELECT -amount, source AS account FROM transactions
			)
			GROUP BY account;`
		rows, err := db.Query(q)
		if err != nil {
			log.Fatalf("summarizing transactions: %v", err)
		}
		for rows.Next() {
			var total int
			var account string
			if err := rows.Scan(&total, &account); err != nil {
				log.Fatal(err)
			}
			log.Printf("%s: %d \n", account, total)
		}
	}
}
