package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/shulganew/GophKeeper/internal/entities"
	"go.uber.org/zap"
)

// TODO add UNICQUE siteURL+login with error duplicated.
func (r *Repo) AddSite(ctx context.Context, site entities.Secret) (secretID *uuid.UUID, err error) {
	query := `
	INSERT INTO secrets (user_id, type, data, key, uploaded) 
	VALUES ($1, $2, $3, $4, $5)
	RETURNING secret_id
	`
	zap.S().Infof("Site data: %+v \n", site)
	secretID = &uuid.UUID{}
	err = r.db.GetContext(ctx, secretID, query, site.UserID, site.Stype.String(), site.Data, time.Now(), time.Now())
	if err != nil {
		return nil, fmt.Errorf("db error during add Site credentials, error: %w", err)
	}
	return
}

// TODO add UNICQUE siteURL+login with error duplicated.
func (r *Repo) GetSites(ctx context.Context, userID string, stype entities.SecretType) (sites []entities.Secret, err error) {
	query := `
	SELECT secret_id, description, data, key, uploaded
	FROM secrets 
	WHERE type = $1 AND user_id = $2
	ORDER BY uploaded DESC
	`
	err = r.db.SelectContext(ctx, &sites, query, stype, userID)
	if err != nil {
		return nil, fmt.Errorf("db error during getting list Site credentials, error: %w", err)
	}
	return
}
