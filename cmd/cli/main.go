package main

import (
	"database/sql"
	"flag"
	"ledger/pkg/csvreader"
	"ledger/pkg/ledger"
	"log"
	"time"
)

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
	csvMode := flag.Bool("csv", false, "inert a transaction using a csv")
	// zeroMode := flag.Bool("zero", false, "find when a bucket zeroes out")
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
	} else if !*insertMode && !*summaryMode{
		// instruct user to pick a mode
		log.Printf("specify one of -insert or -summary or --zero")
		return
	} else if *insertMode && *csvMode {
		entries, err := csvreader.CsvToEntries("records.csv")
		if err != nil {
			log.Fatalf("Reading csv")
		}
		// begin the sql transaction
		tx, err := db.Begin()
		if err != nil {
			log.Fatalf("beginning sql transaction: %v", err)
		}
		// insert all entries
		for _, e := range entries {
			if err := ledger.Insert(tx, e); err != nil {
				log.Fatalf("inserting single entry")
			}
		}
		// commit the sql transaction
		if err := tx.Commit(); err != nil {
			log.Fatalf("committing sql transaction: %v", err)
		}

	} else if *insertMode && *repeat != "" {
		// insert entry that repeats through 2 years from today
		d, err := ledger.ParseDate(*happenedAt)
		if err != nil {
			log.Print(err)
			return
		}
		e := ledger.Entry{
			Source:      *source,
			Destination: *destination,
			HappenedAt:  d,
			Amount:      *amount,
		}
		tx, err := db.Begin()
		if err != nil {
			log.Fatalf("beginning sql transaction: %v", err)
		}
		if err := ledger.InsertRepeating(tx, e, *repeat); err != nil {
			log.Fatalf("inserting a repeating entry: %v", err)
		}
		if err := tx.Commit(); err != nil {
			log.Fatalf("committing sql transaction: %v", err)
		}
	} else if *insertMode {
		// insert a transaction to the db
		d, err := ledger.ParseDate(*happenedAt)
		if err != nil {
			log.Print(err)
			return
		}
		e := ledger.Entry{
			Source:      *source,
			Destination: *destination,
			HappenedAt:  d,
			Amount:      *amount,
		}
		tx, err := db.Begin()
		if err != nil {
			log.Fatalf("beginning sql transaction: %v", err)
		}
		if err := ledger.Insert(tx, e); err != nil {
			log.Fatalf("inserting single entry")
		}
		if err := tx.Commit(); err != nil {
			log.Fatalf("committing sql transaction: %v", err)
		}
	} else if *summaryMode && *through != "" && *assets {
		// summarize all assets through a given date
		td, err := ledger.ParseDate(*through)
		if err != nil {
			log.Print(err)
			return
		}
		tx, err := db.Begin()
		if err != nil {
			log.Fatalf("beginning sql transaction: %v", err)
		}
		sum, err := ledger.SumAssets(tx, td)
		if err != nil {
			log.Fatalf("summing assets: %v", err)
		}
		if err := tx.Commit(); err != nil {
			log.Fatalf("committing sql transaction: %v", err)
		}
		log.Printf("All assets as of %v: %d", *through, sum)
	} else if *summaryMode && *through != "" {
		// summarize all buckets through a given date
		td, err := ledger.ParseDate(*through)
		if err != nil {
			log.Print(err)
			return
		}
		tx, err := db.Begin()
		if err != nil {
			log.Fatalf("beginning sql transaction: %v", err)
		}
		result, err := ledger.SummarizeAllThroughDate(tx, td)
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
		result, err := ledger.SummarizeAllThroughDate(tx, time.Now())
		if err != nil {
			log.Fatalf("summarizing buckets: %v", err)
		}
		if err := tx.Commit(); err != nil {
			log.Fatalf("committing sql transaction: %v", err)
		}
		for b, v := range result {
			log.Printf("%s: %d", b, v)
		}
	}
}
