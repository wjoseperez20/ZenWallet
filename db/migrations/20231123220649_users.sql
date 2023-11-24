-- migrate:up

-- Create the sequence
CREATE SEQUENCE seq_users_id START WITH 1;

-- Create the table
CREATE TABLE users
(
    id         integer                  NOT NULL DEFAULT nextval('seq_users_id'),
    username   varchar(255)             NOT NULL unique,
    password   varchar(255)             NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    PRIMARY KEY (id)
);

-- migrate:down

-- Drop de the users table
DROP TABLE if exists users;

-- Drop the sequence
DROP SEQUENCE seq_users_id;
