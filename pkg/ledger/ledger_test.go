package ledger

import (
	"database/sql"
	"io/ioutil"
	"ledger/pkg/ledgerbucket"
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

func TestSummarizeLedger(t *testing.T) {
	// Given
	db := testdb(t)
	earlyDate := time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local)
	laterDate := time.Date(2021, 4, 1, 0, 0, 0, 0, time.Local)
	entries := []Entry{
		{
			Source:      "savings",
			Destination: "checking",
			EntryDate:   earlyDate,
			Amount:      1000,
		},
		{
			Source:      "checking",
			Destination: "credit card",
			EntryDate:   earlyDate,
			Amount:      200,
		},
		{
			Source:      "checking",
			Destination: "credit card",
			EntryDate:   laterDate,
			Amount:      100,
		},
	}
	buckets := []ledgerbucket.Bucket{
		{
			Name:      "savings",
			Asset:     1,
			Liquidity: "full",
		},
		{
			Name:      "checking",
			Asset:     1,
			Liquidity: "full",
		},
		{
			Name:      "credit card",
			Asset:     0,
			Liquidity: "",
		},
		{
			Name:      "paycheck",
			Asset:     0,
			Liquidity: "",
		},
	}
	tx := testtx(t, db)
	for _, e := range entries {
		err := InsertEntry(tx, e)
		assertNoError(t, err, "inserting entries")
	}
	for _, b := range buckets {
		err := ledgerbucket.InsertBucket(tx, b)
		assertNoError(t, err, "classifying buckets")
	}
	testcommit(t, tx)

	// When
	tx = testtx(t, db)
	result, err := SummarizeLedger(tx, earlyDate)
	assertNoError(t, err, "summarizing all buckets through date")
	testcommit(t, tx)
	want := []balanceDetail{
		{
			bucket:    "checking",
			amount:    1000 - 200,
			asset:     1,
			liquidity: "full",
		},
		{
			bucket:    "savings",
			amount:    -1000,
			asset:     1,
			liquidity: "full",
		},
		{
			bucket:    "credit card",
			amount:    200,
			asset:     0,
			liquidity: "",
		},
	}

	// Then
	assertEqual(t, want, result, "")
}

func TestSummarizeBalanceOverTime(t *testing.T) {
	// Given
	db := testdb(t)
	day1 := time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local)
	entries := []Entry{
		{
			Source:      "checking",
			Destination: "IRA",
			Amount:      50,
			EntryDate:   day1,
		},
		{
			Source:      "checking",
			Destination: "IRA",
			Amount:      50,
			EntryDate:   day1.AddDate(0, 0, 1),
		},
	}
	buckets := []ledgerbucket.Bucket{
		{
			Name:      "checking",
			Asset:     1,
			Liquidity: "full",
		},
		{
			Name:      "IRA",
			Asset:     1,
			Liquidity: "low",
		},
	}
	{
		tx := testtx(t, db)
		for _, e := range entries {
			err := InsertEntry(tx, e)
			assertNoError(t, err, "inserting entries")
		}
		for _, b := range buckets {
			err := ledgerbucket.InsertBucket(tx, b)
			assertNoError(t, err, "classifying buckets")
		}
		testcommit(t, tx)
	}

	// When
	tx := testtx(t, db)
	end := day1.AddDate(0, 0, 2)
	bot, err := SummarizeBalanceOverTime(tx, day1, end)
	assertNoError(t, err, "summarizing balances over time")

	// Then
	want := []dailyBalanceDetail{
		{
			day: day1,
			balances: []balanceDetail{
				{
					bucket:    "checking",
					amount:    -50,
					asset:     1,
					liquidity: "full",
				},
				{
					bucket:    "IRA",
					amount:    50,
					asset:     1,
					liquidity: "low",
				},
			},
		},
		{
			day: day1.AddDate(0, 0, 1),
			balances: []balanceDetail{
				{
					bucket:    "checking",
					amount:    -100,
					asset:     1,
					liquidity: "full",
				},
				{
					bucket:    "IRA",
					amount:    100,
					asset:     1,
					liquidity: "low",
				},
			},
		},
		{
			day: day1.AddDate(0, 0, 2),
			balances: []balanceDetail{
				{
					bucket:    "checking",
					amount:    -100,
					asset:     1,
					liquidity: "full",
				},
				{
					bucket:    "IRA",
					amount:    100,
					asset:     1,
					liquidity: "low",
				},
			},
		},
	}
	assertEqual(t, want, bot, "")

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
