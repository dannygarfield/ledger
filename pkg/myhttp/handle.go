package myhttp

import (
	"database/sql"
	"fmt"
	"ledger/pkg/budget"
	"ledger/pkg/mytemplate"
	"net/http"
	"time"
)

func HandleBudgetOverTime(tx *sql.Tx, r *http.Request, w http.ResponseWriter) error {
	r.ParseForm()
	// set start date
	startDate, err := SetStartDate(tx, r.Form)
	if err != nil {
		return fmt.Errorf("Calling SetStartDate: %v", err)
	}
	// set end date
	endDate, err := SetEndDate(tx, r.Form)
	if err != nil {
		return fmt.Errorf("Calling SetEndDate: %v", err)
	}
	// set interval
	timeInterval, err := SetTimeInterval(r.Form)
	if err != nil {
		return fmt.Errorf("Calling SetTimeInterval: %v", err)
	}
	// set categories
	filterCategories, allCategories, err := SetBudgetCategories(tx, r.Form)
	if err != nil {
		return fmt.Errorf("Calling SetBudgetCategories: %v", err)
	}
	// get summary of spending over time
	spendSummary, err := budget.SummarizeSpendsOverTime(tx, filterCategories, startDate, endDate, timeInterval)
	if err != nil {
		return fmt.Errorf("Calling ledger.SummarizeBalanceOverTime (%v)", err)
	}
	// plot
	plot := budget.MakePlot(spendSummary, startDate, timeInterval)
	// construct data for html template
	htmlTemplateData := struct {
		Start, End    time.Time
		TimeInterval  int
		AllCategories []string
		Plot          budget.PlotData
	}{
		startDate,
		endDate,
		timeInterval,
		allCategories,
		*plot,
	}
	// call template function
	err = mytemplate.BudgetOverTime(w, htmlTemplateData)
	if err != nil {
		return fmt.Errorf("Could not call mytemplate.BudgetOverTime: %v", err)
	}
	return nil
}
