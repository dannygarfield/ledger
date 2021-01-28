package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"ledger/pkg/budget"
	"ledger/pkg/csvreader"
	"ledger/pkg/ledger"
	"ledger/pkg/myhttp"
	"ledger/pkg/mytemplate"
	"ledger/pkg/utils"
	"log"
	"net/http"
	"text/template"
	"time"

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

func (s *server) handleBudgetList(w http.ResponseWriter, r *http.Request) {
	utils.Tx(s.db, r, func(tx *sql.Tx) error {
		err := myhttp.HandleBudgetList(tx, r, w)
		if err != nil {
			return fmt.Errorf("Calling myhttp.HandleBudgetList (%v)", err)
		}
		return nil
	})
}

func (s *server) handleBudgetOverTime(w http.ResponseWriter, r *http.Request) {
	utils.Tx(s.db, r, func(tx *sql.Tx) error {
		err := myhttp.HandleBudgetOverTime(tx, r, w)
		if err != nil {
			return fmt.Errorf("Could not call myhttp.HandleBudgetList: %v", err)
		}
		return nil
	})
}

func (s *server) appHandler(w http.ResponseWriter, r *http.Request) {
	err := func() error {
		// parse html template
		t, err := template.ParseFiles("pkg/mytemplate/app.html")
		if err != nil {
			return fmt.Errorf("Could not parse app.html (%v)", err)
		}
		// execute template
		if err = t.Execute(w, nil); err != nil {
			return fmt.Errorf("Could not execute html template")
		}
		return nil
	}()
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not handle route /app: %v", err), http.StatusInternalServerError)
	}
}

func (s *server) dataHandler(w http.ResponseWriter, r *http.Request) {
	jsonData, err := json.Marshal(struct{ Color string }{Color: "darksalmon"})
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not call json.Marshal: %v", err), http.StatusInternalServerError)
	}
	if _, err := w.Write(jsonData); err != nil {
		fmt.Printf("Could not write to response: %v", err)
	}
}

func (s *server) getBudget(w http.ResponseWriter, r *http.Request) {
	startDate, endDate := utils.BigBang, time.Now()
	var entries []budget.Entry
	utils.Tx(s.db, r, func(tx *sql.Tx) (err error) {
		entries, err = budget.GetBudgetEntries(tx, startDate, endDate)
		return err
	})
	w.Header().Add("content-type", "application/json")
	out, err := json.Marshal(entries)
	if err != nil {
		log.Printf("marshaling entries: %v", err)
	}
	if _, err := io.Copy(w, bytes.NewBuffer(out)); err != nil {
		log.Printf("writing response: %v", err)
	}
}

func (s *server) insertBudgetViaJson(w http.ResponseWriter, r *http.Request) {
	form := budget.Entry
	json.Unmarshal(r.Body, form)
	utils.Tx(s.db, r, func(tx *sql.Tx) (err error) {
		err = budget.InsertEntry(tx, form)
		return err
	})
	w.WriteHeader(200)
}

func main() {
	db, err := sql.Open("sqlite3", "./db.sqlite3")
	if err != nil {
		log.Fatalf("opening database: %v", err)
	}
	//
	s := &server{db: db}
	//
	http.HandleFunc("/app", s.appHandler)
	http.HandleFunc("/data", s.dataHandler)
	http.HandleFunc("/budget.json", s.getBudget)
	http.HandleFunc("/insert.json", s.insertBudgetViaJson)
	//
	http.HandleFunc("/budget", s.handleBudgetList)
	http.HandleFunc("/budgetseries", s.handleBudgetOverTime)
	//
	http.HandleFunc("/ledger", s.ledgerHandler)
	http.HandleFunc("/balance", s.balanceOverTimeHandler)
	http.HandleFunc("/ledgerseries", s.ledgerOverTimeHandler)
	//
	http.HandleFunc("/insert", mytemplate.Insert)
	http.HandleFunc("/upload_csv", s.uploadCsvHandler)
	http.HandleFunc("/insert_ledger_entry", s.insertLedgerEntryHandler)
	http.HandleFunc("/insert_budget_entry", s.insertBudgetEntryHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
