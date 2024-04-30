package entities

import "github.com/gofrs/uuid"

type User struct {
	UUID     uuid.UUID `db:"user_id"`
	Login    string    `db:"login"`
	Password string    `db:"password"`
	PassHash string    `db:"password_hash"`
	Email    string    `db:"email"`
}
