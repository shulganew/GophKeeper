// Main service level. Keeper inratface implements in storage.
//
// File secrets has main methods for crypto GophKeeper tasks.
// Key types, whitch GophKeeper use:
// - mKey - master key, use for encoding ephemeral keys.
// - eKey - Ephemeral key, use for dKey saving in database. Store in separate table "ekeys", encoded by mKey.
// It stored opened in memory during service loading.
// - dKey - Data key (use for data coding in table secretes. Saved in the same table, encoded by eKey
// Create Ephemeral key from master key, user time stamp as key id. Time stamp == key version.
// Postfix "c" in key name means crypted, for ex "eKeyc" mean ephemeral key enceded by master key, "dKeyc" = data key enceded by eKey.
package services

import (
	"context"
	"io"

	"github.com/gofrs/uuid"
	"github.com/shulganew/GophKeeper/internal/api/jwt"
	"github.com/shulganew/GophKeeper/internal/api/oapi"
	"github.com/shulganew/GophKeeper/internal/app/config"
	"github.com/shulganew/GophKeeper/internal/entities"
	"github.com/shulganew/GophKeeper/internal/storage/s3"
)

// User creation, registration, validation and autentification service.
type Keeper struct {
	stor  Keeperer
	fstor FileKeeper
	conf  config.Config
	ua    *jwt.UserAuthenticator
	eKeys []entities.EKeyMem // Decoded ephemeral keys.
}

func NewKeeper(ctx context.Context, stor Keeperer, fstor FileKeeper, conf config.Config, ua *jwt.UserAuthenticator) *Keeper {
	keeper := &Keeper{stor: stor, fstor: fstor, conf: conf, eKeys: []entities.EKeyMem{}, ua: ua}
	// Load eKeys.
	keeper.LoadKeyRing(ctx)
	return keeper
}

type Keeperer interface {
	AddUser(ctx context.Context, login, hash, email, otpKey string) (userID *uuid.UUID, err error)
	GetByLogin(ctx context.Context, login string) (userID *entities.User, err error)

	// Entities credentials methods (site, card, text, file)
	AddSecretStor(ctx context.Context, entity entities.NewSecretEncoded, stype entities.SecretType) (siteID *uuid.UUID, err error)
	GetSecretsStor(ctx context.Context, userID string, stype entities.SecretType) (site []*entities.SecretEncoded, err error)
	GetSecretStor(ctx context.Context, secretID string) (site *entities.SecretEncoded, err error)
	UpdateSecretStor(ctx context.Context, entity entities.NewSecretEncoded, secretID string) (err error)
	DeleteSecretStor(ctx context.Context, secretID string) (err error)
	// Operations with keys.
	// Get Ephemeral encoded keys from storage
	SaveEKeysc(ctx context.Context, eKeysc []entities.EKeyDB) (err error) // Many keys
	SaveEKeyc(ctx context.Context, eKeyc entities.EKeyDB) (err error)     // One key
	LoadEKeysc(ctx context.Context) (eKeysc []entities.EKeyDB, err error)
	DropKeys(ctx context.Context) (err error)
}

type FileKeeper interface {
	UploadFile(ctx context.Context, backet string, fileID string, fr io.Reader) (err error)
	DownloadFile(ctx context.Context, backet string, fileID string) (fr io.ReadCloser, err error)
	DeleteFile(ctx context.Context, backet string, fileID string) (err error)
}

// Check interfaces.
var _ oapi.ServerInterface = (*Keeper)(nil)
var _ FileKeeper = (*s3.FileRepo)(nil)
