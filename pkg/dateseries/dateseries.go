package dateseries

import (
	"database/sql"
	"fmt"
	"ledger/pkg/ledger"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func CreateSeries(tx *sql.Tx) error {
	q := `INSERT INTO dateseries (day) VALUES ($1)`

	today := ledger.ConvertToDate(time.Now())
	fmt.Printf("TODAY: %v\n", today)
	for t := today; t.Before(today.AddDate(2, 0, 0)); t = t.AddDate(0, 0, 1) {
		_, err := tx.Exec(q, t)
		if err != nil {
			return fmt.Errorf("CreateSeries() executing insert: %w", err)
		}
	}
	return nil
}

func GetMaxDate(tx *sql.Tx) (time.Time, error) {
	zeroTime := time.Date(0001, 1, 1, 00, 00, 00, 00, time.Local)
	q := `SELECT MAX(day) FROM dateseries`
	row := tx.QueryRow(q)
	var d string
	if err := row.Scan(&d); err != nil {
		return zeroTime, fmt.Errorf("GetMaxDate() scanning row")
	}
	fmt.Printf("var d: %s\n", d)
	day, err := time.Parse("2006-01-02 15:04:05-07:00", d)
	if err != nil {
		return zeroTime, fmt.Errorf("GetMaxDate() parsing date")
	}
	return day, nil
}
