package testutils

import (
	"database/sql"
	"io/ioutil"
	"reflect"
	"testing"
	"time"
)

var BigBang = time.Date(1996, 04, 11, 0, 0, 0, 0, time.UTC)

var Dec31 = time.Date(2020, 12, 31, 0, 0, 0, 0, time.UTC)

var JanOne = time.Date(2021, 01, 01, 0, 0, 0, 0, time.UTC)

var JanTwo = time.Date(2021, 01, 02, 0, 0, 0, 0, time.UTC)

func Db(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("opening database: %v", err)
	}
	schema, err := ioutil.ReadFile("../../pkg/testutils/testdata/schema.sql")
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

func AssertEqual(t *testing.T, want, got interface{}) {
	t.Helper()
	if equal := reflect.DeepEqual(want, got); !equal {
		t.Fatalf("want: %v, got: %v", want, got)
	}
}
