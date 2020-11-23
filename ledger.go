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
	source	string
	destination	string
	date	time.Time
	amount	uint
}

// String implements fmt.Stringer.
//
// It lets us "pretty-print" this structure more easily in fmt.Printf.
func (t transaction) String() string {
	return fmt.Sprintf("source=%s,destination=%s,date=%s,amount=%d\n", t.source, t.destination, t.date, t.amount)
}

// ledger represents a double-entry accounting journal.
type ledger []transaction

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	ledger := ledger{}

	insertMode := flag.Bool("insert", false, "insert a transaction")
	summaryMode := flag.Bool("summary", false, "get balances of all buckets")

	source := flag.String("source", "", "bucket from which the amount is taken")
	destination := flag.String("destination", "", "bucket into which the amount is deposited")
	date := flag.String("date", "", "date of transaction")
	amount := flag.Uint("amount", 0, "amount in cents of the transaction")
	flag.Parse()

	if *insertMode && *summaryMode {
		log.Printf("only use one of -insert or -summary")
		return
	} else if !*insertMode && !*summaryMode {
		log.Printf("specify one of -insert or -summary")
		return
	} else if *insertMode {
		d, err := time.Parse("2006-01-02", *date)
		if err != nil {
			log.Fatalf("parsing time: %v", err)
		}

		db, err := sql.Open("sqlite3", "./db.sqlite3")
		if err != nil {
			log.Fatalf("opening database: %v", err)
		}
		row := db.QueryRow("SELECT COUNT(*) FROM transactions;")
		var count int
		if err := row.Scan(&count); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("transactions: %d\n", count)

		// insert the transaction indestination the database
		q := "INSERT INTO transactions (source, destination, date, amount) VALUES ($1, $2, $3, $4);"
		if _, err := db.Exec(q, *source, *destination, d, *amount); err != nil {
			log.Fatalf("inserting the transaction: %v", err)
		}

		// query the database for all transactions
		rows, err := db.Query("SELECT source, destination, date, amount FROM transactions;")
		if err != nil {
			log.Fatalf("fetching all transactions: %v", err)
		}
		for rows.Next() {
			tx := transaction{}
			var dstring string
			if err := rows.Scan(&tx.source, &tx.destination, &dstring, &tx.amount); err != nil {
				log.Fatalf("unmarshaling row: %v", err)
			}
			d, err := time.Parse("2006-01-02", *date)
			if err != nil {
				log.Fatalf("parsing time: %v", err)
			}
			tx.date = d
			ledger = append(ledger, tx)
		}

		fmt.Printf("ledger: %s\n", ledger)
	} else if *summaryMode {
		const q = `
SELECT SUM(a.amount) FROM (
	SELECT -amount AS amount
	FROM transactions WHERE origin = "checking"
UNION
	SELECT amount
	FROM transactions WHERE destination = "checking"
) as a;`
		// execute one or more queries to summarize all buckets
		// print each bucket
		log.Printf("TODO")
	}

}
