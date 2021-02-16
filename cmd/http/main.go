package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"ledger/pkg/budget"
	"ledger/pkg/csvreader"
	"ledger/pkg/ledger"
	"ledger/pkg/myhttp"
	"ledger/pkg/mytemplate"
	"ledger/pkg/usd"
	"ledger/pkg/utils"
	"log"
	"net/http"
	"strconv"
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

// begin react handlers
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

func (s *server) getBudget(w http.ResponseWriter, r *http.Request) {
	var entries []budget.Entry
	var startDate time.Time
	var endDate time.Time
	utils.Tx(s.db, r, func(tx *sql.Tx) (err error) {
		startDate, err = myhttp.SetStartDate(tx, r.URL.Query())
		endDate, err = myhttp.SetEndDate(tx, r.URL.Query())
		entries, err = budget.GetBudgetEntries(tx, startDate, endDate)
		return err
	})
	//
	group := struct {
		StartDate time.Time
		EndDate   time.Time
		Entries   []budget.Entry
	}{
		StartDate: startDate,
		EndDate:   endDate,
		Entries:   entries,
	}
	//
	output, err := json.Marshal(group)
	if err != nil {
		log.Printf("marshaling entries: %v", err)
	}
	// os.Stdout.Write(out)
	//
	w.Header().Add("content-type", "application/json")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Add("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	//
	if _, err := io.Copy(w, bytes.NewBuffer(output)); err != nil {
		log.Printf("writing response: %v", err)
	}
}

// get budget over time
func (s *server) getBudgetSeries(w http.ResponseWriter, r *http.Request) {
	// allocate variables
	var startDate time.Time
	var endDate time.Time
	var allCategories []string
	var filterCategories []string
	var timeInterval int
	var spendSummary []map[string]usd.USD
	//
	utils.Tx(s.db, r, func(tx *sql.Tx) (err error) {
		q := r.URL.Query()
		startDate, err = myhttp.SetStartDate(tx, q)
		endDate, err = myhttp.SetEndDate(tx, q)
		timeInterval, err = myhttp.SetTimeInterval(q)
		filterCategories, allCategories, err = myhttp.SetBudgetCategories(tx, q)
		spendSummary, err = budget.SummarizeSpendsOverTime(tx, filterCategories, startDate, endDate, timeInterval)
		return err
	})
	budgetOverTimeTable := budget.MakePlot(spendSummary, startDate, timeInterval)

	//
	group := struct {
		StartDate     time.Time
		EndDate       time.Time
		TimeInterval  int
		AllCategories []string
		Table         budget.PlotData
	}{
		StartDate:     startDate,
		EndDate:       endDate,
		TimeInterval:  timeInterval,
		AllCategories: allCategories,
		Table:         *budgetOverTimeTable,
	}
	//
	output, err := json.Marshal(group)
	if err != nil {
		log.Printf("marshaling entries: %v", err)
	}
	// os.Stdout.Write(out)
	//
	w.Header().Add("content-type", "application/json")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Add("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	//
	if _, err := io.Copy(w, bytes.NewBuffer(output)); err != nil {
		log.Printf("writing response: %v", err)
	}

}

func (s *server) insertBudgetViaJson(w http.ResponseWriter, r *http.Request) {
	type StringEntry struct {
		EntryDate   string
		Amount      string
		Category    string
		Description string
	}
	//
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}
	//
	var stringEntry StringEntry
	json.Unmarshal(body, &stringEntry)
	//
	var entry budget.Entry
	entry.EntryDate, err = time.Parse("2006-01-02", stringEntry.EntryDate)
	if err != nil {
		fmt.Printf("Parsing start time (%v)", err)
		return
	}
	amountInt, err := strconv.Atoi(stringEntry.Amount)
	if err != nil {
		fmt.Printf("Could not convert form string %s to number: %v", stringEntry.Amount, err)
		return
	}
	entry.Amount = usd.USD(amountInt)
	entry.Category = stringEntry.Category
	entry.Description = stringEntry.Description
	//
	utils.Tx(s.db, r, func(tx *sql.Tx) error {
		if err := budget.InsertEntry(tx, entry); err != nil {
			http.Error(w, fmt.Sprintf("Calling budget.InsertEntry() (%v)", err), http.StatusInternalServerError)
			return err
		}
		return nil
	})
	//
	w.Header().Add("content-type", "application/json")
	json, err := json.Marshal(entry)
	if err != nil {
		log.Printf("marshaling entry: %v", err)
	}
	if _, err := io.Copy(w, bytes.NewBuffer(json)); err != nil {
		log.Printf("writing response: %v", err)
	}
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
	http.HandleFunc("/budget.json", s.getBudget)
	http.HandleFunc("/budgetseries.json", s.getBudgetSeries)
	http.HandleFunc("/insert.json", s.insertBudgetViaJson)
	//
	// http.HandleFunc("/budgetseries", s.handleBudgetOverTime)
	//
	http.HandleFunc("/ledger", s.ledgerHandler)
	http.HandleFunc("/balance", s.balanceOverTimeHandler)
	http.HandleFunc("/ledgerseries", s.ledgerOverTimeHandler)
	//
	http.HandleFunc("/insert", mytemplate.Insert)
	http.HandleFunc("/upload_csv", s.uploadCsvHandler)
	http.HandleFunc("/insert_ledger_entry", s.insertLedgerEntryHandler)
	http.HandleFunc("/insert_budget_entry", s.insertBudgetEntryHandler)
	//
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	//
	log.Fatal(http.ListenAndServe(":8080", nil))
}
