CREATE TABLE IF NOT EXISTS entries
(source TEXT
, destination TEXT
, happened_at TEXT
, amount TEXT
);

CREATE TABLE IF NOT EXISTS buckets
(id INTEGER PRIMARY KEY
, name TEXT
, asset INT
, liquidity TEXT
);

CREATE TABLE IF NOT EXISTS dateseries
(day TEXT PRIMARY KEY
);
