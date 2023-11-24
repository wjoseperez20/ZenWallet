-- migrate:up

-- Create the sequence
CREATE SEQUENCE seq_accounts_id START WITH 1;
CREATE SEQUENCE seq_accounts_number
    START WITH 10001 -- Starting value for the sequence
    INCREMENT BY 1 -- Increment by 1 for each new value
    MINVALUE 10001 -- Minimum value for the sequence
    MAXVALUE 99999 -- Maximum value for the sequence
    NO CYCLE;

-- Create the table
CREATE TABLE accounts
(
    id         integer                  NOT NULL DEFAULT nextval('seq_accounts_id'),
    client     varchar(255),
    email      varchar(255)             NOT NULL UNIQUE,
    number     integer                  NOT NULL DEFAULT nextval('seq_accounts_number') UNIQUE,
    balance    DECIMAL(10, 2)                    DEFAULT 0.0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    PRIMARY KEY (number)
);

-- migrate:down

-- Drop de the users table
DROP TABLE if exists accounts;

-- Drop the sequence
DROP SEQUENCE seq_accounts_number;
DROP SEQUENCE seq_accounts_id;
