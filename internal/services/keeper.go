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
	stor Keeperer
	conf config.Config
}

type Keeperer interface {
	AddUser(ctx context.Context, login string, hash string) (userID *uuid.UUID, err error)
	GetByLogin(ctx context.Context, login string) (userID *entities.User, err error)
	AddSite(ctx context.Context, site entities.Secret) (siteID *uuid.UUID, err error)
	GetSites(ctx context.Context, userID string, stype entities.SecretType) (site []entities.Secret, err error)
}

var _ oapi.ServerInterface = (*Keeper)(nil)

func NewKeeper(stor Keeperer, conf config.Config) *Keeper {
	return &Keeper{stor: stor, conf: conf}
}
