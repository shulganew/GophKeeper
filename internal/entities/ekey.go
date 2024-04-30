package entities

import (
	"time"
)

type EKeyMem struct {
	EKey []byte
	TS   time.Time
}

func NewEKeyMem(key []byte, ts time.Time) *EKeyMem {
	return &EKeyMem{EKey: key, TS: ts}
}

type EKeyDB struct {
	EKeyc []byte    `db:"ekeyc"`
	TS   time.Time `db:"ts"`
}

func NewEKeyDB(key []byte, ts time.Time) *EKeyDB {
	return &EKeyDB{EKeyc: key, TS: ts}
}
