package ledger_test

import (
	"database/sql"
	"ledger/pkg/ledger"
	"ledger/pkg/testutils"
	"reflect"
	"testing"
	"time"
)

func TestGetLedger(t *testing.T) {
	db := testutils.Db(t)
	t.Run("one entry",
		func(t *testing.T) {
			start := time.Date(1992, 8, 16, 0, 0, 0, 0, time.Local)
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

func TestSummarizeLedger(t *testing.T) {
	db := testutils.Db(t)
	t.Run("one entry",
		func(t *testing.T) {
			// GIVEN
			entryDate := time.Date(1992, 8, 16, 0, 0, 0, 0, time.Local)
			sourceBucket := "savings"
			destBucket := "checking"

			inputEntry := ledger.Entry{sourceBucket, destBucket, entryDate, 100}

			// insert entry
			testutils.Tx(t, db, func(tx *sql.Tx) error {
				err := ledger.InsertEntry(tx, inputEntry)
				return err
			})

			want := map[string]int{"savings": -100, "checking": 100}

			// When
			var got map[string]int
			testutils.Tx(t, db, func(tx *sql.Tx) (err error) {
				got, err = ledger.SummarizeLedger(
					tx,
					[]string{sourceBucket, destBucket},
					entryDate.AddDate(0, 0, 1),
				)
				return err
			})

			// Then
			assertEqual(t, want, got)
		})
}

func TestSummarizeBalanceOverTime(t *testing.T) {
	db := testutils.Db(t)
	t.Run("three buckets over three days, including a zero value bucket", func(t *testing.T) {
		bucket1 := "our source bucket"
		bucket2 := "our destination bucket"
		bucket3 := "our bucket with zero entries"
		start := time.Date(1992, 8, 16, 0, 0, 0, 0, time.Local)

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
			got, err = ledger.SummarizeLedgerOverTime(
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

func TestGetBuckets(t *testing.T) {
	db := testutils.Db(t)
	t.Run("one transaction, two buckets", func(t *testing.T) {
		input := ledger.Entry{
			"savings",
			"checking",
			time.Date(1992, 8, 16, 0, 0, 0, 0, time.Local),
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

func TestMakePlot(t *testing.T) {
	t.Run("empty summary (zero transactions)", func(t *testing.T) {
		summary := []map[string]int{}
		start := time.Now()
		want := &ledger.PlotData{}
		got := ledger.MakePlot(summary, start)
		assertEqual(t, want, got)
	})

	t.Run("one transaction", func(t *testing.T) {
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

	t.Run("two transactions over two days", func(t *testing.T) {
		summary := []map[string]int{{"savings": -100, "checking": 100}, {"savings": -200, "checking": 200}}
		start := time.Now()
		startString := start.Format("2006-01-02")
		tomorrowString := start.AddDate(0, 0, 1).Format("2006-01-02")
		want := &ledger.PlotData{
			[]string{"checking", "savings"},
			[]string{startString, tomorrowString},
			[][]int{{100, -100}, {200, -200}},
		}

		got := ledger.MakePlot(summary, start)
		assertEqual(t, want, got)
	})
}

func assertEqual(t *testing.T, want, got interface{}) {
	t.Helper()
	if b := reflect.DeepEqual(want, got); !b {
		t.Fatalf("want: %v, got: %v", want, got)
	}
}
