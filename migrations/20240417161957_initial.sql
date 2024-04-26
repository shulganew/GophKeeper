-- +goose Up
-- +goose StatementBegin
DO $$
BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'processing') THEN
		CREATE TYPE secret_type AS ENUM ('SITE', 'CARD', 'TEXT', 'BIN');
	END IF;
END$$;

CREATE TABLE IF NOT EXISTS users (
	user_id UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(), 
	login TEXT NOT NULL UNIQUE, 
	password_hash TEXT NOT NULL,
	email TEXT NOT NULL
	);

CREATE TABLE IF NOT EXISTS secrets (
	secret_id UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
	definition TEXT NOT NULL DEFAULT 'Description of user data', 
	user_id UUID NOT NULL REFERENCES users(user_id), 
	type secret_type NOT NULL DEFAULT 'SITE',
	data BYTEA NOT NULL,
	key TIMESTAMPTZ NOT NULL,
	uploaded TIMESTAMPTZ NOT NULL
	);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
DROP TABLE secrets;
-- +goose StatementEnd
