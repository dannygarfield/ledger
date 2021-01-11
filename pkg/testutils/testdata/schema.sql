CREATE TABLE entries
(   source TEXT,
    destination TEXT,
    happened_at TEXT,
    amount TEXT
);

CREATE TABLE budget_entries
(
    happened_at TEXT,
    amount INT,
    category TEXT,
    description TEXT
);

INSERT INTO budget_entries
    (happened_at, amount, category, description)
VALUES
    (date("2021-01-01"), 3000, "rent", "-"),
    (date("2021-01-01"), 100, "groceries", "whole foods delivery"),
    (date("2021-01-02"), 200, "groceries", "food train")
;
