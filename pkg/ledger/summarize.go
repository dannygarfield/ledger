package ledger

import (
	"database/sql"
	"fmt"
	"time"
)

func GetLedger(tx *sql.Tx, start, end time.Time) ([]Entry, error) {
	q := `SELECT * FROM entries
		WHERE date(happened_at) >= date($1) AND date(happened_at) < date($2)
		ORDER BY happened_at;`

	rows, err := tx.Query(q, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ledger []Entry
	for rows.Next() {
		e := Entry{}
		var datestring string
		if err := rows.Scan(&e.Source, &e.Destination, &datestring, &e.Amount); err != nil {
			return nil, err
		}
		if e.EntryDate, err = time.Parse("2006-01-02 15:04:05-07:00", datestring); err != nil {
			return nil, err
		}
		ledger = append(ledger, e)
	}
	return ledger, nil
}

// get net amount of a single bucket through a given date
func SummarizeBucket(tx *sql.Tx, bucket string, through time.Time) (int, error) {
	q := `SELECT COALESCE(sum(amount), 0) FROM (
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

// get net amounts of all buckets through a given date
func SummarizeLedger(tx *sql.Tx, buckets []string, through time.Time) (map[string]int, error) {
	out := map[string]int{}
	for _, b := range buckets {
		val, err := SummarizeBucket(tx, b, through)
		if err != nil {
			return nil, err
		}
		out[b] = val
	}
	return out, nil
}

func SummarizeLedgerOverTime(tx *sql.Tx, buckets []string, start, end time.Time) ([]map[string]int, error) {
	output := []map[string]int{}
	for d := start; d.Before(end); d = d.AddDate(0, 0, 1) {
		l, err := SummarizeLedger(tx, buckets, d)
		if err != nil {
			return nil, fmt.Errorf("ledger.SummarizeLedgerOverTime() summarizing ledger (%w)", err)
		}
		output = append(output, l)
	}
	return output, nil
}

func GetBuckets(tx *sql.Tx) ([]string, error) {
	q := `SELECT DISTINCT buckets FROM (
		    SELECT source AS buckets FROM entries
		    UNION
		    SELECT destination AS buckets FROM entries
		) ORDER BY buckets
	;`

	rows, err := tx.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	buckets := []string{}
	for rows.Next() {
		var b string
		if err := rows.Scan(&b); err != nil {
			return nil, err
		}
		buckets = append(buckets, b)
	}
	return buckets, nil

}
