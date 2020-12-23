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

func LedgerHandler(tx *sql.Tx, start, end time.Time, w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("pkg/mytemplate/ledger.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not parse ledger.html (%v)", err), http.StatusInternalServerError)
		return
	}

	myledger, err := ledger.GetLedger(tx, start, end)
	data := struct {
		Start, End  time.Time
		Ledger []ledger.Entry
	}{
		start,
		end,
		myledger,
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, data)
}

func DailyLedgerHandler(tx *sql.Tx, w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("pkg/mytemplate/dailyledger.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not parse dailyledger.html (%v)", err), http.StatusInternalServerError)
		return
	}

	start := time.Date(2020, 12, 8, 0, 0, 0, 0, time.Local)
	end := time.Date(2021, 12, 16, 0, 0, 0, 0, time.Local)

	buckets, err := ledger.GetBuckets(tx)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not call GetBuckets() (%v)", err), http.StatusInternalServerError)
	}

	dailyledger, err := ledger.MakePlot(tx, buckets, start, end)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not call MakePlot() (%v)", err), http.StatusInternalServerError)
	}

	data := struct{
		Buckets []string
		DailyLedger []map[string]int
	}{
		buckets,
		dailyledger,
	}

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
