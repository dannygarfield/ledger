package sqlstatements

import (
	"database/sql"
	"fmt"
)

// begin a sql transaction
func BeginTx(db *sql.DB) (*sql.Tx, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("BeginTx() - beginning a sql transaction: %w", err)
	}
	return tx, nil
}
