package main

import (
	"database/sql"
	"flag"
	"ledger/pkg/csvreader"
	"ledger/pkg/ledger"
	"ledger/pkg/ledgerbucket"
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
	object := flag.String("object", "entry", "object to work on: 'entry' or 'bucket'")
	repeat := flag.String("repeat", "", "how often an entry repeats: 'weekly' or 'monthly'")

	asset := flag.Int("asset", 0, "denotes if a bucket is considered an asset (in your posession)")
	liquidity := flag.String("liquidity", "", "liquidity of a bucket: 'low', 'medium', 'full', or the empty string")

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
	} else if *insertMode && *csvMode && *object == "bucket" {
		// insert buckets from a csv
		buckets, err := csvreader.CsvToBuckets(*filepath)
		if err != nil {
			log.Fatalf("Reading csv")
		}
		// begin the sql transaction
		tx, err := db.Begin()
		if err != nil {
			log.Fatalf("beginning sql transaction: %v", err)
		}
		// insert all entries
		for _, b := range buckets {
			if err := ledgerbucket.InsertBucket(tx, b); err != nil {
				log.Fatalf("inserting single bucket")
			}
		}
		// commit the sql transaction
		if err := tx.Commit(); err != nil {
			log.Fatalf("committing sql transaction: %v", err)
		}
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
	} else if *insertMode && *object == "bucket" {
		b := ledgerbucket.Bucket{
			Name:      *source,
			Asset:     *asset,
			Liquidity: *liquidity,
		}
		tx, err := db.Begin()
		if err != nil {
			log.Fatalf("beginning sql transaction: %v", err)
		}
		if err := ledgerbucket.InsertBucket(tx, b); err != nil {
			log.Fatalf("inserting single bucket")
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
		// summarize all buckets through a given date
		td := time.Now()

		if *through != "" {
			td, err = ledger.ParseDate(*through)
			if err != nil {
				log.Print(err)
				return
			}
		}

		tx, err := db.Begin()
		if err != nil {
			log.Fatalf("beginning sql transaction: %v", err)
		}

		bucketList, err := ledgerbucket.GetBuckets(tx)
		if err != nil {
			log.Fatalf("summarizing buckets: %v", err)
		}

		// this will return nothing. second arg should take in all bucket names in buckets table
		ledgerMap, err := ledger.SummarizeLedger(tx, bucketList, td)
		if err != nil {
			log.Fatalf("summarizing buckets: %v", err)
		}
		if err := tx.Commit(); err != nil {
			log.Fatalf("committing sql transaction: %v", err)
		}
		for k, v := range ledgerMap {
			log.Printf("%s: %v", k, v)
		}
	}
}
