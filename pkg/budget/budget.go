package budget

import (
	"database/sql"
	"fmt"
	"ledger/pkg/utils"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Entry struct {
	EntryDate   time.Time
	Amount      int
	Category    string
	Description string
}

func InsertEntry(tx *sql.Tx, e Entry) error {
	q := `INSERT INTO budget_entries
		(happened_at, amount, category, description)
		VALUES ($1, $2, $3, $4);`
	happened_at := utils.ConvertToDate(e.EntryDate)
	_, err := tx.Exec(q, happened_at, e.Amount, e.Category, e.Description)
	if err != nil {
		return fmt.Errorf("calling budget.InsertEntry() (%w)", err)
	}
	return nil
}

func GetBudgetEntries(tx *sql.Tx, start, end time.Time) ([]Entry, error) {
	q := `SELECT * FROM budget_entries
		WHERE date(happened_at) BETWEEN date($1) AND date($2)
		ORDER BY happened_at;`

	rows, err := tx.Query(q, start, end)
	if err != nil {
		return nil, fmt.Errorf("Could not query sql (%v)", err)
	}
	defer rows.Close()
	var budget []Entry
	for rows.Next() {
		e := Entry{}
		var datestring string
		if err := rows.Scan(&datestring, &e.Amount, &e.Category, &e.Description); err != nil {
			return nil, err
		}
		if e.EntryDate, err = time.Parse("2006-01-02 15:04:05-07:00", datestring); err != nil {
			return nil, err
		}
		budget = append(budget, e)
	}
	return budget, nil
}
