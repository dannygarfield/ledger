CREATE TABLE IF NOT EXISTS transactions
(source TEXT
, destination TEXT
, happened_at TEXT
, amount TEXT
);

CREATE TABLE IF NOT EXISTS buckets
(name TEXT PRIMARY KEY
, asset INT
, liquidity TEXT
);

CREATE TABLE IF NOT EXISTS dateseries
(day TEXT PRIMARY KEY
);
