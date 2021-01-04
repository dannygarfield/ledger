package csvwriter

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"ledger/pkg/csvreader"
	"ledger/pkg/ledger"
	"net/http"
)

func UploadCsv(tx *sql.Tx, w http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(10 << 20) // max 10mb files

	// retrieve file from posted form-data
	file, _, err := r.FormFile("user_csv")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving file from form-data (%v)", err), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	//  write temporary file
	tempFile, err := ioutil.TempFile("tempcsv", "upload-*.csv")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error writing temp file (%v)", err), http.StatusInternalServerError)
		return
	}
	defer tempFile.Close()
	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error copying data to tempfile (%v)", err), http.StatusInternalServerError)
		return
	}
	tempFile.Write(fileContent)
	filepath := tempFile.Name()

	// convert to entries
	entries, err := csvreader.CsvToEntries(filepath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not convert csv to entries (%v)", err), http.StatusInternalServerError)
	}

	for _, e := range entries {
		err := ledger.InsertEntry(tx, e)
		if err != nil {
			http.Error(w, fmt.Sprintf("Could not insert entries (%v)", err), http.StatusInternalServerError)
		}
	}
	if err := tx.Commit(); err != nil {
		http.Error(w, fmt.Sprintf("Could not commit sql transaction (%v)", err), http.StatusInternalServerError)
	} else {
		html := `<p>successfully uploaded file</p>
			<p>Return to <a href="/insert">insert</a></p>
			<p>View <a href="/ledger">ledger</a></p>
			<p>View <a href="/dailyledger">dailyledger</a></p>`

		fmt.Fprintf(w, html)
	}

}
