package csvreader

import (
	"encoding/csv"
	"fmt"
	"io"
	"ledger/pkg/ledger"
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

func InsertCsv(filepath string) ([]ledger.Entry, error) {
	// create the reader
	reader, err := csvReader(filepath)
	if err != nil {
		return nil, fmt.Errorf("creating the reader: %w", err)
	}

	// skip the header
	_, err = reader.Read()
	if err != nil {
		return nil, fmt.Errorf("An error encountered ::", err)
	}
	// Validate order -- DO LATER

	// construct slice of entries to return
	var entries []ledger.Entry

	// Read rows and construct and append Entry objects
	for i := 0; ; i += 1 {
		record, err := reader.Read()
		if err == io.EOF {
			break // reached end of the file
		} else if err != nil {
			return nil, fmt.Errorf("Reading a row: %v", err)
		}

		// convert HappenedAt value to time.Time
		happenedAt, err := ledger.ParseDate(record[2])
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
			HappenedAt:  happenedAt,
			Amount:      amount,
		}

		entries = append(entries, e)
	}
	return entries, nil
}

func printRow(row []string) {
	for _, v := range row {
		fmt.Println(v)
	}
}
