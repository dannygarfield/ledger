package ledger

import (
	"database/sql"
	"fmt"
	"ledger/pkg/utils"
	"sort"
	"time"
)

type PlotData struct {
	BucketHeaders []string
	DateHeaders   []string
	Data          [][]int
}

// get all entries in the ledger from start through finish
func GetLedger(tx *sql.Tx, start, end time.Time) ([]Entry, error) {
	q := `SELECT * FROM entries
		WHERE date(happened_at) >= date($1) AND date(happened_at) < date($2)
		ORDER BY happened_at;`

	rows, err := tx.Query(q, start, end)
	if err != nil {
		return nil, fmt.Errorf("Could not Query sql (%v)", err)
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

// get net amount of a single bucket over a given time
func SummarizeBucket(tx *sql.Tx, bucket string, start, end time.Time) (int, error) {
	q := `SELECT COALESCE(sum(amount), 0) FROM (
		SELECT amount, happened_at FROM entries WHERE destination = $1
		UNION ALL
		SELECT -amount, happened_at from entries where source = $1
		)
		WHERE date(happened_at) BETWEEN date($2) AND date($3)
		ORDER BY sum(amount) DESC;`
	row := tx.QueryRow(q, bucket, start, end)
	var sum int
	if err := row.Scan(&sum); err != nil {
		return -1, fmt.Errorf("summarizeBucket() - querying rows: %w", err)
	}
	return sum, nil
}

// get net amounts of provided buckets over a given time
func SummarizeBalance(tx *sql.Tx, buckets []string, from, through time.Time) (map[string]int, error) {
	out := map[string]int{}
	for _, b := range buckets {
		val, err := SummarizeBucket(tx, b, from, through)
		if err != nil {
			return nil, err
		}
		out[b] = val
	}
	return out, nil
}

// get daily balances (starting from bigBang) of provided buckets over a given time
func SummarizeBalanceOverTime(tx *sql.Tx, buckets []string, start, end time.Time) ([]map[string]int, error) {
	bigBang := utils.BigBang()
	output := []map[string]int{}
	for d := start; d.Before(end); d = d.AddDate(0, 0, 1) {
		// summarize from beginning of time through date iterator
		balance, err := SummarizeBalance(tx, buckets, bigBang, d)
		if err != nil {
			return nil, fmt.Errorf("ledger.SummarizeBalanceOverTime() summarizing ledger (%w)", err)
		}
		output = append(output, balance)
	}
	return output, nil
}

// get totals over time, grouped into provided intervals of time
func SummarizeLedgerOverTime(tx *sql.Tx, buckets []string, start, end time.Time, interval int) ([]map[string]int, error) {
	output := []map[string]int{}
	for d := start; d.Before(end); d = d.AddDate(0, 0, interval) {
		// summarize from the start to end of an interval period
		l, err := SummarizeBalance(tx, buckets, d, d.AddDate(0, 0, interval-1))
		if err != nil {
			return nil, fmt.Errorf("ledger.SummarizeEntriesOverTime() summarizing ledger (%w)", err)
		}
		output = append(output, l)
	}
	return output, nil
}

func MakePlot(summary []map[string]int, start time.Time, interval int) *PlotData {
	output := &PlotData{}

	if len(summary) > 0 {
		for b := range summary[0] {
			output.BucketHeaders = append(output.BucketHeaders, b)
		}
		sort.Strings(output.BucketHeaders)

		for i, day := range summary {
			output.DateHeaders = append(output.DateHeaders, start.AddDate(0, 0, i*interval).Format("2006-01-02"))
			row := []int{}
			for _, b := range output.BucketHeaders {
				row = append(row, day[b])
			}
			output.Data = append(output.Data, row)
		}
	}

	return output
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
