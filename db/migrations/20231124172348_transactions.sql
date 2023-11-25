-- migrate:up

-- Create the sequence
CREATE SEQUENCE seq_transactions_id START WITH 1;

-- Create the table
CREATE TABLE transactions
(
    id         integer                  NOT NULL DEFAULT nextval('seq_transactions_id'),
    amount     DECIMAL(10, 2)           NOT NULL,
    date       date                     NOT NULL,
    account_id integer                  NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    PRIMARY KEY (id)
);
-- Alter table for foreign keys
ALTER TABLE transactions
    ADD CONSTRAINT fk_account_number FOREIGN KEY (account_id) REFERENCES accounts (account);

-- migrate:down

-- Drop de the users table
DROP TABLE if exists transactions;

-- Drop the sequence
DROP SEQUENCE seq_transactions_id;
