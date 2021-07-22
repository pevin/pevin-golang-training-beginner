CREATE TABLE IF NOT EXISTS inquiries(
    id  VARCHAR (255) PRIMARY KEY,
    transaction_id VARCHAR (255) NOT NULL,
    payment_code VARCHAR (255) NOT NULL,
    amount decimal(16,4),
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL
);
