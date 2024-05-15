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
	FILE SecretType = "BIN"
)

func (s *SecretType) String() string {
	return string(*s)
}

// DB DTO type for storing secter data.
// OAPI pattern - new mean struct without id (new constructor),
// id will be retruned by DB
// !!! Each encodeted date type has own ID, for ex Gfile has gfileID, whitch equal secredtID in database. For ex secredtID type FILE == gfileID.
type NewSecret struct {
	Type     SecretType `db:"type"` // Type of data - Site data, Credit card, Text or file.
	EKeyVer  time.Time  `db:"ekey_version"`
	Uploaded time.Time  `db:"uploaded"`
	UserID   string     `db:"user_id"`
	DKeyCr   []byte     `db:"dkey"`
}
type SecretDecoded struct {
	NewSecret
	Data     []byte
	SecretID uuid.UUID
}

type SecretEncoded struct {
	NewSecret
	DataCr   []byte    `db:"data"`      // Decrypted data.
	SecretID uuid.UUID `db:"secret_id"` // Stored secretID.
}

type NewSecretDecoded struct {
	NewSecret
	Data []byte // Stored decod data.
}

// NewSecretEncoded data in struct with crypted data.
type NewSecretEncoded struct {
	NewSecret
	DataCr []byte `db:"data"` // Stored crypted data - data crypted.
}
