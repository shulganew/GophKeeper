package storage

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/shulganew/GophKeeper/internal/entities"
)

// TODO add UNICQUE siteURL+login with error duplicated.
func (r *Repo) AddSite(ctx context.Context, site entities.NewSecretEncoded) (secretID *uuid.UUID, err error) {
	query := `
	INSERT INTO secrets (user_id, type, data, ekey_version, dkey, uploaded) 
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING secret_id
	`
	secretID = &uuid.UUID{}
	err = r.db.GetContext(ctx, secretID, query, site.UserID, site.Type, site.DataCr, site.EKeyVer, site.DKey, site.Uploaded)
	if err != nil {
		return nil, fmt.Errorf("db error during add Site credentials, error: %w", err)
	}
	return
}

// TODO add UNICQUE siteURL+login with error duplicated.
func (r *Repo) GetSites(ctx context.Context, userID string, stype entities.SecretType) (sites []entities.SecretEncoded, err error) {
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
