package csvwriter

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"ledger/pkg/csvreader"
	"ledger/pkg/ledger"
	"net/http"
)

func UploadCsv(tx *sql.Tx, w http.ResponseWriter, r *http.Request) error {

	r.ParseMultipartForm(10 << 20) // max 10mb files

	// retrieve file from posted form-data
	file, _, err := r.FormFile("user_csv")
	if err != nil {
		return fmt.Errorf("Error retrieving file from form-data (%v)", err)
	}
	defer file.Close()

	//  write temporary file
	tempFile, err := ioutil.TempFile("tempcsv", "upload-*.csv")
	if err != nil {
		return fmt.Errorf("Error error writing temp file (%v)", err)
	}
	defer tempFile.Close()
	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("Error copying data to tempfile (%v)", err)
	}
	tempFile.Write(fileContent)
	filepath := tempFile.Name()

	// convert to entries
	entries, err := csvreader.CsvToEntries(filepath)
	if err != nil {
		return fmt.Errorf("Could not convert csv to entries (%v)", err)
	}

	for _, e := range entries {
		err := ledger.InsertEntry(tx, e)
		if err != nil {
			return fmt.Errorf("Could not insert entries (%v)", err)
		}
	}
	return nil
}
