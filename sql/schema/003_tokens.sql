-- +goose Up
CREATE TABLE token (
	token TEXT 	PRIMARY KEY,
	created_at	TIMESTAMP NOT NULL,
	updated_at	TIMESTAMP NOT NULL,
	user_id 	UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	expire_at	TIMESTAMP NOT NULL,
	revoke_at 	TIMESTAMP DEFAULT NULL
);
-- +goose Down
DROP TABLE IF EXISTS token;
