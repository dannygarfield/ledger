package main

import (
	"database/sql"
	"fmt"
	"ledger/pkg/csvwriter"
	"ledger/pkg/ledger"
	"ledger/pkg/mytemplate"
	"log"
	"net/http"
	"strconv"
	"time"

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
	tx, err := s.db.Begin()
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not open sql transaction (%v)", err), http.StatusInternalServerError)
	}

	csvwriter.UploadCsv(tx, w, r)

}

func (s *server) uploadEntryHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	entrydate, _ := time.Parse("2006-01-02", r.PostForm["happened_at"][0])
	amount, err := strconv.Atoi(r.PostForm["amount"][0])
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not convert amount field to int (%v)", err), http.StatusInternalServerError)
	}

	entry := ledger.Entry{
		Source:      r.PostForm["source"][0],
		Destination: r.PostForm["destination"][0],
		EntryDate:   entrydate,
		Amount:      amount,
	}

	tx, err := s.db.Begin()
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not open sql transaction (%v)", err), http.StatusInternalServerError)
	}

	if err := ledger.InsertEntry(tx, entry); err != nil {
		http.Error(w, fmt.Sprintf("Could not insert entries (%v)", err), http.StatusInternalServerError)
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
	http.HandleFunc("/upload_entries_csv", s.uploadCsvHandler)
	http.HandleFunc("/upload_entry", s.uploadEntryHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
