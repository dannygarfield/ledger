package main

import (
	"database/sql"
	"fmt"
	"ledger/pkg/mytemplate"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type server struct{ db *sql.DB }

func (s *server) ledgerHandler(w http.ResponseWriter, r *http.Request) {
	tx, err := s.db.Begin()
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not open sql transaction (%v)", err), http.StatusInternalServerError)
	}
	// providing abritrary dates -- eventually these should be user inputs
	start := time.Date(1992, 8, 16, 0, 0, 0, 0, time.Local)
	end := time.Date(2024, 8, 16, 0, 0, 0, 0, time.Local)
	mytemplate.LedgerHandler(tx, start, end, w, r)
}

func (s *server) dailyLedgerHandler(w http.ResponseWriter, r *http.Request) {
	tx, err := s.db.Begin()
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not open sql transaction (%v)", err), http.StatusInternalServerError)
	}
	// providing abritrary dates -- eventually these should be user inputs
	buckets := []string{"checking", "savings", "rent"}
	start := time.Date(2020, 12, 8, 0, 0, 0, 0, time.Local)
	end := time.Date(2021, 12, 16, 0, 0, 0, 0, time.Local)
	mytemplate.DailyLedgerHandler(tx, buckets, start, end, w, r)
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
	http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %v", r.URL)
	})
	// http.HandleFunc("/insert", mytemplate.insertHandler)
	// http.HandleFunc("/ledger", mytemplate.LedgerHandler)
	http.HandleFunc("/ledger", s.ledgerHandler)
	http.HandleFunc("/dailyledger", s.dailyLedgerHandler)
	http.HandleFunc("/insert", mytemplate.InsertHandler)
	http.HandleFunc("/save", mytemplate.SaveHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
