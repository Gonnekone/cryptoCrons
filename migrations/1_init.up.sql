CREATE TABLE IF NOT EXISTS coins
(
    id   TEXT PRIMARY KEY,
    name VARCHAR(20) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS prices
(
    id    TEXT NOT NULL,
    ts    INT  NOT NULL,
    price INT  NOT NULL,
    CONSTRAINT fk_coins FOREIGN KEY (id) REFERENCES coins (id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_prices_timestamp ON prices (ts);
