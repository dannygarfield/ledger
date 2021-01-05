package ledger

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// transaction represents a double-Entry accounting item in the ledger.
type Entry struct {
	Source      string
	Destination string
	EntryDate   time.Time
	Amount      int
}

// Insert an entry
func InsertEntry(tx *sql.Tx, e Entry) error {
	q := `INSERT INTO entries
		(source, destination, happened_at, amount)
		VALUES ($1, $2, $3, $4);`
	_, err := tx.Exec(q, e.Source, e.Destination, e.EntryDate, e.Amount)
	if err != nil {
		return fmt.Errorf("insert() - executing the insert: %w", err)
	}
	return nil
}

// insert a transaction that repeats weekly or monthly
func InsertRepeatingEntry(tx *sql.Tx, e Entry, freq string) error {
	q := `INSERT INTO entries
		(source, destination, happened_at, amount)
		VALUES ($1, $2, $3, $4);`
	var freqMonth int
	var freqDay int
	if freq == "monthly" {
		freqMonth = 1
		freqDay = 0
	} else if freq == "weekly" {
		freqMonth = 0
		freqDay = 7
	}
	endDate := time.Now().AddDate(2, 0, 0)
	for e.EntryDate.Before(endDate) {
		if _, err := tx.Exec(q, e.Source, e.Destination, e.EntryDate, e.Amount); err != nil {
			return fmt.Errorf("insertRepeating() - inserting transactions: %w", err)
		}
		e.EntryDate = e.EntryDate.AddDate(0, freqMonth, freqDay)
	}
	return nil
}

func PrepareEntryForInsert(r *http.Request) (Entry, error) {
	r.ParseForm()
	entrydate, err := time.Parse("2006-01-02", r.PostForm["happened_at"][0])
	if err != nil {
		return Entry{}, fmt.Errorf("Could not parse entrydate (%v)", err)
	}
	amount, err := strconv.Atoi(r.PostForm["amount"][0])
	if err != nil {
		return Entry{}, fmt.Errorf("Could not convert amount field to int (%v)", err)
	}
	entry := Entry{
		Source:      r.PostForm["source"][0],
		Destination: r.PostForm["destination"][0],
		EntryDate:   entrydate,
		Amount:      amount,
	}
	return entry, nil
}

// parse a date
func ParseDate(s string) (time.Time, error) {
	d, err := time.Parse("2006-01-02", s)
	if err != nil {
		return time.Time{}, fmt.Errorf("parseDate() - parsing time: %w", err)
	}
	return d, nil
}

// return a time with year, month, and day values; all other values equal 0
func ConvertToDate(t time.Time) time.Time {
	year := t.Year()
	month := t.Month()
	day := t.Day()
	return time.Date(year, month, day, 0, 0, 0, 0, time.Local)
}
