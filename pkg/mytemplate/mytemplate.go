package mytemplate

import (
	"database/sql"
	"fmt"
	"html/template"
	"ledger/pkg/ledger"
	"net/http"
	"time"
)

// define a struct to feed into a template
type DayLedger struct {
	Day       string
	LedgerMap map[string]int
}

func LedgerHandler(tx *sql.Tx, w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("pkg/mytemplate/ledger.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := prepareDayLedger(tx)
	t.Execute(w, data)
}

func DailyLedgerHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("pkg/mytemplate/dailyledger.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := map[string][]int{
		"checking": {-100, -200, -300},
		"savings":  {100, 200, 300},
		"401k":     {0, 0, 0},
	}
	// data := []map[string]int{
	// 	{bucket1: -100, bucket2: 100, bucket3: 0},
	// 	{bucket1: -200, bucket2: 200, bucket3: 0},
	// 	{bucket1: -300, bucket2: 300, bucket3: 0},
	// }

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

	// convert date to time.Time
	// convert amount to int
	// open db
	// start sql Tx
	// insert entry
	// commit

	for k, v := range formdata {
		fmt.Printf("%s: %s\n", k, v)
	}
	http.Redirect(w, r, "/ledger", http.StatusFound)
}

func prepareDayLedger(tx *sql.Tx) DayLedger {

	today := time.Now().Format("01/02/2006")

	allbuckets, _ := ledger.GetBuckets(tx)
	l, _ := ledger.SummarizeLedger(tx, allbuckets, time.Now())

	return DayLedger{today, l}
}
