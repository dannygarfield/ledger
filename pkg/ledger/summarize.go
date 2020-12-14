package ledger

import (
	"database/sql"
	"fmt"
	"time"
)

type balanceDetail struct {
	bucket    string
	amount    int
	asset     int
	liquidity string
}

type dailyBalanceDetail struct {
	day      time.Time
	balances []balanceDetail
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

func MakePlot(tx *sql.Tx, buckets []string, start, end time.Time) ([]map[string]int, error) {
	out := []map[string]int{}
	for d := start; d.Before(end); d = d.AddDate(0, 0, 1) {
		m := map[string]int{}
		out = append(out, m)
		for _, bucket := range buckets {
			val, err := SummarizeBucket(tx, bucket, d)
			if err != nil {
				return nil, err
			}
			m[bucket] = val
		}
	}
	return out, nil
}
