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
func LedgerHandler(tx *sql.Tx, w http.ResponseWriter, r *http.Request) error {
	t, err := template.ParseFiles("pkg/mytemplate/ledger.html")
	if err != nil {
		return fmt.Errorf("Could not parse ledger.html (%v)", err)
	}

	start := time.Date(1992, 8, 16, 0, 0, 0, 0, time.Local)
	end := time.Date(2024, 8, 16, 0, 0, 0, 0, time.Local)

	myledger, err := ledger.GetLedger(tx, start, end)
	if err != nil {
		return fmt.Errorf("Calling ledger.GetLedger() (%v)", err)
	}
	data := struct {
		Start, End time.Time
		Ledger     []ledger.Entry
	}{
		start,
		end,
		myledger,
	}
	if err = t.Execute(w, data); err != nil {
		return fmt.Errorf("Could not Execute template (%v)", err)
	}
	return nil
}

// display the ledger over the course of 2+ days
func DailyLedgerHandler(tx *sql.Tx, w http.ResponseWriter, r *http.Request) error {
	t, err := template.ParseFiles("pkg/mytemplate/dailyledger.html")
	if err != nil {
		return fmt.Errorf("Could not parse dailyledger.html (%v)", err)
	}
	// arbitrary start and end dates for now
	start := time.Date(2020, 12, 8, 0, 0, 0, 0, time.Local)
	end := time.Date(2021, 12, 16, 0, 0, 0, 0, time.Local)
	// get all buckets
	buckets, err := ledger.GetBuckets(tx)
	if err != nil {
		return fmt.Errorf("Calling ledger.GetBuckets() (%v)", err)
	}
	summary, err := ledger.SummarizeLedgerOverTime(tx, buckets, start, end)
	if err != nil {
		return fmt.Errorf("Calling ledger.SummarizeLedgerOverTime (%v)", err)
	}
	plot := ledger.MakePlot(summary, start)
	// execute template
	if err = t.Execute(w, plot); err != nil {
		return fmt.Errorf("Could not Execute template (%v)", err)
	}
	return nil
}

func Insert(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("pkg/mytemplate/insert.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not parse insert.html (%v)", err), http.StatusInternalServerError)
		return
	}
	t.Execute(w, nil)
}
