package main

import (
	"database/sql"
	"io/ioutil"
	"log"
	"reflect"
	"testing"
	"time"
)

func TestInsertOne(t *testing.T) {
	// Given
	db := testdb(t)
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("beginning the sql transaction")
	}
	// When
	e := entry {
		source: "checking",
		destination: "credit card",
		happenedAt: time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local),
		amount: 12500,
	}
	if err := insert(tx, e); err != nil {
		t.Fatalf("inserting record: %v", err)
	}
	if err := tx.Commit(); err != nil {
		log.Fatalf("committing the transaction: %v", err)
	}

	// Then
	{
		result, err := summary(db, e.source, e.happenedAt)
		assertNoError(t, err, "summary(source)")
		assertEqual(t, -e.amount, result, "source")
	}
	{
		result, err := summary(db, e.destination, e.happenedAt)
		assertNoError(t, err, "summary(destination)")
		assertEqual(t, e.amount, result, "destination")
	}
}

func TestSummarizeAllThroughDate(t *testing.T) {
	// Given
	db := testdb(t)
	earlyDate := time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local)
	laterDate := time.Date(2021, 4, 1, 0, 0, 0, 0, time.Local)
	entries := []entry{
		{
			source:      "checking",
			destination: "credit card",
			happenedAt: earlyDate,
			amount:      125000,
		},
		{
			source:      "checking",
			destination: "credit card",
			happenedAt: laterDate,
			amount:      2000,
		},
		{
			source:      "savings",
			destination: "checking",
			happenedAt: earlyDate,
			amount:      50000,
		},
	}
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("beginning the sql transaction")
	}
	for _, e := range entries {
		err := insert(tx, e)
		assertNoError(t, err, "inserting transaction into tx")
	}
	if err := tx.Commit(); err != nil {
		log.Fatalf("committing the transaction: %v", err)
	}

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
