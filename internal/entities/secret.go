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
// !!! Each encodeted date type has own ID, for ex Gfile has gfileID, whitch equal secredtID in database. i.e. secredtID type FILE == gfileID
type NewSecret struct {
	UserID   string     `db:"user_id"`
	Type     SecretType `db:"type"` // Type of data - Site data, Credit card, Text or file.
	EKeyVer  time.Time  `db:"ekey_version"`
	DKeyCr   []byte     `db:"dkey"`
	Uploaded time.Time  `db:"uploaded"`
}
type SecretDecoded struct {
	NewSecret
	SecretID uuid.UUID
	Data     []byte
}

type SecretEncoded struct {
	NewSecret
	SecretID uuid.UUID `db:"secret_id"` // Stored secretID.
	DataCr   []byte    `db:"data"`      // Decrypted data.
}

type NewSecretDecoded struct {
	NewSecret
	Data []byte // Stored decod data.
}

// Data in struct with crypted data.
type NewSecretEncoded struct {
	NewSecret
	DataCr []byte `db:"data"` // Stored crypted data - data crypted.
}
