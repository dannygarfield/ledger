package utils

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"
)

var BigBang = time.Date(1996, 4, 11, 0, 0, 0, 0, time.Local)

// var bb = time.Date(1996, 4, 11, 0, 0, 0, 0, time.Local)

func Tx(db *sql.DB, r *http.Request, work func(tx *sql.Tx) error) {
	ctx := r.Context()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		fmt.Printf("Could not call BeginTx() (%v)", err)
		return
	}

	workErr := work(tx)
	if workErr != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			fmt.Printf("Error on rollback (%v) -- rollback caused by work() (%v)", rollbackErr, workErr)
			return
		}
		fmt.Printf("Could not execute work() (%v)", err)
		return
	}
	if err := tx.Commit(); err != nil {
		fmt.Printf("Could not commit sql tx (%v)", err)
		return
	}
}

// return a time with year, month, and day values; all other values equal 0
func ConvertToDate(t time.Time) time.Time {
	year := t.Year()
	month := t.Month()
	day := t.Day()
	return time.Date(year, month, day, 0, 0, 0, 0, time.Local)
}

// parse a date
func ParseDate(s string) (time.Time, error) {
	d, err := time.Parse("2006-01-02", s)
	if err != nil {
		return time.Time{}, fmt.Errorf("parseDate() - parsing time: %w", err)
	}
	return d, nil
}
