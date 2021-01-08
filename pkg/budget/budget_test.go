package budget_test

import (
	"database/sql"
	"ledger/pkg/budget"
	"ledger/pkg/testutils"
	"ledger/pkg/utils"
	"testing"
	"time"
)

func TestInsertEntry(t *testing.T) {
	// initialize db and test variables
	db := testutils.Db(t)
	bigBang := testutils.BigBang()
	entryDate := utils.ConvertToDate(time.Now())
	t.Run("one entry",
		func(t *testing.T) {
			entry := budget.Entry{
				EntryDate:   entryDate,
				Amount:      100,
				Category:    "groceries",
				Description: "whole foods delivery",
			}
			// insert entry
			testutils.Tx(t, db, func(tx *sql.Tx) error {
				err := budget.InsertEntry(tx, entry)
				return err
			})
			// summarize
			want := []budget.Entry{entry}
			var got []budget.Entry
			testutils.Tx(t, db, func(tx *sql.Tx) (err error) {
				got, err = budget.GetBudgetEntries(tx, bigBang, entryDate.AddDate(0, 0, 1))
				return err
			})
			testutils.AssertEqual(t, want, got)
		})
}
