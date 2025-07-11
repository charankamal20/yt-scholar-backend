-- +goose Up
-- SQL in this section is executed when the migration is applied

CREATE TABLE users (
    user_id     VARCHAR PRIMARY KEY,
    name        VARCHAR NOT NULL,
    email       VARCHAR NOT NULL UNIQUE,
    profile_pic TEXT    NOT NULL,
    updated_at  TIMESTAMP,
    created_at  TIMESTAMP DEFAULT now()
);

CREATE TABLE refresh_tokens (
    user_id     VARCHAR NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    token       TEXT UNIQUE,
    time_issued TIMESTAMP DEFAULT now(),
    expire_at   TIMESTAMP NOT NULL,
    deleted     TIMESTAMP DEFAULT NULL,
    created_at  TIMESTAMP DEFAULT now(),
    PRIMARY KEY (user_id)
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back

DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS users;
