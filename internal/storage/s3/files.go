package s3

import (
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
)

// Updload file to MINIO storage.
func (r *FileRepo) UploadFile(ctx context.Context, fileID string, fr io.Reader) (err error) {
	zap.S().Debugln("Key Upload: ", fileID)
	// Put object to minio.
	_, err = r.mio.PutObject(ctx, "gohpkeeper", fileID, fr, int64(-1), minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		return fmt.Errorf("upload file to MIO problem: %w", err)
	}
	return
}

// Download file from MINIO storage.
func (r *FileRepo) DownloadFile(ctx context.Context, fileID string) (fr *minio.Object, err error) {
	// Put object to minio.
	zap.S().Debugln("Key Download: ", fileID)
	fr, err = r.mio.GetObject(ctx, "gohpkeeper", fileID, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("downlad file to MIO problem: %w", err)
	}
	return
}
