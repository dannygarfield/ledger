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
		VALUES (date($1), $2, $3, $4);`
	happened_at := utils.ConvertToDate(e.EntryDate)
	fmt.Println("happened_at:", happened_at)
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
		fmt.Println("parsing e.EntryDate")
		if e.EntryDate, err = time.Parse("2006-01-02", datestring); err != nil {
			return nil, err
		}
		fmt.Println("parsed e.EntryDate")
		budget = append(budget, e)
	}
	return budget, nil
}

func SummarizeCategory(tx *sql.Tx, category string, start, end time.Time) (int, error) {
	q := `SELECT COALESCE(sum(amount), 0)
		FROM budget_entries
		WHERE category = $1
		AND
		date(happened_at) BETWEEN date($2) AND date($3)`
	row := tx.QueryRow(q, category, start, end)
	var sum int
	if err := row.Scan(&sum); err != nil {
		return -1, fmt.Errorf("calling row.Scan() (%w)", err)
	}
	return sum, nil
}
