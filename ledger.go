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
	from	string
	to	string
	date	time.Time
	amt	uint
}

// String implements fmt.Stringer.
//
// It lets us "pretty-print" this structure more easily in fmt.Printf.
func (t transaction) String() string {
	return fmt.Sprintf("from=%s,to=%s,date=%s,amt=%d\n", t.from, t.to, t.date, t.amt)
}

// ledger represents a double-entry accounting journal.
type ledger []transaction

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	ledger := ledger{}

	var from = flag.String("from", "", "bucket from which the amount is taken")
	var to = flag.String("to", "", "bucket into which the amount is deposited")
	var date = flag.String("date", "", "date of transaction")
	var amt = flag.Uint("amt", 0, "amount in cents of the transaction")
	flag.Parse()

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

	// insert the transaction into the database
	q := "INSERT INTO transactions (origin, destination, happened_at, amount) VALUES ($1, $2, $3, $4);"
	if _, err := db.Exec(q, *from, *to, d, *amt); err != nil {
		log.Fatalf("inserting the transaction: %v", err)
	}

	// query the database for all transactions
	rows, err := db.Query("SELECT origin, destination, happened_at, amount FROM transactions;")
	if err != nil {
		log.Fatalf("fetching all transactions: %v", err)
	}
	for rows.Next() {
		tx := transaction{}
		var dstring string
		if err := rows.Scan(&tx.from, &tx.to, &dstring, &tx.amt); err != nil {
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
}
