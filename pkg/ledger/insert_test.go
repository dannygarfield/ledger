package ledger_test

import (
	"database/sql"
	"io/ioutil"
	"ledger/pkg/ledger"
	"ledger/pkg/testutils"
	"testing"
	"time"
)

func TestInsertEntry(t *testing.T) {
	// initialize db and test variables
	db := testutils.Db(t)
	bigBang := testutils.BigBang()
	entryDate := time.Now()
	t.Run("one entry",
		func(t *testing.T) {
			entry := ledger.Entry{
				Source:      "savings",
				Destination: "checking",
				EntryDate:   entryDate,
				Amount:      100,
			}
			// insert entry
			testutils.Tx(t, db, func(tx *sql.Tx) error {
				err := ledger.InsertEntry(tx, entry)
				return err
			})
			// summarize ledger source
			{
				want := -100
				var got int
				testutils.Tx(t, db, func(tx *sql.Tx) (err error) {
					got, err = ledger.SummarizeBucket(tx, "savings", bigBang, entryDate)
					return err
				})
				testutils.AssertEqual(t, want, got)
			}
			// summarize ledger destination
			{
				want := 100
				var got int
				testutils.Tx(t, db, func(tx *sql.Tx) (err error) {
					got, err = ledger.SummarizeBucket(tx, "checking", bigBang, entryDate)
					return err
				})
				testutils.AssertEqual(t, want, got)
			}
		})
}

func TestInsertRepeatingEntry(t *testing.T) {
	// initialize db and test vars
	db := testdb(t)
	bigBang := testutils.BigBang()
	entryDate := time.Now()
	t.Run("one monthly repeating entry",
		func(t *testing.T) {
			entry := ledger.Entry{
				Source:      "savings",
				Destination: "checking",
				Amount:      100,
				EntryDate:   entryDate,
			}
			// insert repeating entry
			testutils.Tx(t, db, func(tx *sql.Tx) error {
				err := ledger.InsertRepeatingEntry(tx, entry, "monthly")
				return err
			})
			// summarize ledger source
			want := -100 * 25
			var got int
			testutils.Tx(t, db, func(tx *sql.Tx) (err error) {
				got, err = ledger.SummarizeBucket(tx, "savings", bigBang, entryDate.AddDate(2, 0, 0))
				return err
			})
			testutils.AssertEqual(t, want, got)
		})
}

// helper functions
func testdb(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("opening database: %v", err)
	}
	schema, err := ioutil.ReadFile("../../schema.sql")
	if err != nil {
		t.Fatalf("opening schema file: %v", err)
	}
	if _, err := db.Exec(string(schema)); err != nil {
		t.Fatalf("loading schema: %v", err)
	}
	return db
}

func testtx(t *testing.T, db *sql.DB) *sql.Tx {
	t.Helper()

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("beginning sql transactions: %v", err)
	}
	return tx
}

func testcommit(t *testing.T, tx *sql.Tx) {
	t.Helper()

	if err := tx.Commit(); err != nil {
		t.Fatalf("committing sql transaction: %v", err)
	}
}

func assertNoError(t *testing.T, err error, msg string) {
	t.Helper()
	if err != nil {
		t.Fatalf("%s: %v", msg, err)
	}
}
