package mytemplate

import (
	"database/sql"
	"fmt"
	"html/template"
	"ledger/pkg/ledger"
	"net/http"
	"time"
)

// display a ledger on a single day
func LedgerHandler(tx *sql.Tx, w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("pkg/mytemplate/ledger.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not parse ledger.html (%v)", err), http.StatusInternalServerError)
		return
	}

	start := time.Date(1992, 8, 16, 0, 0, 0, 0, time.Local)
	end := time.Date(2024, 8, 16, 0, 0, 0, 0, time.Local)

	myledger, err := ledger.GetLedger(tx, start, end)
	data := struct {
		Start, End time.Time
		Ledger     []ledger.Entry
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

// display the ledger over the course of 2+ days
func DailyLedgerHandler(tx *sql.Tx, w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("pkg/mytemplate/dailyledger.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not parse dailyledger.html (%v)", err), http.StatusInternalServerError)
		return
	}
	// arbitrary start and end dates for now
	start := time.Date(2020, 12, 8, 0, 0, 0, 0, time.Local)
	end := time.Date(2021, 12, 16, 0, 0, 0, 0, time.Local)
	// get all buckets
	buckets, err := ledger.GetBuckets(tx)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not call GetBuckets() (%v)", err), http.StatusInternalServerError)
	}
	summary, err := ledger.SummarizeLedgerOverTime(tx, buckets, start, end)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not call SummarizeLedgerOverTime() (%v)", err), http.StatusInternalServerError)
	}

	plot := ledger.MakePlot(summary, start)

	// execute template
	t.Execute(w, plot)
}
