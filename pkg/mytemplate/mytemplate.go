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
	// parse form
	r.ParseForm()
	fmt.Println("PostForm:", r.PostForm)
	formStart := r.PostForm["start"]
	formEnd := r.PostForm["end"]

	start := time.Date(1992, 8, 16, 0, 0, 0, 0, time.Local)
	if len(formStart) > 0 && formStart[0] != "" {
		start, err = time.Parse("2006-01-02", r.PostForm["start"][0])
		if err != nil {
			return fmt.Errorf("Parsing start time (%v)", err)
		}
	}
	end := time.Date(2024, 8, 16, 0, 0, 0, 0, time.Local)
	if len(formEnd) > 0 && formEnd[0] != "" {
		end, err = time.Parse("2006-01-02", r.PostForm["end"][0])
		if err != nil {
			return fmt.Errorf("Parsing end time (%v)", err)
		}
	}

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

	r.ParseForm()
	fmt.Println("PostForm:", r.PostForm)
	formStart := r.PostForm["start"]
	formEnd := r.PostForm["end"]
	formBuckets := r.PostForm["buckets"]

	start := time.Now().AddDate(0, -1, 0)
	if len(formStart) > 0 && formStart[0] != "" {
		start, err = time.Parse("2006-01-02", r.PostForm["start"][0])
		if err != nil {
			return fmt.Errorf("Parsing start time (%v)", err)
		}
	}
	end := time.Now().AddDate(0, 1, 0)
	if len(formEnd) > 0 && formEnd[0] != "" {
		end, err = time.Parse("2006-01-02", r.PostForm["end"][0])
		if err != nil {
			return fmt.Errorf("Parsing end time (%v)", err)
		}
	}

	// get all buckets
	allBuckets, err := ledger.GetBuckets(tx)
	if err != nil {
		return fmt.Errorf("Calling ledger.GetBuckets() (%v)", err)
	}
	if len(formBuckets) == 0 {
		formBuckets = allBuckets
	}

	summary, err := ledger.SummarizeLedgerOverTime(tx, formBuckets, start, end)
	if err != nil {
		return fmt.Errorf("Calling ledger.SummarizeLedgerOverTime (%v)", err)
	}
	plot := ledger.MakePlot(summary, start)

	data := struct {
		AllBuckets []string
		Plot       ledger.PlotData
	}{
		allBuckets,
		*plot,
	}

	// execute template
	if err = t.Execute(w, data); err != nil {
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
