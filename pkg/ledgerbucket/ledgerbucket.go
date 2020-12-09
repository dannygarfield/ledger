package ledgerbucket

import (
	"database/sql"
	"fmt"
)

// a Bucket describes ownership and accessibility of money
type Bucket struct {
	Name      string
	Asset     bool
	Liquidity string
}

// add a bucket to the db
func AddBucket(tx *sql.Tx, bucket Bucket) error {
	q := `INSERT INTO buckets
		(name, asset, liquidity)
		VALUES ($1, $2, $3)`
	var x int
	if bucket.Asset == true {
		x = 1
	} else {
		x = 0
	}
	_, err := tx.Exec(q, bucket.Name, x, bucket.Liquidity)
	if err != nil {
		return fmt.Errorf("addBuckets() - executing query: %w", err)
	}
	return nil
}
