-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS users
(
    id             uuid primary key      not null default uuid_generate_v4(),
    email          varchar unique        not null,
    verified_email bool    default false not null,
    login          varchar unique        not null,
    password       varchar               not null,
    picture        varchar default null  null
);

CREATE INDEX IF NOT EXISTS users_login_idx ON users USING hash (login);
CREATE INDEX IF NOT EXISTS users_email_idx ON users USING hash (email);

CALL register_updated_at_created_at_columns('users')

-- +goose StatementEnd
