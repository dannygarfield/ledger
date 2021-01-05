package utils

import (
	"database/sql"
	"fmt"
	"net/http"
)

func Tx(db *sql.DB, r *http.Request, work func(tx *sql.Tx) error) {
	ctx := r.Context()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		fmt.Printf("Could not call BeginTx() (%v)", err)
	}

	workErr := work(tx)
	if workErr != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			fmt.Printf("Error on rollback (%v) -- rollback caused by work() (%v)", rollbackErr, workErr)
		}
		fmt.Printf("Could not execute work() (%v)", err)
	}
	if err := tx.Commit(); err != nil {
		fmt.Printf("Could not commit sql tx (%v)", err)
	}
}
