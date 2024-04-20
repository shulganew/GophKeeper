package services

import (
	"context"

	"go.uber.org/zap"
)

// Work with saved site credentianls.
type SiteService struct {
	stor siteRepo
}

type siteRepo interface {
	AddSite(ctx context.Context, userID, site, slogin, spw string) error
}

func NewSiteService(stor siteRepo) *SiteService {
	return &SiteService{stor: stor}
}

// Add new site credential: site, login and password.
func (r *SiteService) AddSite(ctx context.Context, userID, site, slogin, spw string) (err error) {
	err = r.stor.AddSite(ctx, userID, site, slogin, spw)
	if err != nil {
		zap.S().Errorln("Error adding site credentials: ", err)
	}
	return
}

