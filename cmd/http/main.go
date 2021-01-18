package main

import (
	"database/sql"
	"fmt"
	"ledger/pkg/budget"
	"ledger/pkg/csvreader"
	"ledger/pkg/ledger"
	"ledger/pkg/mytemplate"
	"ledger/pkg/utils"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type server struct{ db *sql.DB }

// ledger handlers
func (s *server) ledgerHandler(w http.ResponseWriter, r *http.Request) {
	utils.Tx(s.db, r, func(tx *sql.Tx) error {
		err := mytemplate.Ledger(tx, w, r)
		if err != nil {
			http.Error(w, fmt.Sprintf("Calling mytemplate.LedgerHandler (%v)", err), http.StatusInternalServerError)
			return err
		}
		return nil
	})
}

func (s *server) insertLedgerEntryHandler(w http.ResponseWriter, r *http.Request) {
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

func (s *server) balanceOverTimeHandler(w http.ResponseWriter, r *http.Request) {
	utils.Tx(s.db, r, func(tx *sql.Tx) error {
		if err := mytemplate.BalanceOverTime(tx, w, r); err != nil {
			http.Error(w, fmt.Sprintf("Calling mytemplate.BalanceOverTime (%v)", err), http.StatusInternalServerError)
			return err
		}
		return nil
	})
}

func (s *server) ledgerOverTimeHandler(w http.ResponseWriter, r *http.Request) {
	utils.Tx(s.db, r, func(tx *sql.Tx) error {
		if err := mytemplate.LedgerOverTime(tx, w, r); err != nil {
			http.Error(w, fmt.Sprintf("Calling mytemplate.LedgerOverTime (%v)", err), http.StatusInternalServerError)
			return err
		}
		return nil
	})
}

// budget handlers
func (s *server) budgetHandler(w http.ResponseWriter, r *http.Request) {
	utils.Tx(s.db, r, func(tx *sql.Tx) error {
		err := mytemplate.Budget(tx, w, r)
		if err != nil {
			http.Error(w, fmt.Sprintf("Calling mytemplate.Budget() (%v)", err), http.StatusInternalServerError)
			return err
		}
		return nil
	})
}

func (s *server) insertBudgetEntryHandler(w http.ResponseWriter, r *http.Request) {
	entry, err := budget.PrepareEntryForInsert(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Calling budget.PrepareEntryForInsert() (%v)", err), http.StatusInternalServerError)
		return
	}

	utils.Tx(s.db, r, func(tx *sql.Tx) error {
		if err := budget.InsertEntry(tx, entry); err != nil {
			http.Error(w, fmt.Sprintf("Calling budget.InsertEntry() (%v)", err), http.StatusInternalServerError)
			return err
		}
		return nil
	})
	// mytemplate.Insert(w, r)
	s.ledgerHandler(w, r)
}

func (s *server) budgetOverTimeHandler(w http.ResponseWriter, r *http.Request) {
	utils.Tx(s.db, r, func(tx *sql.Tx) error {
		if err := mytemplate.BudgetOverTime(tx, w, r); err != nil {
			http.Error(w, fmt.Sprintf("Calling mytemplate.BudgetOverTime (%v)", err), http.StatusInternalServerError)
			return err
		}
		return nil
	})
}

// insert by CSV
func (s *server) uploadCsvHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20) // max upload 10mb
	// create tempfile and return filepath
	filepath, err := csvreader.CreateTempFile(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Calling csvreader.CreateTempFile() (%v)", err), http.StatusInternalServerError)
		return
	}
	// call ledger or budget uploader
	if len(r.PostForm["entry_type"]) > 0 && r.PostForm["entry_type"][0] == "ledger" {
		fmt.Println("uploading ledger entries...")
		// convert csv to ledger entries
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
		fmt.Println("success")
	} else {
		fmt.Println("uploading budget entries...")
		entries, err := csvreader.CsvToBudgetEntries(filepath)
		if err != nil {
			http.Error(w, fmt.Sprintf("Calling csvreader.CsvToBudgetEntries() (%v)", err), http.StatusInternalServerError)
			return
		}
		// insert entries
		utils.Tx(s.db, r, func(tx *sql.Tx) error {
			for _, e := range entries {
				err := budget.InsertEntry(tx, e)
				if err != nil {
					http.Error(w, fmt.Sprintf("Calling budget.InsertEntry (%v)", err), http.StatusInternalServerError)
					return err
				}
			}
			return nil
		})
		fmt.Println("success")
	}
	mytemplate.Insert(w, r)
}

func main() {
	db, err := sql.Open("sqlite3", "./db.sqlite3")
	if err != nil {
		log.Fatalf("opening database: %v", err)
	}

	s := &server{db: db}

	http.HandleFunc("/ledger", s.ledgerHandler)
	http.HandleFunc("/balance", s.balanceOverTimeHandler)
	http.HandleFunc("/ledgerseries", s.ledgerOverTimeHandler)
	http.HandleFunc("/budget", s.budgetHandler)
	http.HandleFunc("/budgetseries", s.budgetOverTimeHandler)
	http.HandleFunc("/insert", mytemplate.Insert)
	http.HandleFunc("/upload_csv", s.uploadCsvHandler)
	http.HandleFunc("/insert_ledger_entry", s.insertLedgerEntryHandler)
	http.HandleFunc("/insert_budget_entry", s.insertBudgetEntryHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
