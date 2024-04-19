-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
	id SERIAL, 
	user_id UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(), 
	login TEXT NOT NULL UNIQUE, 
	password_hash TEXT NOT NULL
	);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
