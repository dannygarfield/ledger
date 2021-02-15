package myhttp

import (
	"database/sql"
	"fmt"
	"ledger/pkg/budget"
	"ledger/pkg/utils"
	"net/url"
	"strconv"
	"time"
)

func SetStartDate(tx *sql.Tx, values url.Values) (time.Time, error) {
	formStart := values.Get("startDate")
	if formStart != "" && formStart != "undefined" {
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
	formEnd := values.Get("endDate")
	if formEnd != "" && formEnd != "undefined" {
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

func SetTimeInterval(values url.Values) (int, error) {
	formInterval := values.Get("interval")
	if formInterval == "" {
		return 1, nil
	}
	interval, err := strconv.Atoi(formInterval)
	if err != nil {
		return -1, fmt.Errorf("Could not convert form interval %s to number: %v", formInterval, err)
	}
	return interval, nil
}

func SetBudgetCategories(tx *sql.Tx, values url.Values) ([]string, []string, error) {
	allCategories, err := budget.GetCategories(tx)
	if err != nil {
		return nil, nil, fmt.Errorf("Calling budget.GetCategories (%v)", err)
	}
	if values.Get("categories") != "" {
		formCategories := values["categories"]
		// fmt.Println("formCategories:", formCategories)
		// fmt.Println("values.Get():", values.Get("categories"))
		return formCategories, allCategories, nil
	}
	return allCategories, allCategories, nil
}
