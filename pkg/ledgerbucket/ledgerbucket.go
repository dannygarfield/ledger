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
func GetBuckets(tx *sql.Tx) ([]string, error) {
	q := `SELECT name FROM buckets`
	rows, err := tx.Query(q)
	if err != nil {
		return nil, fmt.Errorf("querying db (%w)", err)
	}
	var bucketNames []string
	for rows.Next() {
		var b string
		if err := rows.Scan(&b); err != nil {
			return nil, fmt.Errorf("scanning rows (%w)", err)
		}
		bucketNames = append(bucketNames, b)
	}
	return bucketNames, nil
}
