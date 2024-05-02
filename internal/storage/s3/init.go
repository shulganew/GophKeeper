package s3

import (
	"context"

	"github.com/minio/minio-go/v7"
)

type FileRepo struct {
	mio *minio.Client
}

func NewFileRepo(ctx context.Context, mio *minio.Client) (*FileRepo, error) {
	fr := FileRepo{mio: mio}
	// TODO add confilt param.
	_, err := fr.mio.BucketExists(ctx, "gohpkeeper")
	if err != nil {
		return nil, err
	}
	return &fr, nil
}
