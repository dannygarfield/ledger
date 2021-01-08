package main

import (
	"database/sql"
	"fmt"
	"ledger/pkg/csvreader"
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

func (s *server) dailyBalanceHandler(w http.ResponseWriter, r *http.Request) {
	utils.Tx(s.db, r, func(tx *sql.Tx) error {
		if err := mytemplate.DailyBalanceHandler(tx, w, r); err != nil {
			http.Error(w, fmt.Sprintf("Calling mytemplate.DailyBalanceHandler (%v)", err), http.StatusInternalServerError)
			return err
		}
		return nil
	})
}

func (s *server) ledgerSeriesHandler(w http.ResponseWriter, r *http.Request) {
	utils.Tx(s.db, r, func(tx *sql.Tx) error {
		if err := mytemplate.LedgerSeriesHandler(tx, w, r); err != nil {
			http.Error(w, fmt.Sprintf("Calling mytemplate.LedgerSeriesHandler (%v)", err), http.StatusInternalServerError)
			return err
		}
		return nil
	})
}

func (s *server) uploadCsvToLedgerHandler(w http.ResponseWriter, r *http.Request) {
	// create tempfile and return filepath: CreateTempFile()
	filepath, err := csvreader.CreateTempFile(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Calling csvreader.CreateTempFile() (%v)", err), http.StatusInternalServerError)
		return
	}
	// create ledger entries from file: CsvToLedgerEntries()
	entries, err := csvreader.CsvToLedgerEntries(filepath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Calling csvreader.CsvToLedgerEntries() (%v)", err), http.StatusInternalServerError)
		return
	}
	// insert entries
	utils.Tx(s.db, r, func(tx *sql.Tx) error {
		for _, e := range entries {
			err := ledger.InsertEntry(tx, e)
			if err != nil {
				http.Error(w, fmt.Sprintf("Calling ledger.InsertEntry (%v)", err), http.StatusInternalServerError)
				return err
			}
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
	http.HandleFunc("/dailybalance", s.dailyBalanceHandler)
	http.HandleFunc("/ledgerseries", s.ledgerSeriesHandler)
	http.HandleFunc("/insert", mytemplate.Insert)
	http.HandleFunc("/upload_ledger_entries_csv", s.uploadCsvToLedgerHandler)
	http.HandleFunc("/upload_entry", s.uploadEntryHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
