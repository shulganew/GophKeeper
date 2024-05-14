package entities

import (
	"time"
)

type EKeyMem struct {
	TS   time.Time
	EKey []byte
}

func NewEKeyMem(key []byte, ts time.Time) *EKeyMem {
	return &EKeyMem{EKey: key, TS: ts}
}

type EKeyDB struct {
	TS    time.Time `db:"ts"`
	EKeyc []byte    `db:"ekeyc"`
}

func NewEKeyDB(key []byte, ts time.Time) *EKeyDB {
	return &EKeyDB{EKeyc: key, TS: ts}
}
