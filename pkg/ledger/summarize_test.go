package ledger_test

import (
	"database/sql"
	"ledger/pkg/ledger"
	"ledger/pkg/ledgerbucket"
	"ledger/pkg/testutils"
	"reflect"
	"testing"
	"time"
)

func TestSummarizeBalanceOverTime(t *testing.T) {
	db := testutils.Db(t)
	t.Run("three buckets over three days, including a zero value bucket", func(t *testing.T) {
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

func TestSummarizeLedger(t *testing.T) {
	db := testutils.Db(t)
	t.Run("simplest case",
		func(t *testing.T) {
			// Given
			throughdate := time.Now()
			input := ledgerbucket.Bucket{"checking", 1, "full"}

			testutils.Tx(t, db, func(tx *sql.Tx) error {
				if err := ledgerbucket.InsertBucket(tx, input); err != nil {
					return err
				}
				return nil
			})

			want := map[string]int{"checking": 0}

			// When
			var got map[string]int
			testutils.Tx(t, db, func(tx *sql.Tx) (err error) {
				got, err = ledger.SummarizeLedger(
					tx,
					throughdate,
				)
				return err
			})
			assertEqual(t, want, got)
		})
}

// 	earlyDate := time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local)
// 	laterDate := time.Date(2021, 4, 1, 0, 0, 0, 0, time.Local)
// 	entries := []Entry{
// 		{
// 			Source:      "savings",
// 			Destination: "checking",
// 			EntryDate:   earlyDate,
// 			Amount:      1000,
// 		},
// 		{
// 			Source:      "checking",
// 			Destination: "credit card",
// 			EntryDate:   earlyDate,
// 			Amount:      200,
// 		},
// 		{
// 			Source:      "checking",
// 			Destination: "credit card",
// 			EntryDate:   laterDate,
// 			Amount:      100,
// 		},
// 	}
// 	buckets := []ledgerbucket.Bucket{
// 		{
// 			Name:      "savings",
// 			Asset:     1,
// 			Liquidity: "full",
// 		},
// 		{
// 			Name:      "checking",
// 			Asset:     1,
// 			Liquidity: "full",
// 		},
// 		{
// 			Name:      "credit card",
// 			Asset:     0,
// 			Liquidity: "",
// 		},
// 		{
// 			Name:      "paycheck",
// 			Asset:     0,
// 			Liquidity: "",
// 		},
// 	}
// 	tx := testtx(t, db)
// 	for _, e := range entries {
// 		err := InsertEntry(tx, e)
// 		assertNoError(t, err, "inserting entries")
// 	}
// 	for _, b := range buckets {
// 		err := ledgerbucket.InsertBucket(tx, b)
// 		assertNoError(t, err, "classifying buckets")
// 	}
// 	testcommit(t, tx)
//
// 	// When
// 	tx = testtx(t, db)
// 	result, err := SummarizeLedger(tx, earlyDate)
// 	assertNoError(t, err, "summarizing all buckets through date")
// 	testcommit(t, tx)
// 	want := []balanceDetail{
// 		{
// 			bucket:    "checking",
// 			amount:    1000 - 200,
// 			asset:     1,
// 			liquidity: "full",
// 		},
// 		{
// 			bucket:    "savings",
// 			amount:    -1000,
// 			asset:     1,
// 			liquidity: "full",
// 		},
// 		{
// 			bucket:    "credit card",
// 			amount:    200,
// 			asset:     0,
// 			liquidity: "",
// 		},
// 	}
//
// 	// Then
// 	assertEqual(t, want, result, "")
// }

func assertEqual(t *testing.T, want, got interface{}) {
	t.Helper()
	if b := reflect.DeepEqual(want, got); !b {
		t.Fatalf("want: %v, got: %v", want, got)
	}
}
