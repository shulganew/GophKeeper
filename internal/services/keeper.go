// Main service level. Keeper inratface implements in storage.
//
// File secrets has main methods for crypto GophKeeper tasks.
// Key types, whitch GophKeeper use:
// mKey - master key, use for encoding ephemeral keys.
// eKey - Ephemeral key, use for dKey saving in database. store in separate table "ekeys", encoded by mKey.
// It stored opened in memory during service loading.
// dKey - Data key (use for data coding in table secretes. Saved in the same table, encoded by eKey
// Create Ephemeral key from master key, user time stamp as key id. Time stamp == key version.
// Postfix "c" in key name means crypted, for ex "eKeyc" mean ephemeral key enceded by master key, "dKeyc" = data key enceded by eKey.
package services

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/shulganew/GophKeeper/internal/app/config"
	"github.com/shulganew/GophKeeper/internal/entities"
	"github.com/shulganew/GophKeeper/internal/rest/oapi"
)

// User creation, registration, validation and autentification service.
type Keeper struct {
	stor  Keeperer
	conf  config.Config
	eKeys []entities.EKeyMem // Decoded ephemeral keys.
}

type Keeperer interface {
	AddUser(ctx context.Context, login, hash, email string) (userID *uuid.UUID, err error)
	GetByLogin(ctx context.Context, login string) (userID *entities.User, err error)
	AddSite(ctx context.Context, site entities.Secret) (siteID *uuid.UUID, err error)
	GetSites(ctx context.Context, userID string, stype entities.SecretType) (site []entities.Secret, err error)

	// Operations with keys.
	// Get Ephemeral encoded keys from storage
	SaveEKeysc(ctx context.Context, eKeysc []entities.EKeyDB) (err error) // Many keys
	SaveEKeyc(ctx context.Context, eKeyc entities.EKeyDB) (err error)     // One key
	LoadEKeysc(ctx context.Context) (eKeysc []entities.EKeyDB, err error)
}

var _ oapi.ServerInterface = (*Keeper)(nil)

func NewKeeper(ctx context.Context, stor Keeperer, conf config.Config) *Keeper {
	keeper := &Keeper{stor: stor, conf: conf, eKeys: []entities.EKeyMem{}}
	// Load eKeys.
	keeper.LoadKeyRing(ctx)
	return keeper
}
