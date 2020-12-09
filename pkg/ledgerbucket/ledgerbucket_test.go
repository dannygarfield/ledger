package ledgerbucket

import (
	"database/sql"
	"io/ioutil"
	"reflect"
	"testing"
)

func TestInsertBucket(t *testing.T) {
	// Given
	db := testdb(t)
	b := Bucket{
		Name:      "Checking",
		Asset:     1,
		Liquidity: "Full",
	}

	// When
	tx := testtx(t, db)
	err := InsertBucket(tx, b)
	assertNoError(t, err, "")
	testcommit(t, tx)

	// Then
	tx = testtx(t, db)
	result, err := ShowBuckets(tx)
	assertNoError(t, err, "")
	testcommit(t, tx)
	bs := []Bucket{b}
	assertEqual(t, bs, result, "")
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
