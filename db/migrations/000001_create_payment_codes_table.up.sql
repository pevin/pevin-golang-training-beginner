CREATE TABLE IF NOT EXISTS payment_codes(
   id  VARCHAR (255) PRIMARY KEY,
   payment_code VARCHAR (255) NOT NULL,
   name VARCHAR (255) NOT NULL,
   status VARCHAR (255) NOT NULL,
   expiration_date VARCHAR (255) NOT NULL,
  created_at timestamptz NOT NULL,
  updated_at timestamptz NOT NULL
);
