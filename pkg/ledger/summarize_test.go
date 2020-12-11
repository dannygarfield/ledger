package ledger_test

import (
	"database/sql"
	"ledger/pkg/ledger"
	"ledger/pkg/testutils"
	"reflect"
	"testing"
	"time"
)

func TestSummarizeBalanceOverTime(t *testing.T) {
	db := testutils.Db(t)
	t.Run("basic", func(t *testing.T) {
		bucket1 := "our source bucket"
		bucket2 := "our destination bucket"
		bucket3 := "our bucket with zero entries"
		start := time.Date(1992, 8, 16, 0, 0, 0, 0, time.UTC)

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
			got, err = ledger.MakePlot(
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

func assertEqual(t *testing.T, want, got interface{}) {
	t.Helper()
	if b := reflect.DeepEqual(want, got); !b {
		t.Fatalf("want: %v, got: %v", want, got)
	}
}
