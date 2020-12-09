package ledger

import (
	"database/sql"
	"fmt"
	"time"
)

// get net amounts of all buckets through a given date
func SummarizeLedger(tx *sql.Tx, through time.Time) (map[string]int, error) {
	q := `SELECT account, sum(amount) FROM (
		SELECT amount, happened_at, destination AS account FROM entries
		UNION ALL
		SELECT -amount, happened_at, source AS account FROM entries
		)
		WHERE date(happened_at) <= date($1)
		GROUP BY account
		ORDER BY sum(amount) DESC;`
	rows, err := tx.Query(q, through)
	if err != nil {
		return nil, fmt.Errorf("summarizeLedger() - querying rows: %w", err)
	}
	result := make(map[string]int)
	for rows.Next() {
		var account string
		var total int
		if err := rows.Scan(&account, &total); err != nil {
			return nil, fmt.Errorf("summarizeLedger() - iterating through rows: %w", err)
		}
		result[account] = total
		// log.Printf("%s: %d \n", account, total)
	}
	return result, nil
}

// get net amount of a single bucket through a given date
func SummarizeBucket(tx *sql.Tx, bucket string, through time.Time) (int, error) {
	q := `SELECT sum(amount) FROM (
		SELECT amount, happened_at FROM entries WHERE destination = $1
		UNION ALL
		SELECT -amount, happened_at from entries where source = $1
		)
		WHERE date(happened_at) <= date($2)
		ORDER BY sum(amount) DESC;`
	row := tx.QueryRow(q, bucket, through)
	var sum int
	if err := row.Scan(&sum); err != nil {
		return -1, fmt.Errorf("summarizeBucket() - querying rows: %w", err)
	}
	return sum, nil
}

// get total assets owned on a given date
func SumAssets(tx *sql.Tx, through time.Time) (int, error) {
	q := `SELECT sum(amount) FROM (
    	SELECT destination AS account, amount, happened_at
		FROM entries
        UNION ALL
        SELECT source AS account, -amount, happened_at
		FROM entries
        ) t
        LEFT JOIN buckets b
        ON t.account = b.name
        WHERE date(t.happened_at) < date($1) AND b.asset = 1;`
	row := tx.QueryRow(q, through)
	var sum int
	if err := row.Scan(&sum); err != nil {
		return -1, fmt.Errorf("sumAssets() - scanning rows: %w", err)
	}
	return sum, nil
}

// find the first date after today that a bucket becomes <= 0
func FindWhenZero(tx *sql.Tx, bucket string) (time.Time, error) {
	today := ConvertToDate(time.Now())
	for t := today; t.Before(today.AddDate(2, 0, 0)); t = t.AddDate(0, 0, 1) {
		balance, err := SummarizeBucket(tx, bucket, t)
		if err != nil {
			return time.Now(), fmt.Errorf("findWhenZero() - summarizing bucket: %w", err)
		}
		if balance <= 0 {
			return t, nil
		}
	}
	return time.Now(), nil
}
