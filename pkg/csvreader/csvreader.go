package csvreader

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"ledger/pkg/budget"
	"ledger/pkg/ledger"
	"ledger/pkg/usd"
	"ledger/pkg/utils"
	"net/http"
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

// convert a CSV to a slice of ledger entries
func CsvToLedgerEntries(filepath string) ([]ledger.Entry, error) {
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
		return nil, fmt.Errorf("Columns must be in order: source, destination, entrydate, amount (%w)", err)
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
		EntryDate, err := utils.ParseDate(record[2])
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
			EntryDate:   EntryDate,
			Amount:      amount,
		}
		entries = append(entries, e)
	}
	return entries, nil
}

// convert a CSV to a slice of budget entries
func CsvToBudgetEntries(filepath string) ([]budget.Entry, error) {
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
	if header[0] != "entrydate" || header[1] != "amount" || header[2] != "category" || header[3] != "description" {
		return nil, fmt.Errorf("Columns must be in order: entrydate, amount, category, description (%w)", err)
	}

	// construct slice of buckets to return
	var entries []budget.Entry
	// Read rows and construct and append Entry objects
	for i := 0; ; i += 1 {
		record, err := reader.Read()
		if err == io.EOF {
			break // reached end of the file
		} else if err != nil {
			return nil, fmt.Errorf("Reading a row: %v", err)
		}
		// convert EntryDate value to time.Time
		EntryDate, err := utils.ParseDate(record[0])
		if err != nil {
			return nil, fmt.Errorf("Parsing string to time.Time: %w", err)
		}
		// convert amount value to int
		amount, err := usd.StringToUsd(record[1])
		if err != nil {
			return nil, fmt.Errorf("Calling usd.StringToUsd: %w", err)
		}
		// construct the entry
		e := budget.Entry{
			EntryDate:   EntryDate,
			Amount:      amount,
			Category:    record[2],
			Description: record[3],
		}
		entries = append(entries, e)
	}
	return entries, nil
}

func CreateTempFile(r *http.Request) (string, error) {
	// retrieve file from posted form-data
	// expects server has already called r.ParseMultipartForm()
	file, _, err := r.FormFile("user_csv")
	if err != nil {
		return "", fmt.Errorf("Error retrieving file from form-data (%v)", err)
	}
	defer file.Close()

	//  write temporary file
	tempFile, err := ioutil.TempFile("tempcsv", "upload-*.csv")
	if err != nil {
		return "", fmt.Errorf("Error error writing temp file (%v)", err)
	}
	defer tempFile.Close()
	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("Error copying data to tempfile (%v)", err)
	}
	tempFile.Write(fileContent)
	return tempFile.Name(), nil
}
