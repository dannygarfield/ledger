package csvreader

import (
	"encoding/csv"
	"fmt"
	"io"
	"ledger/pkg/ledger"
	"ledger/pkg/ledgerbucket"
	"log"
	"os"
	"strconv"
)

func csvReader(filepath string) (*csv.Reader, error) {
	// Open the file
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("Opening the file: %w", err)
	}
	// Set up the reader
	reader := csv.NewReader(file)
	return reader, nil
}

// convert a CSV to a slice of entries
func CsvToEntries(filepath string) ([]ledger.Entry, error) {
	// create the reader
	reader, err := csvReader(filepath)
	if err != nil {
		return nil, fmt.Errorf("creating the reader: %w", err)
	}
	// read the header
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("Reading the header row: %w", err)
	}
	// Validate order of columns
	if header[0] != "source" || header[1] != "destination" || header[2] != "entrydate" || header[3] != "amount" {
		log.Fatalln("Columns must be in order: source, destination, entrydate, amount")
	}

	// construct slice of buckets to return
	var entries []ledger.Entry
	// Read rows and construct and append Entry objects
	for i := 0; ; i += 1 {
		record, err := reader.Read()
		if err == io.EOF {
			break // reached end of the file
		} else if err != nil {
			return nil, fmt.Errorf("Reading a row: %v", err)
		}
		// convert EntryDate value to time.Time
		EntryDate, err := ledger.ParseDate(record[2])
		if err != nil {
			return nil, fmt.Errorf("Parsing string to time.Time: %w", err)
		}
		// convert amount value to int
		amount, err := strconv.Atoi(record[3])
		if err != nil {
			return nil, fmt.Errorf("Converting string to int: %w", err)
		}
		// construct the entry
		e := ledger.Entry{
			Source:      record[0],
			Destination: record[1],
			EntryDate:  EntryDate,
			Amount:      amount,
		}
		entries = append(entries, e)
	}
	return entries, nil
}

// convert a CSV to a slice of buckets
func CsvToBuckets(filepath string) ([]ledgerbucket.Bucket, error) {
	// create the reader
	reader, err := csvReader(filepath)
	if err != nil {
		return nil, fmt.Errorf("creating the reader: %w", err)
	}
	// read the header
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("Reading the header row: %w", err)
	}
	// Validate order of columns
	if header[0] != "name" || header[1] != "asset" || header[2] != "liquidity" {
		log.Fatalln("Columns must be in order: source, destination, EntryDate, amount")
	}
	// construct slice of entries to return
	var buckets []ledgerbucket.Bucket
	// Read rows and construct and append Entry objects
	for i := 0; ; i += 1 {
		record, err := reader.Read()
		if err == io.EOF {
			break // reached end of the file
		} else if err != nil {
			return nil, fmt.Errorf("Reading a row: %v", err)
		}

		// convert string to int
		a, err := strconv.Atoi(record[1])
		if err != nil {
			return nil, fmt.Errorf("Converting string to int: %w", err)
		}
		// construct the bucket
		b := ledgerbucket.Bucket{
			Name:      record[0],
			Asset:     a,
			Liquidity: record[2],
		}
		buckets = append(buckets, b)
	}
	return buckets, nil
}

func printRow(row []string) {
	for _, v := range row {
		fmt.Println(v)
	}
}
