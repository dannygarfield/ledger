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
		bucket := "our bucket"
		start := time.Date(1992, 8, 16, 0, 0, 0, 0, time.UTC)
		input := []struct {
			bucket string
			date   time.Time
			val    int
		}{
			{bucket, start, 100},
			{bucket, start.AddDate(0, 0, 1), 100},
			{bucket, start.AddDate(0, 0, 2), 100},
		}
		testutils.Tx(t, db, func(tx *sql.Tx) error {
			for _, i := range input {
				if err := ledger.InsertEntry(tx, ledger.Entry{
					Source:      "whatever",
					Destination: i.bucket,
					EntryDate:   i.date,
					Amount:      i.val,
				}); err != nil {
					return err
				}
			}
			return nil
		})
		want := []map[string]int{
			{bucket: 100},
			{bucket: 200},
			{bucket: 300},
		}

		var got []map[string]int
		testutils.Tx(t, db, func(tx *sql.Tx) (err error) {
			got, err = ledger.MakePlot(
				tx,
				[]string{bucket},
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
