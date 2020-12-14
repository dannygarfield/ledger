package ledger

import (
	"database/sql"
	"io/ioutil"
	"reflect"
	"testing"
	"time"
)

func TestInsertEntry(t *testing.T) {
	// Given
	db := testdb(t)
	e := Entry{
		Source:      "checking",
		Destination: "credit card",
		EntryDate:   time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local),
		Amount:      120,
	}

	// When
	tx := testtx(t, db)
	err := InsertEntry(tx, e)
	assertNoError(t, err, "")
	testcommit(t, tx)

	// Then
	{
		tx := testtx(t, db)
		result, err := SummarizeBucket(tx, e.Source, e.EntryDate)
		assertNoError(t, err, "summary(source)")
		testcommit(t, tx)
		assertEqual(t, -e.Amount, result, "source")
	}
	{
		tx := testtx(t, db)
		result, err := SummarizeBucket(tx, e.Destination, e.EntryDate)
		assertNoError(t, err, "summary(destination)")
		testcommit(t, tx)
		assertEqual(t, e.Amount, result, "destination")
	}
}

func TestInsertRepeatingEntry(t *testing.T) {
	// Given
	db := testdb(t)
	e1 := Entry{
		Source:      "checking",
		Destination: "IRA",
		Amount:      50,
		EntryDate:   time.Now(), // repeating write until 2 years from now. setting EntryDate to time.Now() requires less math
	}
	e2 := Entry{
		Source:      "checking",
		Destination: "rent",
		Amount:      50,
		EntryDate:   time.Now(),
	}

	// When
	tx := testtx(t, db)
	{
		err := InsertRepeatingEntry(tx, e1, "weekly")
		assertNoError(t, err, "inserting weekly entry")
	}
	{
		err := InsertRepeatingEntry(tx, e2, "monthly")
		assertNoError(t, err, "inserting repeating entry")
	}
	testcommit(t, tx)

	// Then
	endDate := time.Now().AddDate(2, 0, 0)
	{
		tx := testtx(t, db)
		result, err := SummarizeBucket(tx, e1.Source, endDate)
		assertNoError(t, err, "")
		testcommit(t, tx)
		assertEqual(t, -e1.Amount*105-e2.Amount*25, result, "inserting weekly")
	}
	{
		tx := testtx(t, db)
		result, err := SummarizeBucket(tx, e2.Destination, endDate)
		assertNoError(t, err, "")
		testcommit(t, tx)
		assertEqual(t, e2.Amount*25, result, "inserting monthly")
	}

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

func assertEqual(t *testing.T, want, got interface{}, msg string) {
	t.Helper()
	if b := reflect.DeepEqual(want, got); !b {
		t.Fatalf("%s: want: %v, got: %v", msg, want, got)
	}
}
