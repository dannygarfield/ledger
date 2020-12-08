package dateseries

import (
	"database/sql"
	"io/ioutil"
	"ledger/pkg/ledger"
	"reflect"
	"testing"
	"time"
)

func TestCreateSeries(t *testing.T) {
	// Given
	db := testdb(t)

	// When
	end := time.Now().AddDate(2, 0, 0)
	{
		tx := testtx(t, db)
		err := UpdateSeries(tx, end)
		assertNoError(t, err, "")
		testcommit(t, tx)
	}

	// Then
	tx := testtx(t, db)
	maxDate, err := GetMaxDate(tx)
	assertNoError(t, err, "test: getting max date")
	testcommit(t, tx)
	endFormatted := ledger.ConvertToDate(end)
	assertEqual(t, endFormatted, maxDate, "")
}

func TestUpdateSeries(t *testing.T) {
	// Given
	db := testdb(t)
	end1 := time.Now().AddDate(2, 0, 0)
	tx := testtx(t, db)
	err := UpdateSeries(tx, end1)
	assertNoError(t, err, "")
	testcommit(t, tx)

	// When
	newEnd := time.Now().AddDate(2, 0, 2)
	tx = testtx(t, db)
	err = UpdateSeries(tx, newEnd)
	assertNoError(t, err, "")
	testcommit(t, tx)

	// Then
	tx = testtx(t, db)
	maxDate, err := GetMaxDate(tx)
	assertNoError(t, err, "test: getting max date")
	testcommit(t, tx)
	newEndFormatted := ledger.ConvertToDate(newEnd)

	assertEqual(t, newEndFormatted, maxDate, "")

}

// Helper functions
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

func assertEqual(t *testing.T, want, got interface{}, msg string) {
	t.Helper()
	if b := reflect.DeepEqual(want, got); !b {
		t.Fatalf("%s: want: %v, got: %v", msg, want, got)
	}
}
