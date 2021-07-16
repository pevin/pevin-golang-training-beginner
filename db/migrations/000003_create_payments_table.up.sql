CREATE TABLE IF NOT EXISTS payments(
    id  VARCHAR (255) PRIMARY KEY,
    transaction_id VARCHAR (255) NOT NULL,
    payment_code VARCHAR (255) NOT NULL,
    name VARCHAR (255) NOT NULL,
    amount decimal(16,4),
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL
);
