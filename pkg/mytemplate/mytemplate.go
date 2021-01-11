package mytemplate

import (
	"database/sql"
	"fmt"
	"html/template"
	"ledger/pkg/ledger"
	"net/http"
	"strconv"
	"time"
)

// display a ledger on a single day
func Ledger(tx *sql.Tx, w http.ResponseWriter, r *http.Request) error {
	// parse html template
	t, err := template.ParseFiles("pkg/mytemplate/ledger.html")
	if err != nil {
		return fmt.Errorf("Could not parse ledger.html (%v)", err)
	}
	// parse html form
	r.ParseForm()
	formStart := r.PostForm["start"]
	formEnd := r.PostForm["end"]
	// set start date
	start := time.Now().AddDate(0, -1, 0)
	if len(formStart) > 0 && formStart[0] != "" {
		start, err = time.Parse("2006-01-02", r.PostForm["start"][0])
		if err != nil {
			return fmt.Errorf("Parsing start time (%v)", err)
		}
	}
	// set end date
	end := time.Now().AddDate(0, 1, 0)
	if len(formEnd) > 0 && formEnd[0] != "" {
		end, err = time.Parse("2006-01-02", r.PostForm["end"][0])
		if err != nil {
			return fmt.Errorf("Parsing end time (%v)", err)
		}
	}
	// get ledger data
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

// display the ledger's net balances over time, daily
func BalanceOverTime(tx *sql.Tx, w http.ResponseWriter, r *http.Request) error {
	// parse html template
	t, err := template.ParseFiles("pkg/mytemplate/balance.html")
	if err != nil {
		return fmt.Errorf("Could not parse balance.html (%v)", err)
	}
	// parse html form
	r.ParseForm()
	formStart := r.PostForm["start"]
	formEnd := r.PostForm["end"]
	formBuckets := r.PostForm["buckets"]
	// set start date
	start := time.Now().AddDate(0, -1, 0)
	if len(formStart) > 0 && formStart[0] != "" {
		start, err = time.Parse("2006-01-02", r.PostForm["start"][0])
		if err != nil {
			return fmt.Errorf("Parsing start time (%v)", err)
		}
	}
	// set end date
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
	// if we don't get buckets from user input, show all buckets
	if len(formBuckets) == 0 {
		formBuckets = allBuckets
	}
	// get summary data and format for html
	summary, err := ledger.SummarizeBalanceOverTime(tx, formBuckets, start, end)
	if err != nil {
		return fmt.Errorf("Calling ledger.SummarizeBalanceOverTime (%v)", err)
	}
	plot := ledger.MakePlot(summary, start, 1)
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

// display the ledger over time, grouped into a given interval period
func LedgerOverTime(tx *sql.Tx, w http.ResponseWriter, r *http.Request) error {
	// parse html template
	t, err := template.ParseFiles("pkg/mytemplate/ledgerseries.html")
	if err != nil {
		return fmt.Errorf("Could not parse balance.html (%v)", err)
	}
	// parse html form
	r.ParseForm()
	formStart := r.PostForm["start"]
	formEnd := r.PostForm["end"]
	formBuckets := r.PostForm["buckets"]
	formInterval := r.PostForm["interval"]
	// set start date
	start := time.Now().AddDate(0, -1, 0)
	if len(formStart) > 0 && formStart[0] != "" {
		start, err = time.Parse("2006-01-02", r.PostForm["start"][0])
		if err != nil {
			return fmt.Errorf("Parsing start time (%v)", err)
		}
	}
	// set end date
	end := time.Now().AddDate(0, 1, 0)
	if len(formEnd) > 0 && formEnd[0] != "" {
		end, err = time.Parse("2006-01-02", r.PostForm["end"][0])
		if err != nil {
			return fmt.Errorf("Parsing end time (%v)", err)
		}
	}
	// set interval
	interval := 1
	if len(formInterval) > 0 && formInterval[0] != "" {
		interval, err = strconv.Atoi(formInterval[0])
		if err != nil {
			return fmt.Errorf("calling strconv.Atoi() (%v)", err)
		}
	}
	// get all buckets
	allBuckets, err := ledger.GetBuckets(tx)
	if err != nil {
		return fmt.Errorf("Calling ledger.GetBuckets() (%v)", err)
	}
	// if we don't get buckets from user input, show all buckets
	if len(formBuckets) == 0 {
		formBuckets = allBuckets
	}
	// get summary data and format for html
	summary, err := ledger.SummarizeLedgerOverTime(tx, formBuckets, start, end, interval)
	if err != nil {
		return fmt.Errorf("Calling ledger.SummarizeBalanceOverTime (%v)", err)
	}
	plot := ledger.MakePlot(summary, start, interval)
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
