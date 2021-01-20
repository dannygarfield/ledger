package myhttp

import (
	"database/sql"
	"fmt"
	"ledger/pkg/budget"
	"ledger/pkg/utils"
	"net/http"
	"net/url"
	"time"
)

func HandleBudgetList(tx *sql.Tx, r *http.Request, w http.ResponseWriter) error {
	// set start date
	start, err := SetStartDate(tx, r.Form)
	if err != nil {
		http.Error(w, fmt.Sprintf("Calling myhttp.SetStartDate: %v", err), http.StatusInternalServerError)
		return err
	}
	// set end date
	end, err := SetEndDate(tx, r.Form)
	if err != nil {
		http.Error(w, fmt.Sprintf("Calling myhttp.SetStartDate: %v", err), http.StatusInternalServerError)
		return err
	}
	// get budget entries
	budgetEntries, err := budget.GetBudgetEntries(tx, start, end)
	if err != nil {
		return fmt.Errorf("Calling budget.GetBudgetEntries() (%v)", err)
	}
    fmt.Println(budgetEntries)
    return nil
}

func SetStartDate(tx *sql.Tx, values url.Values) (time.Time, error) {
	formStart := values.Get("start")
	if formStart != "" {
		start, err := time.Parse("2006-01-02", formStart)
		if err != nil {
			return utils.BigBang, fmt.Errorf("Parsing time (%v)", err)
		}
		return start, nil
	} else {
		start, err := budget.GetEarliestBudgetDate(tx)
		if err != nil {
			return utils.BigBang, fmt.Errorf("Calling budget.GetEarliestBudgetDate() (%v)", err)
		}
		return start, nil

	}
}

func SetEndDate(tx *sql.Tx, values url.Values) (time.Time, error) {
	formEnd := values.Get("end")
	if formEnd != "" {
		end, err := time.Parse("2006-01-02", formEnd)
		if err != nil {
			return utils.BigBang, fmt.Errorf("Parsing time (%v)", err)
		}
		return end, nil
	} else {
		end, err := budget.GetLatestBudgetDate(tx)
		if err != nil {
			return utils.BigBang, fmt.Errorf("Calling budget.GetLatestBudgetDate() (%v)", err)
		}
		return end, nil

	}
}
