package ledger

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
	e := Entry{
		Source:      "checking",
		Destination: "credit card",
		HappenedAt:  time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local),
		Amount:      120,
	}

	// When
	tx := testtx(t, db)
	err := Insert(tx, e)
	assertNoError(t, err, "")
	testcommit(t, tx)

	// Then
	{
		tx := testtx(t, db)
		result, err := SummarizeBucket(tx, e.Source, e.HappenedAt)
		assertNoError(t, err, "summary(source)")
		testcommit(t, tx)
		assertEqual(t, -e.Amount, result, "source")
	}
	{
		tx := testtx(t, db)
		result, err := SummarizeBucket(tx, e.Destination, e.HappenedAt)
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
		HappenedAt:  time.Now(), // repeating write until 2 years from now. setting happenedAt to time.Now() requires less math
	}
	e2 := Entry{
		Source:      "checking",
		Destination: "rent",
		Amount:      50,
		HappenedAt:  time.Now(),
	}

	// When
	tx := testtx(t, db)
	{
		err := InsertRepeating(tx, e1, "weekly")
		assertNoError(t, err, "inserting weekly entry")
	}
	{
		err := InsertRepeating(tx, e2, "monthly")
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

func TestSummarizeAllThroughDate(t *testing.T) {
	// Given
	db := testdb(t)
	earlyDate := time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local)
	laterDate := time.Date(2021, 4, 1, 0, 0, 0, 0, time.Local)
	entries := []Entry{
		{
			Source:      "savings",
			Destination: "checking",
			HappenedAt:  earlyDate,
			Amount:      500,
		},
		{
			Source:      "checking",
			Destination: "credit card",
			HappenedAt:  earlyDate,
			Amount:      1250,
		},
		{
			Source:      "checking",
			Destination: "credit card",
			HappenedAt:  laterDate,
			Amount:      20,
		},
	}
	tx := testtx(t, db)
	for _, e := range entries {
		err := Insert(tx, e)
		assertNoError(t, err, "inserting transaction")
	}
	testcommit(t, tx)

	// When
	tx = testtx(t, db)
	result, err := SummarizeAllThroughDate(tx, earlyDate)
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
	entries := []Entry{
		{
			Source:      "savings",
			Destination: "checking",
			HappenedAt:  entryDate,
			Amount:      500,
		},
		{
			Source:      "checking",
			Destination: "credit card",
			HappenedAt:  entryDate,
			Amount:      1250,
		},
		{
			Source:      "paycheck",
			Destination: "checking",
			HappenedAt:  entryDate,
			Amount:      200,
		},
	}
	buckets := []Bucket{
		{
			Name:      "savings",
			Asset:     true,
			Liquidity: "full",
		},
		{
			Name:      "checking",
			Asset:     true,
			Liquidity: "full",
		},
		{
			Name:      "credit card",
			Asset:     false,
			Liquidity: "",
		},
		{
			Name:      "paycheck",
			Asset:     false,
			Liquidity: "",
		},
	}
	tx := testtx(t, db)
	for _, e := range entries {
		err := Insert(tx, e)
		assertNoError(t, err, "inserting entries")
	}
	for _, b := range buckets {
		err := AddBucket(tx, b)
		assertNoError(t, err, "classifying buckets")
	}
	testcommit(t, tx)

	// When
	tx = testtx(t, db)
	result, err := SumAssets(tx, entryDate.AddDate(0, 0, 1))
	assertNoError(t, err, "summing assets")
	testcommit(t, tx)
	want := -entries[1].Amount + entries[2].Amount

	// Then
	assertEqual(t, want, result, "checking equality of sumAssets")

}

func TestWhenZero(t *testing.T) {
	// Given
	db := testdb(t)
	e1 := Entry{
		Source:      "savings",
		Destination: "checking",
		HappenedAt:  time.Now(),
		Amount:      500,
	}
	e2 := Entry{
		Source:      "checking",
		Destination: "rent",
		HappenedAt:  time.Now(),
		Amount:      150,
	}
	tx := testtx(t, db)
	err := Insert(tx, e1)
	assertNoError(t, err, "inserting one entry")
	err = InsertRepeating(tx, e2, "monthly")
	assertNoError(t, err, "inserting repeating entry")
	testcommit(t, tx)

	// When
	tx = testtx(t, db)
	result, err := FindWhenZero(tx, e2.Source)
	assertNoError(t, err, "finding when bucket hits zero")
	testcommit(t, tx)
	want := ConvertToDate(time.Now()).AddDate(0, 3, 0)

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
