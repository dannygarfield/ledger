package main

import (
	"database/sql"
	"io/ioutil"
	"reflect"
	"testing"
	"time"
)

func TestInsertOne(t *testing.T) {
	// Given
	db := testdb(t)

	// When
	e := []entry{
		{
			source:      "checking",
			destination: "credit card",
			happenedAt:  time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local),
			amount:      12500,
		},
	}
	err := insert(db, e)
	assertNoError(t, err, "inserting one entry")

	// Then
	{
		result, err := summarizeBucket(db, e[0].source, e[0].happenedAt)
		assertNoError(t, err, "summary(source)")
		assertEqual(t, -e[0].amount, result, "source")
	}
	{
		result, err := summarizeBucket(db, e[0].destination, e[0].happenedAt)
		assertNoError(t, err, "summary(destination)")
		assertEqual(t, e[0].amount, result, "destination")
	}
}

func TestInsertRepeatingEntry(t *testing.T) {
	// Given
	db := testdb(t)

	// When
	e1 := entry{
		source:      "checking",
		destination: "IRA",
		amount:      5000,
		happenedAt:  time.Now(), // repeating write until 2 years from now. setting happenedAt to time.Now() requires less math
	}
	e2 := entry{
		source:      "checking",
		destination: "rent",
		amount:      5000,
		happenedAt:  time.Now(),
	}
	{
		err := insertRepeating(db, e1, "weekly")
		assertNoError(t, err, "inserting weekly entry")
	}
	{
		err := insertRepeating(db, e2, "monthly")
		assertNoError(t, err, "inserting repeating entry")
	}

	// Then
	endDate := time.Now().AddDate(2, 0, 0)
	{
		result, err := summarizeBucket(db, e1.source, endDate)
		assertNoError(t, err, "")
		assertEqual(t, -e1.amount*105-e2.amount*25, result, "inserting weekly")
	}
	{
		result, err := summarizeBucket(db, e2.destination, endDate)
		assertNoError(t, err, "")
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
			amount:      50000,
		},
		{
			source:      "checking",
			destination: "credit card",
			happenedAt:  earlyDate,
			amount:      125000,
		},
		{
			source:      "checking",
			destination: "credit card",
			happenedAt:  laterDate,
			amount:      2000,
		},
	}
	err := insert(db, entries)
	assertNoError(t, err, "inserting transaction")

	// When
	result, err := summarizeAllThroughDate(db, earlyDate)
	assertNoError(t, err, "summarizing all buckets through date")
	want := map[string]int{
		"checking":    -75000,
		"credit card": 125000,
		"savings":     -50000,
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
			amount:      50000,
		},
		{
			source:      "checking",
			destination: "credit card",
			happenedAt:  entryDate,
			amount:      125000,
		},
		{
			source:      "paycheck",
			destination: "checking",
			happenedAt:  entryDate,
			amount:      20000,
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
	err := insert(db, entries)
	assertNoError(t, err, "inserting entries")
	err = classifyBuckets(db, buckets)
	assertNoError(t, err, "classifying buckets")

	// When
	result, err := sumAssets(db, entryDate)
	assertNoError(t, err, "summing assets")
	want := -entries[1].amount + entries[2].amount

	// Then
	assertEqual(t, want, result, "checking equality of sumAssets")

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
