-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
	user_id UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(), 
	login TEXT NOT NULL UNIQUE, 
	password_hash TEXT NOT NULL
	);

CREATE TABLE IF NOT EXISTS sites_secrets (
	site_id UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(), 
	user_id UUID NOT NULL REFERENCES users(user_id), 
	site_url TEXT NOT NULL, 
	slogin TEXT NOT NULL,
	spw TEXT NOT NULL
	);

CREATE TABLE IF NOT EXISTS site_grand (
	owner_id UUID NOT NULL REFERENCES users(user_id),
	site_id UUID NOT NULL REFERENCES sites_secrets(site_id), 
	grand_id UUID NOT NULL REFERENCES users(user_id)
	);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
