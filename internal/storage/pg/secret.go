package pg


import (
	"context"
	"fmt"

	"github.com/shulganew/GophKeeper/internal/entities"
)

// Save to storage ephemeral encoded keys.
func (r *Repo) SaveEKeysc(ctx context.Context, eKeysc []entities.EKeyDB) (err error) {
	query := `INSERT INTO ekeys (ts, ekeyc) VALUES (:ts, :ekeyc) ON CONFLICT DO NOTHING`
	_, err = r.db.NamedExecContext(ctx, query, eKeysc)
	if err != nil {
		return fmt.Errorf("db error during saving eKeysc, error: %w", err)
	}
	return
}

func (r *Repo) SaveEKeyc(ctx context.Context, eKeyc entities.EKeyDB) (err error) {
	query := `INSERT INTO ekeys (ts, ekeyc) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err = r.db.ExecContext(ctx, query, eKeyc.TS, eKeyc.EKeyc)
	if err != nil {
		return fmt.Errorf("db error during saving key eKeyc, error: %w", err)
	}
	return
}

// Save to storage ephemeral encoded keys.
func (r *Repo) LoadEKeysc(ctx context.Context) (eKeysc []entities.EKeyDB, err error) {
	eKeysc = []entities.EKeyDB{}
	query := `SELECT * FROM ekeys`
	err = r.db.SelectContext(ctx, &eKeysc, query)
	if err != nil {
		return nil, fmt.Errorf("db error during loading eKeysc, error: %w", err)
	}
	return
}
