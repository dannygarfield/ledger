CREATE TABLE IF NOT EXISTS entries
(
    source TEXT,
    destination TEXT,
    happened_at TEXT,
    amount TEXT
);

CREATE TABLE IF NOT EXISTS budget_entries
(
    happened_at TEXT,
    amount INT,
    category TEXT,
    description TEXT
);
