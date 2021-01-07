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

	csvMode := flag.Bool("csv", false, "inert a transaction using a csv")
	filepath := flag.String("filepath", "", "path to csv file to read")
	repeat := flag.String("repeat", "", "how often an entry repeats: 'weekly' or 'monthly'")

	through := flag.String("through", "", "date through which to summarize")

	source := flag.String("source", "", "bucket from which the amount is taken")
	destination := flag.String("destination", "", "bucket into which the amount is deposited")
	entrydate := flag.String("entrydate", "", "date of transaction")
	amount := flag.Int("amount", 0, "amount in cents of the transaction")

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
	} else if !*insertMode && !*summaryMode {
		// instruct user to pick a mode
		log.Printf("specify one of -insert or -summary or --zero")
		return
	} else if *insertMode && *csvMode {
		// insert entries from a csv
		entries, err := csvreader.CsvToEntries(*filepath)
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
			if err := ledger.InsertEntry(tx, e); err != nil {
				log.Fatalf("inserting single entry")
			}
		}
		// commit the sql transaction
		if err := tx.Commit(); err != nil {
			log.Fatalf("committing sql transaction: %v", err)
		}
	} else if *insertMode && *repeat != "" {
		// insert entry that repeats through 2 years from today
		d, err := ledger.ParseDate(*entrydate)
		if err != nil {
			log.Print(err)
			return
		}
		e := ledger.Entry{
			Source:      *source,
			Destination: *destination,
			EntryDate:   d,
			Amount:      *amount,
		}
		tx, err := db.Begin()
		if err != nil {
			log.Fatalf("beginning sql transaction: %v", err)
		}
		if err := ledger.InsertRepeatingEntry(tx, e, *repeat); err != nil {
			log.Fatalf("inserting a repeating entry: %v", err)
		}
		if err := tx.Commit(); err != nil {
			log.Fatalf("committing sql transaction: %v", err)
		}
	} else if *insertMode {
		// insert a transaction to the db
		d, err := ledger.ParseDate(*entrydate)
		if err != nil {
			log.Print(err)
			return
		}
		e := ledger.Entry{
			Source:      *source,
			Destination: *destination,
			EntryDate:   d,
			Amount:      *amount,
		}
		tx, err := db.Begin()
		if err != nil {
			log.Fatalf("beginning sql transaction: %v", err)
		}
		if err := ledger.InsertEntry(tx, e); err != nil {
			log.Fatalf("inserting single entry")
		}
		if err := tx.Commit(); err != nil {
			log.Fatalf("committing sql transaction: %v", err)
		}
	} else if *summaryMode {
		bigBang := time.Date(1996, 04, 11, 0, 0, 0, 0, time.Local)
		// summarize all buckets through a given date
		td := time.Now()
		if *through != "" {
			td, err = ledger.ParseDate(*through)
			if err != nil {
				log.Print(err)
				return
			}
		}
		// begin sql transaction
		tx, err := db.Begin()
		if err != nil {
			log.Fatalf("beginning sql transaction: %v", err)
		}
		// get list of bucket names
		bucketList, err := ledger.GetBuckets(tx)
		if err != nil {
			log.Fatalf("summarizing buckets: %v", err)
		}
		// get ledger summary
		ledgerMap, err := ledger.SummarizeBalance(tx, bucketList, bigBang, td)
		if err != nil {
			log.Fatalf("summarizing buckets: %v", err)
		}
		if err := tx.Commit(); err != nil {
			log.Fatalf("committing sql transaction: %v", err)
		}
		for b, v := range ledgerMap {
			log.Printf("%s: %v", b, v)
		}
	}
}
