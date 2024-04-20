package app

import (
	"github.com/shulganew/GophKeeper/internal/app/config"
	"github.com/shulganew/GophKeeper/internal/services"
	"github.com/shulganew/GophKeeper/internal/storage"
)

// A container pattern.
type UseCases struct {
	stor    *storage.Repo
	conf    *config.Config
	userSrv *services.UserService
	siteSrv *services.SiteService
}

func NewUseCases(conf *config.Config, stor *storage.Repo) *UseCases {
	cases := &UseCases{}
	cases.conf = conf
	cases.userSrv = services.NewUserService(stor)
	cases.siteSrv = services.NewSiteService(stor)
	cases.stor = stor
	return cases
}

func (c *UseCases) UserService() *services.UserService {
	return c.userSrv
}

func (c *UseCases) SiteService() *services.SiteService {
	return c.siteSrv
}

func (c *UseCases) Config() *config.Config {
	return c.conf
}

func (c *UseCases) Repo() *storage.Repo {
	return c.stor
}
