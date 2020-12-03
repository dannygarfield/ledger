CREATE TABLE transactions (source TEXT, destination TEXT, happened_at TEXT, amount TEXT);

CREATE TABLE buckets (name TEXT PRIMARY KEY, asset INT, liquidity TEXT);
