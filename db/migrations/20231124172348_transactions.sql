-- migrate:up

-- Create the sequence
CREATE SEQUENCE seq_transactions_id START WITH 1;

-- Create the table
CREATE TABLE transactions
(
    id             integer                  NOT NULL DEFAULT nextval('seq_transactions_id'),
    amount         numeric                  NOT NULL DEFAULT 0,
    date           date                     NOT NULL,
    account_number integer                  NOT NULL,
    created_at     TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at     TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    PRIMARY KEY (id)
);
-- Alter table for foreign keys
ALTER TABLE transactions
    ADD CONSTRAINT fk_account_number FOREIGN KEY (account_number) REFERENCES accounts (number);

-- migrate:down

-- Drop de the users table
DROP TABLE if exists transactions;

-- Drop the sequence
DROP SEQUENCE seq_transactions_id;
