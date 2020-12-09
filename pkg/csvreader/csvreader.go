package csvreader

import (
	"encoding/csv"
	"fmt"
	"io"
	"ledger/pkg/ledger"
	"os"
	"strconv"
)

func CsvReaderRow() error {
	// Open the file
	file, err := os.Open("records.csv")
	if err != nil {
		return fmt.Errorf("Opening the file: %w", err)
	}

	// Set up the reader
	reader := csv.NewReader(file)

	// Read the first row
	header, err := reader.Read()
	if err != nil {
		return fmt.Errorf("An error encountered ::", err)
	}
	// Validate order -- DO LATER
	fmt.Printf("HEADERS : %v \n", header)
	// printRow(header)

	// Read rows and construct Entry objects
	for i := 0; ; i += 1 {
		record, err := reader.Read()
		fmt.Printf("ROW : %v \n", record)
		if err == io.EOF {
			break // reached end of the file
		} else if err != nil {
			return fmt.Errorf("Reading a row: %v", err)
		}

		// convert HappenedAt value to time.Time
		happenedAt, err := ledger.ParseDate(record[2])
		if err != nil {
			return fmt.Errorf("Parsing string to time.Time: %w", err)
		}

		// convert amount value to int
		amount, err := strconv.Atoi(record[3])
		if err != nil {
			fmt.Errorf("Converting string to int: %w", err)
		}

		// construct the entry
		e := ledger.Entry{
			Source:      record[0],
			Destination: record[1],
			HappenedAt:  happenedAt,
			Amount:      amount,
		}
		fmt.Printf("Entry: %v \n", e)
	}
	return nil
}

func printRow(row []string) {
	for _, v := range row {
		fmt.Println(v)
	}
}
