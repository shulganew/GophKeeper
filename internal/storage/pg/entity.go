package pg

import (
	"context"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/shulganew/GophKeeper/internal/entities"
)

// TODO add UNICQUE siteURL+login with error duplicated.

func (r *Repo) AddSecretStor(ctx context.Context, entity entities.NewSecretEncoded, stype entities.SecretType) (secretID *uuid.UUID, err error) {
	query := `
	INSERT INTO secrets (user_id, type, data, ekey_version, dkey, uploaded) 
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING secret_id
	`
	secretID = &uuid.UUID{}
	err = r.db.GetContext(ctx, secretID, query, entity.UserID, stype, entity.DataCr, entity.EKeyVer, entity.DKeyCr, entity.Uploaded)
	if err != nil {
		return nil, fmt.Errorf("db error during add Site credentials, error: %w", err)
	}
	return
}

// TODO add UNICQUE siteURL+login with error duplicated.
func (r *Repo) GetSecretsStor(ctx context.Context, userID string, stype entities.SecretType) (sites []entities.SecretEncoded, err error) {
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

// Get secret particular type by id.
func (r *Repo) GetSecretStor(ctx context.Context, secretID string) (sectet *entities.SecretEncoded, err error) {
	query := `
	SELECT secret_id, data, ekey_version, dkey, uploaded
	FROM secrets 
	WHERE secret_id = $1
	ORDER BY uploaded DESC
	`
	se := entities.SecretEncoded{}
	err = r.db.GetContext(ctx, &se, query, secretID)
	if err != nil {
		return nil, fmt.Errorf("db error during getting list Site credentials, error: %w", err)
	}
	return &se, nil
}

func (r *Repo) UpdateSecretStor(ctx context.Context, entity entities.NewSecretEncoded, secretID string) (err error) {
	query := `
	UPDATE secrets 
	SET  data = $1, ekey_version = $2, dkey = $3, uploaded = $4
	WHERE secret_id = $5
	`

	_, err = r.db.ExecContext(ctx, query, entity.DataCr, entity.EKeyVer, entity.DKeyCr, time.Now(), secretID)
	if err != nil {
		return fmt.Errorf("db error during getting list Site credentials, error: %w", err)
	}

	return nil
}
func (r *Repo) DeleteSecretStor(ctx context.Context, secretID string) (err error) {
	query := `
	DELETE from secrets 
	WHERE secret_id = $1
	`

	_, err = r.db.ExecContext(ctx, query, secretID)
	if err != nil {
		return fmt.Errorf("db error during getting list Site credentials, error: %w", err)
	}

	return nil
}
