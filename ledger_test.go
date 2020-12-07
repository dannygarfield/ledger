package main

import (
	"database/sql"
	"io/ioutil"
	"reflect"
	"testing"
	"time"
)

func TestInsert(t *testing.T) {
	// Given
	db := testdb(t)
	e := entry{
		source:      "checking",
		destination: "credit card",
		happenedAt:  time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local),
		amount:      120,
	}

	// When
	tx := testtx(t, db)
	err := insert(tx, e)
	assertNoError(t, err, "")
	testcommit(t, tx)

	// Then
	{
		tx := testtx(t, db)
		result, err := summarizeBucket(tx, e.source, e.happenedAt)
		assertNoError(t, err, "summary(source)")
		testcommit(t, tx)
		assertEqual(t, -e.amount, result, "source")
	}
	{
		tx := testtx(t, db)
		result, err := summarizeBucket(tx, e.destination, e.happenedAt)
		assertNoError(t, err, "summary(destination)")
		testcommit(t, tx)
		assertEqual(t, e.amount, result, "destination")
	}
}

func TestInsertRepeatingEntry(t *testing.T) {
	// Given
	db := testdb(t)
	e1 := entry{
		source:      "checking",
		destination: "IRA",
		amount:      50,
		happenedAt:  time.Now(), // repeating write until 2 years from now. setting happenedAt to time.Now() requires less math
	}
	e2 := entry{
		source:      "checking",
		destination: "rent",
		amount:      50,
		happenedAt:  time.Now(),
	}

	// When
	tx := testtx(t, db)
	{
		err := insertRepeating(tx, e1, "weekly")
		assertNoError(t, err, "inserting weekly entry")
	}
	{
		err := insertRepeating(tx, e2, "monthly")
		assertNoError(t, err, "inserting repeating entry")
	}
	testcommit(t, tx)

	// Then
	endDate := time.Now().AddDate(2, 0, 0)
	{
		tx := testtx(t, db)
		result, err := summarizeBucket(tx, e1.source, endDate)
		assertNoError(t, err, "")
		testcommit(t, tx)
		assertEqual(t, -e1.amount*105-e2.amount*25, result, "inserting weekly")
	}
	{
		tx := testtx(t, db)
		result, err := summarizeBucket(tx, e2.destination, endDate)
		assertNoError(t, err, "")
		testcommit(t, tx)
		assertEqual(t, e2.amount*25, result, "inserting monthly")
	}

}

func TestSummarizeAllThroughDate(t *testing.T) {
	// Given
	db := testdb(t)
	earlyDate := time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local)
	laterDate := time.Date(2021, 4, 1, 0, 0, 0, 0, time.Local)
	entries := []entry{
		{
			source:      "savings",
			destination: "checking",
			happenedAt:  earlyDate,
			amount:      500,
		},
		{
			source:      "checking",
			destination: "credit card",
			happenedAt:  earlyDate,
			amount:      1250,
		},
		{
			source:      "checking",
			destination: "credit card",
			happenedAt:  laterDate,
			amount:      20,
		},
	}
	tx := testtx(t, db)
	for _, e := range entries {
			err := insert(tx, e)
			assertNoError(t, err, "inserting transaction")
	}
	testcommit(t, tx)

	// When
	tx = testtx(t, db)
	result, err := summarizeAllThroughDate(tx, earlyDate)
	assertNoError(t, err, "summarizing all buckets through date")
	testcommit(t, tx)
	want := map[string]int{
		"checking":    -750,
		"credit card": 1250,
		"savings":     -500,
	}

	// Then
	assertEqual(t, want, result, "")
}

func TestGetAssets(t *testing.T) {
	// Given
	db := testdb(t)
	entryDate := time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local)
	entries := []entry{
		{
			source:      "savings",
			destination: "checking",
			happenedAt:  entryDate,
			amount:      500,
		},
		{
			source:      "checking",
			destination: "credit card",
			happenedAt:  entryDate,
			amount:      1250,
		},
		{
			source:      "paycheck",
			destination: "checking",
			happenedAt:  entryDate,
			amount:      200,
		},
	}
	buckets := []bucket{
		{
			name:      "savings",
			asset:     true,
			liquidity: "full",
		},
		{
			name:      "checking",
			asset:     true,
			liquidity: "full",
		},
		{
			name:      "credit card",
			asset:     false,
			liquidity: "",
		},
		{
			name:      "paycheck",
			asset:     false,
			liquidity: "",
		},
	}
	tx := testtx(t, db)
	for _, e := range entries {
		err := insert(tx, e)
		assertNoError(t, err, "inserting entries")
	}
	for _, b := range buckets {
		err := addBucket(tx, buckets)
		assertNoError(t, err, "classifying buckets")
	}
	testcommit(t, tx)

	// When
	result, err := sumAssets(tx, entryDate.AddDate(0, 0, 1))
	assertNoError(t, err, "summing assets")
	want := -entries[1].amount + entries[2].amount

	// Then
	assertEqual(t, want, result, "checking equality of sumAssets")

}

func TestWhenZero(t *testing.T) {
	// Given
	db := testdb(t)
	e1 := entry{
		source:      "savings",
		destination: "checking",
		happenedAt:  time.Now(),
		amount:      500,
	}
	e2 := entry{
		source:      "checking",
		destination: "rent",
		happenedAt:  time.Now(),
		amount:      150,
	}
	tx := testtx(t, db)
	err := insert(tx, e1)
	assertNoError(t, err, "inserting one entry")
	err = insertRepeating(tx, e2, "monthly")
	assertNoError(t, err, "inserting repeating entry")
	testcommit(t, tx)

	// When
	tx = testtx(t, db)
	result, err := findWhenZero(tx, e2.source)
	assertNoError(t, err, "finding when bucket hits zero")
	testcommit(t, tx)
	want := convertToDate(time.Now()).AddDate(0, 3, 0)

	// Then
	assertEqual(t, want, result, "")
}

// helper functions
func testdb(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("opening database: %v", err)
	}
	schema, err := ioutil.ReadFile("./schema.sql")
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
