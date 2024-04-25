package entities

import (
	"time"

	"github.com/gofrs/uuid"
)

// Type of secter data.
type SecretType string

const (
	SITE SecretType = "SITE"
	CARD SecretType = "CARD"
	TEXT SecretType = "TEXT"
	BIN  SecretType = "BIN"
)

func (s *SecretType) String() string {
	return string(*s)
}

// DB DTO type for storing secter data.
type Secret struct {
	UUID   uuid.UUID  `db:"secret_id"`
	UserID string     `db:"user_id"`
	Stype  SecretType `db:"secret_type"`
	Data   []byte     `db:"data"`
	Key    time.Time  `db:"key"`
	Upload time.Time  `db:"upload"`
}
