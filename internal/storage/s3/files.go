package s3

import (
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
)

// Updload file to MINIO storage.
func (r *FileRepo) UploadFile(ctx context.Context, fileID string, fr io.ReadCloser) (err error) {
	zap.S().Debugln("Key Upload: ", fileID)
	// Put object to minio.
	_, err = r.mio.PutObject(ctx, "gohpkeeper", fileID, fr, int64(-1), minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		return fmt.Errorf("upload file to MIO problem: %w", err)
	}
	err = fr.Close()
	if err != nil {
		return fmt.Errorf("upload file to MIO problem with close: %w", err)
	}
	return
}

// Download file from MINIO storage.
func (r *FileRepo) DownloadFile(ctx context.Context, storageID string) (fr *minio.Object, err error) {
	// Put object to minio.
	zap.S().Debugln("Key Download: ", storageID)
	fr, err = r.mio.GetObject(ctx, "gohpkeeper", storageID, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("downlad file to MIO problem: %w", err)
	}
	return
}
