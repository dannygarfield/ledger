package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"ledger/pkg/csvreader"
	"ledger/pkg/ledger"
	"ledger/pkg/mytemplate"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type server struct{ db *sql.DB }

func (s *server) ledgerHandler(w http.ResponseWriter, r *http.Request) {
	tx, err := s.db.Begin()
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not open sql transaction (%v)", err), http.StatusInternalServerError)
	}
	mytemplate.LedgerHandler(tx, w, r)
}

func (s *server) dailyLedgerHandler(w http.ResponseWriter, r *http.Request) {
	tx, err := s.db.Begin()
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not open sql transaction (%v)", err), http.StatusInternalServerError)
	}
	mytemplate.DailyLedgerHandler(tx, w, r)
}

func (s *server) uploadCsvHandler(w http.ResponseWriter, r *http.Request) {
	tempFilepath, _ := uploadFile(w, r)
	entries, err := csvreader.CsvToEntries(tempFilepath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not convert csv to entries (%v)", err), http.StatusInternalServerError)
	}

	// fmt.Printf("entries: %v\n", entries)

	tx, err := s.db.Begin()
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not open sql transaction (%v)", err), http.StatusInternalServerError)
	}

	// fmt.Println("Opened the sql tx")

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

func uploadFile(w http.ResponseWriter, r *http.Request) (string, error) {

	r.ParseMultipartForm(10 << 20) // max 10mb files

	// retrieve file from posted form-data
	file, _, err := r.FormFile("user_csv")
	if err != nil {
		return "", fmt.Errorf("Error retrieving file from form-data (%v)", err)
	}
	defer file.Close()

	//  write temporary file
	tempFile, err := ioutil.TempFile("tempcsv", "upload-*.csv")
	if err != nil {
		return "", fmt.Errorf("Error writing temp file (%v)", err)
	}
	defer tempFile.Close()

	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("Error copying data to tempfile (%v)", err)
	}
	tempFile.Write(fileContent)
	filepath := tempFile.Name()

	return filepath, nil
}

func main() {
	db, err := sql.Open("sqlite3", "./db.sqlite3")
	if err != nil {
		log.Fatalf("opening database: %v", err)
	}

	s := &server{db: db}
	// register a path
	// instead of constructing and returning a response object, we write directly
	// to the response object (w)
	// because of this, Golang http is http2 and websockets compatible
	http.HandleFunc("/ledger", s.ledgerHandler)
	http.HandleFunc("/dailyledger", s.dailyLedgerHandler)
	http.HandleFunc("/insert", mytemplate.Insert)
	http.HandleFunc("/upload", s.uploadCsvHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
