package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/shulganew/GophKeeper/internal/entities"
)

// TODO add UNICQUE siteURL+login with error duplicated.
func (r *Repo) AddSite(ctx context.Context, site entities.Secret) error {
	query := `
	INSERT INTO secrets (user_id, type, data, key, uploaded) 
	VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.ExecContext(ctx, query, site.UserID, site.Stype.String(), site.Data, time.Now(), time.Now())
	if err != nil {
		return fmt.Errorf("error during add Site credentials, error: %w", err)
	}
	return nil
}
