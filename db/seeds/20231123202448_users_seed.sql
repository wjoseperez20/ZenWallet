-- migrate:up

-- Populate the database with some sample data
INSERT INTO users(username, password)
VALUES ('wjoseperez', '$2a$14$o713Q3nBKrXAvITz0WDy5O3xBXH5N4CBGFaEBaWYyPVYlInUm7Wo6');

-- migrate:down

-- Delete sample data
DELETE
* FROM users;



