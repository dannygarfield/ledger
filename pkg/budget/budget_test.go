package budget_test

import (
	"database/sql"
	"ledger/pkg/budget"
	"ledger/pkg/testutils"
	"testing"
	"time"
)

func TestInsertEntry(t *testing.T) {
	// initialize db and test variables
	db := testutils.Db(t)
	bigBang := testutils.BigBang()
	entryDate := time.Date(2020, 12, 31, 0, 0, 0, 0, time.UTC)
	t.Run("one entry",
		func(t *testing.T) {
			entry := budget.Entry{
				EntryDate:   entryDate,
				Amount:      100,
				Category:    "groceries",
				Description: "NYE meal",
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
				got, err = budget.GetBudgetEntries(tx, bigBang, entryDate)
				return err
			})
			testutils.AssertEqual(t, want, got)
		})
}

func TestGetBudget(t *testing.T) {
	db := testutils.Db(t)
	janOne := testutils.JanOne()
	janTwo := testutils.JanTwo()
	t.Run("one entry",
		func(t *testing.T) {
			want := []budget.Entry{
				{
					EntryDate:   janTwo,
					Amount:      200,
					Category:    "groceries",
					Description: "food train",
				},
			}
			var got []budget.Entry
			testutils.Tx(t, db, func(tx *sql.Tx) (err error) {
				got, err = budget.GetBudgetEntries(tx, janTwo, janTwo)
				return err
			})
			testutils.AssertEqual(t, want, got)
		})
	t.Run("three entries",
		func(t *testing.T) {
			want := []budget.Entry{
				{
					EntryDate:   janOne,
					Amount:      3000,
					Category:    "rent",
					Description: "-",
				},
				{
					EntryDate:   janOne,
					Amount:      100,
					Category:    "groceries",
					Description: "whole foods delivery",
				},
				{
					EntryDate:   janTwo,
					Amount:      200,
					Category:    "groceries",
					Description: "food train",
				},
			}
			var got []budget.Entry
			testutils.Tx(t, db, func(tx *sql.Tx) (err error) {
				got, err = budget.GetBudgetEntries(tx, janOne, janTwo)
				return err
			})
			testutils.AssertEqual(t, want, got)
		})
}

func TestSummarizeCategory(t *testing.T) {
	db := testutils.Db(t)
	t.Run("empty summary (zero entries in time period)",
		func (t *testing.T) {
				want := 0
				var got int
				testutils.Tx(t, db, func(tx *sql.Tx) (err error) {
					got, err = budget.SummarizeCategory(
						tx,
						"groceries",
						testutils.BigBang(),
						testutils.BigBang())
					return err
				})
				testutils.AssertEqual(t, want, got)
		})
	t.Run("empty summary (zero entries in category)",
		func (t *testing.T) {
				want := 0
				var got int
				testutils.Tx(t, db, func(tx *sql.Tx) (err error) {
					got, err = budget.SummarizeCategory(
						tx,
						"gifts",
						testutils.JanOne(),
						testutils.JanTwo())
					return err
				})
				testutils.AssertEqual(t, want, got)
		})
	t.Run("one entry",
		func (t *testing.T) {
				want := 3000
				var got int
				testutils.Tx(t, db, func(tx *sql.Tx) (err error) {
					got, err = budget.SummarizeCategory(
						tx,
						"rent",
						testutils.JanOne(),
						testutils.JanTwo())
					return err
				})
				testutils.AssertEqual(t, want, got)
		})
	t.Run("two entries",
		func (t *testing.T) {
				want := 300
				var got int
				testutils.Tx(t, db, func(tx *sql.Tx) (err error) {
					got, err = budget.SummarizeCategory(
						tx,
						"groceries",
						testutils.JanOne(),
						testutils.JanTwo())
					return err
				})
				testutils.AssertEqual(t, want, got)
		})
}
