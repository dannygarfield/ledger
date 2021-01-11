package budget_test

import (
	"database/sql"
	"ledger/pkg/budget"
	"ledger/pkg/testutils"
	"testing"
)

func TestInsertEntry(t *testing.T) {
	// initialize db and test variables
	db := testutils.Db(t)
	bigBang := testutils.BigBang()
	dec31 := testutils.Dec31()
	t.Run("one entry",
		func(t *testing.T) {
			entry := budget.Entry{
				EntryDate:   dec31,
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
				got, err = budget.GetBudgetEntries(tx, bigBang, dec31)
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
	t.Run("empty summary, no category given",
		func(t *testing.T) {
			want := 0
			var got int
			testutils.Tx(t, db, func(tx *sql.Tx) (err error) {
				got, err = budget.SummarizeCategory(
					tx,
					"",
					testutils.BigBang(),
					testutils.BigBang())
				return err
			})
			testutils.AssertEqual(t, want, got)
		})
	t.Run("empty summary, zero entries in time period",
		func(t *testing.T) {
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
	t.Run("empty summary, zero entries in category",
		func(t *testing.T) {
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
		func(t *testing.T) {
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
	t.Run("two entries over two days",
		func(t *testing.T) {
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

func TestSummarizeCategories(t *testing.T) {
	db := testutils.Db(t)
	dec31 := testutils.Dec31()
	janOne := testutils.JanOne()
	janTwo := testutils.JanTwo()
	t.Run("empty, no categories given",
		func(t *testing.T) {
			want := map[string]int{}
			var got map[string]int
			testutils.Tx(t, db, func(tx *sql.Tx) (err error) {
				got, err = budget.SummarizeCategories(tx, []string{}, janOne, janTwo)
				return err
			})
			testutils.AssertEqual(t, want, got)
	})
	t.Run("one category with no entries",
		func(t *testing.T) {
			want := map[string]int{"home": 0}
			var got map[string]int
			testutils.Tx(t, db, func(tx *sql.Tx) (err error) {
				got, err = budget.SummarizeCategories(
					tx,
					[]string{"home"},
					janOne,
					janTwo)
				return err
			})
			testutils.AssertEqual(t, want, got)
	})
	t.Run("one category with no entries in time period",
		func(t *testing.T) {
			want := map[string]int{"groceries": 0}
			var got map[string]int
			testutils.Tx(t, db, func(tx *sql.Tx) (err error) {
				got, err = budget.SummarizeCategories(
					tx,
					[]string{"groceries"},
					dec31,
					dec31)
				return err
			})
			testutils.AssertEqual(t, want, got)
	})
	t.Run("two categories with no entries",
		func(t *testing.T) {
			want := map[string]int{"home": 0, "utilities": 0}
			var got map[string]int
			testutils.Tx(t, db, func(tx *sql.Tx) (err error) {
				got, err = budget.SummarizeCategories(
					tx,
					[]string{"home", "utilities"},
					janOne,
					janTwo)
				return err
			})
			testutils.AssertEqual(t, want, got)
	})
	t.Run("one category with one entry on one day",
		func(t *testing.T) {
			want := map[string]int{"groceries": 100}
			var got map[string]int
			testutils.Tx(t, db, func(tx *sql.Tx) (err error) {
				got, err = budget.SummarizeCategories(
					tx,
					[]string{"groceries"},
					janOne,
					janOne)
				return err
			})
			testutils.AssertEqual(t, want, got)
	})
	t.Run("one category with two entries over two days",
		func(t *testing.T) {
			want := map[string]int{"groceries": 300}
			var got map[string]int
			testutils.Tx(t, db, func(tx *sql.Tx) (err error) {
				got, err = budget.SummarizeCategories(
					tx,
					[]string{"groceries"},
					janOne,
					janTwo)
				return err
			})
			testutils.AssertEqual(t, want, got)
	})
	t.Run("multiple entries over multiple days",
		func(t *testing.T) {
			want := map[string]int{"groceries": 300, "rent": 3000}
			var got map[string]int
			testutils.Tx(t, db, func(tx *sql.Tx) (err error) {
				got, err = budget.SummarizeCategories(
					tx,
					[]string{"groceries", "rent"},
					janOne,
					janTwo)
				return err
			})
			testutils.AssertEqual(t, want, got)
	})
	t.Run("three categories",
		func(t *testing.T) {
			want := map[string]int{"groceries": 200, "home": 0, "rent": 0}
			var got map[string]int
			testutils.Tx(t, db, func(tx *sql.Tx) (err error) {
				got, err = budget.SummarizeCategories(
					tx,
					[]string{"groceries", "rent", "home"},
					janTwo,
					janTwo)
				return err
			})
			testutils.AssertEqual(t, want, got)
	})
}

func TestSummarizeSpendsOverTime(t *testing.T) {
	db := testutils.Db(t)
	janOne := testutils.JanOne()
	janTwo := testutils.JanTwo()
	t.Run("empty, no categories given",
		func(t *testing.T) {
			want := []map[string]int{{}}
			var got []map[string]int
			testutils.Tx(t, db, func(tx *sql.Tx) (err error) {
				got, err = budget.SummarizeSpendsOverTime(
					tx,
					[]string{},
					janOne,
					janOne,
					1, // interval: summarize daily
				)
				return err
			})
			testutils.AssertEqual(t, want, got)
	})
	t.Run("empty, one category given, over one day",
		func(t *testing.T) {
			want := []map[string]int{
				{"COBRA": 0},
			}
			var got []map[string]int
			testutils.Tx(t, db, func(tx *sql.Tx) (err error) {
				got, err = budget.SummarizeSpendsOverTime(
					tx,
					[]string{"COBRA"},
					janOne,
					janOne,
					1,
				)
				return err
			})
			testutils.AssertEqual(t, want, got)
	})
	t.Run("empty, two categories given, over one day",
		func(t *testing.T) {
			want := []map[string]int{
				{"COBRA": 0, "laundry": 0},
			}
			var got []map[string]int
			testutils.Tx(t, db, func(tx *sql.Tx) (err error) {
				got, err = budget.SummarizeSpendsOverTime(
					tx,
					[]string{"COBRA", "laundry"},
					janOne,
					janOne,
					1,
				)
				return err
			})
			testutils.AssertEqual(t, want, got)
	})
	t.Run("empty, two categories given, over two days",
		func(t *testing.T) {
			want := []map[string]int{
				{"COBRA": 0, "laundry": 0},
				{"COBRA": 0, "laundry": 0},
			}
			var got []map[string]int
			testutils.Tx(t, db, func(tx *sql.Tx) (err error) {
				got, err = budget.SummarizeSpendsOverTime(
					tx,
					[]string{"COBRA", "laundry"},
					janOne,
					janTwo,
					1,
				)
				return err
			})
			testutils.AssertEqual(t, want, got)
	})
	t.Run("empty, two categories given, over two days",
		func(t *testing.T) {
			want := []map[string]int{
				{"COBRA": 0, "laundry": 0},
				{"COBRA": 0, "laundry": 0},
			}
			var got []map[string]int
			testutils.Tx(t, db, func(tx *sql.Tx) (err error) {
				got, err = budget.SummarizeSpendsOverTime(
					tx,
					[]string{"COBRA", "laundry"},
					janOne,
					janTwo,
					1,
				)
				return err
			})
			testutils.AssertEqual(t, want, got)
	})
	t.Run("one category given, over one day",
		func(t *testing.T) {
			want := []map[string]int{
				{"groceries": 100},
			}
			var got []map[string]int
			testutils.Tx(t, db, func(tx *sql.Tx) (err error) {
				got, err = budget.SummarizeSpendsOverTime(
					tx,
					[]string{"groceries"},
					janOne,
					janOne,
					1,
				)
				return err
			})
			testutils.AssertEqual(t, want, got)
	})
	t.Run("two categories given, over two days",
		func(t *testing.T) {
			want := []map[string]int{
				{"groceries": 100, "rent": 3000},
				{"groceries": 200, "rent": 0},
			}
			var got []map[string]int
			testutils.Tx(t, db, func(tx *sql.Tx) (err error) {
				got, err = budget.SummarizeSpendsOverTime(
					tx,
					[]string{"groceries", "rent"},
					janOne,
					janTwo,
					1,
				)
				return err
			})
			testutils.AssertEqual(t, want, got)
	})
	t.Run("one category over two days, skip by two",
		func(t *testing.T) {
			want := []map[string]int{
				{"groceries": 300},
			}
			var got []map[string]int
			testutils.Tx(t, db, func(tx *sql.Tx) (err error) {
				got, err = budget.SummarizeSpendsOverTime(
					tx,
					[]string{"groceries"},
					janOne,
					janTwo,
					2,
				)
				return err
			})
			testutils.AssertEqual(t, want, got)
	})
	t.Run("two categories over two days, skip by two",
		func(t *testing.T) {
			want := []map[string]int{
				{"groceries": 300, "rent": 3000},
			}
			var got []map[string]int
			testutils.Tx(t, db, func(tx *sql.Tx) (err error) {
				got, err = budget.SummarizeSpendsOverTime(
					tx,
					[]string{"rent", "groceries"},
					janOne,
					janTwo,
					2,
				)
				return err
			})
			testutils.AssertEqual(t, want, got)
	})
	t.Run("two categories over five days, skip by two",
		func(t *testing.T) {
			want := []map[string]int{
				{"groceries": 300, "rent": 3000},
				{"groceries": 0, "rent": 0},
				{"groceries": 0, "rent": 0},
			}
			var got []map[string]int
			testutils.Tx(t, db, func(tx *sql.Tx) (err error) {
				got, err = budget.SummarizeSpendsOverTime(
					tx,
					[]string{"rent", "groceries"},
					janOne,
					janOne.AddDate(0, 0, 4),
					2,
				)
				return err
			})
			testutils.AssertEqual(t, want, got)
	})

}
