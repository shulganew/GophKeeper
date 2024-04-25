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
	AddUser(ctx context.Context, login string, hash string) (*uuid.UUID, error)
	GetByLogin(ctx context.Context, login string) (*entities.User, error)
	AddSite(ctx context.Context, site entities.Secret) error
}

var _ oapi.ServerInterface = (*Keeper)(nil)

func NewKeeper(stor Keeperer, conf config.Config) *Keeper {
	return &Keeper{stor: stor, conf: conf}
}
