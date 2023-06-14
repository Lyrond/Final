CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    name text NOT NULL,
    email text NOT NULL,
    password_hash text NOT NULL
);

ALTER TABLE users ADD CONSTRAINT users_email_unique UNIQUE (email);