package testutils

import (
	"database/sql"
	"io/ioutil"
	"testing"
)

func Db(t *testing.T) *sql.DB {
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

func Tx(t *testing.T, db *sql.DB, work func(tx *sql.Tx) error) {
	t.Helper()
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("starting tx: %v", err)
	}
	err = work(tx)
	if err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			t.Fatalf("rolling back tx: %v (caused by %v)", err2, err)
		}
		t.Fatalf("error caused rollback: %v", err)
	}
	if err := tx.Commit(); err != nil {
		t.Fatalf("committing tx: %v", err)
	}
}
