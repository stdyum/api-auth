-- +goose Up
-- +goose StatementBegin

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE OR REPLACE FUNCTION update_updated_at_column()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE PROCEDURE register_updated_at_created_at_columns(table_name varchar) AS
$$
BEGIN
    EXECUTE format(
            'ALTER TABLE %I ADD COLUMN created_at timestamp not null default now();',
            table_name
            );

    EXECUTE format(
            'ALTER TABLE %I ADD COLUMN updated_at timestamp not null default now();',
            table_name
            );

    EXECUTE format(
            'CREATE TRIGGER %I BEFORE UPDATE ON %I FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();',
            table_name || '_updated_at_trg', table_name
            );

END;
$$ LANGUAGE plpgsql;


-- +goose StatementEnd
