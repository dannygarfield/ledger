package main

import (
	"database/sql"
	"fmt"
	"ledger/pkg/csvwriter"
	"ledger/pkg/ledger"
	"ledger/pkg/mytemplate"
	"ledger/pkg/utils"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type server struct{ db *sql.DB }

func (s *server) ledgerHandler(w http.ResponseWriter, r *http.Request) {
	utils.Tx(s.db, r, func(tx *sql.Tx) error {
		err := mytemplate.LedgerHandler(tx, w, r)
		if err != nil {
			http.Error(w, fmt.Sprintf("Calling mytemplate.LedgerHandler (%v)", err), http.StatusInternalServerError)
			return err
		}
		return nil
	})
}

func (s *server) dailyLedgerHandler(w http.ResponseWriter, r *http.Request) {
	utils.Tx(s.db, r, func(tx *sql.Tx) error {
		if err := mytemplate.DailyLedgerHandler(tx, w, r); err != nil {
			http.Error(w, fmt.Sprintf("Calling mytemplate.DailyLedgerHandler (%v)", err), http.StatusInternalServerError)
			return err
		}
		return nil
	})
}

func (s *server) uploadCsvHandler(w http.ResponseWriter, r *http.Request) {
	utils.Tx(s.db, r, func(tx *sql.Tx) error {
		if err := csvwriter.UploadCsv(tx, w, r); err != nil {
			http.Error(w, fmt.Sprintf("Calling csvwriter.UploadCsv (%v)", err), http.StatusInternalServerError)
			return err
		}
		return nil
	})
	mytemplate.Insert(w, r)
}

func (s *server) uploadEntryHandler(w http.ResponseWriter, r *http.Request) {
	entry, err := ledger.PrepareEntryForInsert(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Calling ledger.PrepareEntryForInsert() (%v)", err), http.StatusInternalServerError)
		return
	}

	utils.Tx(s.db, r, func(tx *sql.Tx) error {
		if err := ledger.InsertEntry(tx, entry); err != nil {
			http.Error(w, fmt.Sprintf("Calling ledger.InsertEntry() (%v)", err), http.StatusInternalServerError)
			return err
		}
		return nil
	})
	// mytemplate.Insert(w, r)
	s.ledgerHandler(w, r)
}

func main() {
	db, err := sql.Open("sqlite3", "./db.sqlite3")
	if err != nil {
		log.Fatalf("opening database: %v", err)
	}

	s := &server{db: db}

	http.HandleFunc("/ledger", s.ledgerHandler)
	http.HandleFunc("/dailyledger", s.dailyLedgerHandler)
	http.HandleFunc("/insert", mytemplate.Insert)
	http.HandleFunc("/upload_entries_csv", s.uploadCsvHandler)
	http.HandleFunc("/upload_entry", s.uploadEntryHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
