package entities

import (
	"github.com/gofrs/uuid"
)

type User struct {
	UUID     uuid.UUID `json:"-" db:"user_id"`
	Login    string    `json:"login" db:"login"`
	Password string    `json:"password"`
	PassHash string    `db:"password_hash"`
}
