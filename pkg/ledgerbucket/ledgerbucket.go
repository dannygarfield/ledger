package ledgerbucket

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// a Bucket describes ownership and accessibility of money
type Bucket struct {
	Name      string
	Asset     int
	Liquidity string
}

// add a bucket to the db
func InsertBucket(tx *sql.Tx, bucket Bucket) error {
	q := `INSERT INTO buckets
		(name, asset, liquidity)
		VALUES ($1, $2, $3)`
	_, err := tx.Exec(q, bucket.Name, bucket.Asset, bucket.Liquidity)
	if err != nil {
		return fmt.Errorf("addBuckets() - executing query: %w", err)
	}
	return nil
}

// show all buckets in the db
func ShowBuckets(tx *sql.Tx) ([]Bucket, error) {
	q := `SELECT name, asset, liquidity FROM buckets`
	rows, err := tx.Query(q)
	if err != nil {
		return nil, fmt.Errorf("querying db (%w)", err)
	}
	var buckets []Bucket
	for rows.Next() {
		b := Bucket{}
		if err := rows.Scan(&b.Name, &b.Asset, &b.Liquidity); err != nil {
			return nil, fmt.Errorf("scanning rows (%w)", err)
		}
		buckets = append(buckets, b)
	}
	return buckets, nil
}
