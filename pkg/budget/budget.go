package budget

import (
	"database/sql"
	"fmt"
	"ledger/pkg/utils"
	"net/http"
	"strconv"
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
	_, err := tx.Exec(q, happened_at, e.Amount, e.Category, e.Description)
	if err != nil {
		return fmt.Errorf("calling budget.InsertEntry() (%w)", err)
	}
	return nil
}

// get all entries from budget in given time period
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
		if e.EntryDate, err = time.Parse("2006-01-02", datestring); err != nil {
			return nil, err
		}
		budget = append(budget, e)
	}
	return budget, nil
}

// get net spend of category from start through end
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

// get net spend of categories from start through end
func SummarizeCategories(tx *sql.Tx, categories []string, from, through time.Time) (map[string]int, error) {
	output := map[string]int{}
	for _, c := range categories {
		val, err := SummarizeCategory(tx, c, from, through)
		if err != nil {
			return nil, fmt.Errorf("calling SummarizeCategory() (%v)", err)
		}
		output[c] = val
	}
	return output, nil
}

func SummarizeSpendsOverTime(tx *sql.Tx, categories []string, start, end time.Time, interval int) ([]map[string]int, error) {
	output := []map[string]int{}
	for d := start; d.Before(end.AddDate(0, 0, 1)); d = d.AddDate(0, 0, interval) {
		// summarize from the start to end of an interval period
		c, err := SummarizeCategories(tx, categories, d, d.AddDate(0, 0, interval-1))
		if err != nil {
			return nil, fmt.Errorf("calling SummarizeCategories() (%w)", err)
		}
		output = append(output, c)
	}
	return output, nil
}

func PrepareEntryForInsert(r *http.Request) (Entry, error) {
	r.ParseForm()
	entrydate, err := time.Parse("2006-01-02", r.PostForm["happened_at"][0])
	if err != nil {
		return Entry{}, fmt.Errorf("Could not parse entrydate (%v)", err)
	}
	amount, err := strconv.Atoi(r.PostForm["amount"][0])
	if err != nil {
		return Entry{}, fmt.Errorf("Could not convert amount field to int (%v)", err)
	}
	entry := Entry{
		EntryDate:   entrydate,
		Amount:      amount,
		Category:    r.PostForm["category"][0],
		Description: r.PostForm["description"][0],
	}
	return entry, nil
}

func GetCategories(tx *sql.Tx) ([]string, error) {
	q := `SELECT DISTINCT category FROM budget_entries ORDER BY category;`
	rows, err := tx.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	categories := []string{}
	for rows.Next() {
		var c string
		if err := rows.Scan(&c); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}
