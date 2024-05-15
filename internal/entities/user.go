package entities

import "github.com/gofrs/uuid"

type User struct {
	Login    string    `db:"login"`
	Password string    `db:"password"`
	PassHash string    `db:"password_hash"`
	Email    string    `db:"email"`
	UUID     uuid.UUID `db:"user_id"`
}
