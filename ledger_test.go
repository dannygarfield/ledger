package main

import (
	"database/sql"
	"io/ioutil"
	"testing"
	"time"
)

func TestInsertOne(t *testing.T) {
	// Given
	db := testdb(t)
	// When
	source := "checking"
	destination := "credit card"
	happenedAt := time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local)
	amount := 12500
	if err := insert(db, source, destination, happenedAt, amount); err != nil {
		t.Fatalf("inserting record: %v", err)
	}
	// Then
	{
		result, err := summary(db, source, happenedAt)
		assertNoError(t, err, "summary(source)")
		assertEqual(t, -amount, result, "source")
	}
	{
		result, err := summary(db, destination, happenedAt)
		assertNoError(t, err, "summary(destination)")
		assertEqual(t, amount, result, "destination")
	}
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
	if want != got {
		t.Fatalf("%s: want: %v, got: %v", msg, want, got)
	}
}
