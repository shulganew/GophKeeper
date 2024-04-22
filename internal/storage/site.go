package storage

import (
	"context"
	"fmt"
)

// TODO add UNICQUE siteURL+login with error duplicated.
func (r *Repo) AddSite(ctx context.Context, userID, site, slogin, spw string) error {
	query := `
	INSERT INTO sites_secrets (user_id, site_url, slogin, spw) 
	VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.ExecContext(ctx, query, userID, site, slogin, spw)
	if err != nil {
		return fmt.Errorf("error during add Site credentials, error: %w", err)
	}
	return nil
}
