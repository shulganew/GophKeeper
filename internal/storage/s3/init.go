package s3

import (
	"context"

	"github.com/minio/minio-go/v7"
)

type FileRepo struct {
	mio *minio.Client
}

func NewFileRepo(ctx context.Context, backet string, mio *minio.Client) (*FileRepo, error) {
	fr := FileRepo{mio: mio}
	_, err := fr.mio.BucketExists(ctx, backet)
	if err != nil {
		return nil, err
	}
	return &fr, nil
}
