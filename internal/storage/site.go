package storage

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/shulganew/GophKeeper/internal/entities"
	"go.uber.org/zap"
)

// TODO add UNICQUE siteURL+login with error duplicated.
func (r *Repo) AddSite(ctx context.Context, site entities.Secret) (secretID *uuid.UUID, err error) {
	query := `
	INSERT INTO secrets (user_id, type, data, ekey_version, dkey, uploaded) 
	VALUES (:user_id, :type, :data, :ekey_version, :dkey, :uploaded)
	`
	zap.S().Infof("Site data: %+v \n", site)
	secretID = &uuid.UUID{}
	_, err = r.db.NamedExecContext(ctx, query, site)
	if err != nil {
		return nil, fmt.Errorf("db error during add Site credentials, error: %w", err)
	}
	return
}

// TODO add UNICQUE siteURL+login with error duplicated.
func (r *Repo) GetSites(ctx context.Context, userID string, stype entities.SecretType) (sites []entities.Secret, err error) {
	query := `
	SELECT secret_id, data, ekey_version, dkey, uploaded
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
