package ledger_test

import (
	"database/sql"
	"fmt"
	"ledger/pkg/ledger"
	"ledger/pkg/testutils"
	"testing"
	"time"
)

func TestGetLedger(t *testing.T) {
	db := testutils.Db(t)
	t.Run("one entry",
		func(t *testing.T) {
			start := time.Date(2004, 8, 16, 0, 0, 0, 0, time.Local)
			end := start.AddDate(0, 0, 1)
			input := ledger.Entry{
				"savings",
				"checking",
				start,
				100,
			}
			testutils.Tx(t, db, func(tx *sql.Tx) error {
				err := ledger.InsertEntry(tx, input)
				return err
			})
			want := []ledger.Entry{input}
			var got []ledger.Entry
			testutils.Tx(t, db, func(tx *sql.Tx) (err error) {
				got, err = ledger.GetLedger(tx, start, end)
				return err
			})
			assertEqual(t, want, got)
		})
}

func TestSummarizeBalance(t *testing.T) {
	db := testutils.Db(t)
	bigBang := testutils.BigBang()
	t.Run("one entry,",
		func(t *testing.T) {
			entryDate := time.Now()
			entry := ledger.Entry{
				Source:      "savings",
				Destination: "checking",
				Amount:      100,
				EntryDate:   entryDate,
			}
			// insert entry
			testutils.Tx(t, db, func(tx *sql.Tx) error {
				err := ledger.InsertEntry(tx, entry)
				return err
			})
			// summarize from begining of time
			{
				want := map[string]int{"savings": -100, "checking": 100}
				var got map[string]int
				testutils.Tx(t, db, func(tx *sql.Tx) (err error) {
					got, err = ledger.SummarizeBalance(
						tx,
						[]string{"savings", "checking"},
						bigBang,
						entryDate.AddDate(0, 0, 1),
					)
					return err
				})
				assertEqual(t, want, got)
			}
			// summarize before entryDate
			{
				want := map[string]int{"savings": 0, "checking": 0}
				var got map[string]int
				testutils.Tx(t, db, func(tx *sql.Tx) (err error) {
					got, err = ledger.SummarizeBalance(
						tx,
						[]string{"savings", "checking"},
						bigBang,
						entryDate.AddDate(0, 0, -1),
					)
					return err
				})
				assertEqual(t, want, got)
			}

		})
}

func TestGetBalanceOverTime(t *testing.T) {
	db := testutils.Db(t)
	t.Run("three buckets over three days, including a zero value bucket", func(t *testing.T) {
		bucket1 := "our source bucket"
		bucket2 := "our destination bucket"
		bucket3 := "our bucket with zero entries"
		start := time.Now()

		input := []ledger.Entry{
			{bucket1, bucket2, start, 100},
			{bucket1, bucket2, start.AddDate(0, 0, 1), 100},
			{bucket1, bucket2, start.AddDate(0, 0, 2), 100},
		}

		testutils.Tx(t, db, func(tx *sql.Tx) error {
			for _, i := range input {
				if err := ledger.InsertEntry(tx, i); err != nil {
					return err
				}
			}
			return nil
		})
		want := []map[string]int{
			{bucket1: -100, bucket2: 100, bucket3: 0},
			{bucket1: -200, bucket2: 200, bucket3: 0},
			{bucket1: -300, bucket2: 300, bucket3: 0},
		}
		var got []map[string]int
		testutils.Tx(t, db, func(tx *sql.Tx) (err error) {
			got, err = ledger.GetBalanceOverTime(
				tx,
				[]string{bucket1, bucket2, bucket3},
				start,
				start.AddDate(0, 0, 3),
			)
			return err
		})
		assertEqual(t, want, got)
	})
}

func TestSummarizeLedgerOverTime(t *testing.T) {
	db := testutils.Db(t)
	t.Run("empty summary (zero entries)",
		func(t *testing.T) {
			today := time.Now()
			want := []map[string]int{}
			var got []map[string]int
			testutils.Tx(t, db, func(tx *sql.Tx) (err error) {
				got, err = ledger.SummarizeLedgerOverTime(
					tx,
					[]string{},
					today,
					today,
					1, // interval: summarize daily
				)
				return err
			})
			assertEqual(t, want, got)
		})

	t.Run("two entries over three days",
		func(t *testing.T) {
			start := time.Now()
			end := start.AddDate(0, 0, 3)
			entry := ledger.Entry{
				Source:      "savings",
				Destination: "checking",
				EntryDate:   start,
				Amount:      100,
			}
			// insert entries
			testutils.Tx(t, db, func(tx *sql.Tx) (err error) {
				fmt.Println("END:", end)
				for i := 0; i < 2; i++ {
					fmt.Println("ENTRY DATE:", entry.EntryDate)
					err := ledger.InsertEntry(tx, entry)
					if err != nil {
						return err
					}
					entry.EntryDate = entry.EntryDate.AddDate(0, 0, 1)
				}
				return err
			})
			// summarize transactions daily
			interval := 1 // group daily
			want := []map[string]int{
				{"checking": 100, "savings": -100},
				{"checking": 100, "savings": -100},
				{"checking": 0, "savings": 0},
			}
			var got []map[string]int
			testutils.Tx(t, db, func(tx *sql.Tx) (err error) {
				got, err = ledger.SummarizeLedgerOverTime(
					tx,
					[]string{"checking", "savings"},
					start,
					end,
					interval,
				)
				return err
			})
			assertEqual(t, want, got)
		})
}

func TestMakePlot(t *testing.T) {
	t.Run("empty summary (zero entries)", func(t *testing.T) {
		summary := []map[string]int{}
		start := time.Now()
		want := &ledger.PlotData{}
		got := ledger.MakePlot(summary, start)
		assertEqual(t, want, got)
	})
	t.Run("one entry", func(t *testing.T) {
		summary := []map[string]int{{"savings": -100, "checking": 100}}
		start := time.Now()
		startString := start.Format("2006-01-02")
		want := &ledger.PlotData{
			[]string{"checking", "savings"},
			[]string{startString},
			[][]int{{100, -100}},
		}

		got := ledger.MakePlot(summary, start)
		assertEqual(t, want, got)
	})
	t.Run("two entries over two days", func(t *testing.T) {
		summary := []map[string]int{{"savings": -100, "IRA": 0, "checking": 100}, {"savings": -100, "IRA": 50, "checking": 50}}
		start := time.Now()
		startString := start.Format("2006-01-02")
		tomorrowString := start.AddDate(0, 0, 1).Format("2006-01-02")
		want := &ledger.PlotData{
			[]string{"IRA", "checking", "savings"},
			[]string{startString, tomorrowString},
			[][]int{{0, 100, -100}, {50, 50, -100}},
		}
		got := ledger.MakePlot(summary, start)
		assertEqual(t, want, got)
	})
}

func TestGetBuckets(t *testing.T) {
	db := testutils.Db(t)
	t.Run("one transaction, two buckets", func(t *testing.T) {
		input := ledger.Entry{
			"savings",
			"checking",
			time.Date(2004, 8, 16, 0, 0, 0, 0, time.Local),
			100,
		}
		testutils.Tx(t, db, func(tx *sql.Tx) error {
			return ledger.InsertEntry(tx, input)
		})
		want := []string{"checking", "savings"}
		var got []string
		testutils.Tx(t, db, func(tx *sql.Tx) (err error) {
			got, err = ledger.GetBuckets(tx)
			return err
		})
		assertEqual(t, want, got)
	})
}
