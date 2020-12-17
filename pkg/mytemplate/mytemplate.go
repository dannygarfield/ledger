package mytemplate

import (
	"database/sql"
	"html/template"
	"ledger/pkg/ledger"
	"ledger/pkg/ledgerbucket"
	"log"
	"net/http"
	"time"
	"fmt"
)

// define a struct to feed into a template
type DayLedger struct {
	Day         string
	LedgerMap map[string]int
}

func LedgerHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("pkg/mytemplate/ledger.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := prepareDayLedger()
	t.Execute(w, data)
}

func InsertHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("pkg/mytemplate/insert.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, nil)
}

func SaveHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	formdata := r.Form
	for k, v := range formdata {
		fmt.Printf("%s: %s\n", k, v)
	}
	http.Redirect(w, r, "/ledger", http.StatusFound)
}

func prepareDayLedger() DayLedger {
	db, err := sql.Open("sqlite3", "./db.sqlite3")
	if err != nil {
		log.Fatalf("opening database: %v", err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("beginning sql transaction: %v", err)
	}

	today := time.Now().Format("01/02/2006")

	allbuckets, _ := ledgerbucket.GetBuckets(tx)
	l, _ := ledger.SummarizeLedger(tx, allbuckets, time.Now())

	if err := tx.Commit(); err != nil {
		log.Fatalf("committing sql transaction: %v", err)
	}

	return DayLedger{today, l}
}
