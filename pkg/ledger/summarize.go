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

// get net amounts of all buckets through a given date
func SummarizeLedger(tx *sql.Tx, through time.Time) ([]balanceDetail, error) {
	fmt.Printf("BALANCES THROUGH: %v\n", through)
	q := `SELECT account, sum(amount), b.asset, b.liquidity
	    FROM (
	        SELECT amount, happened_at, destination AS account FROM entries
	        UNION ALL
	        SELECT -amount, happened_at, source AS account FROM entries
	    ) l
	    LEFT JOIN buckets b
	    ON l.account = b.name
	    WHERE date(l.happened_at) <= date($1)
	    GROUP BY account
	    ORDER BY b.liquidity DESC, sum(amount) DESC;`
	rows, err := tx.Query(q, through)
	if err != nil {
		return nil, fmt.Errorf("summarizeLedger() - querying rows: %w", err)
	}

	var result []balanceDetail

	for rows.Next() {
		b := balanceDetail{}
		if err := rows.Scan(&b.bucket, &b.amount, &b.asset, &b.liquidity); err != nil {
			return nil, fmt.Errorf("summarize.SummarizeLedger() scanning rows (%w)", err)
		}
		result = append(result, b)
	}
	return result, nil
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
