package s3

import (
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
)

// Updload file to MINIO storage.
func (r *FileRepo) UploadFile(ctx context.Context, backet string, fileID string, fr io.Reader) (err error) {
	zap.S().Debugln("Key Upload: ", fileID)
	// Put object to minio.
	_, err = r.mio.PutObject(ctx, backet, fileID, fr, int64(-1), minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		return fmt.Errorf("upload file to MIO problem: %w", err)
	}
	return
}

// Download file from MINIO storage.
func (r *FileRepo) DownloadFile(ctx context.Context, backet string, fileID string) (fr io.ReadCloser, err error) {
	// Put object to minio.
	zap.S().Debugln("Key Download: ", fileID)
	fr, err = r.mio.GetObject(ctx, backet, fileID, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("downlad file to MIO problem: %w", err)
	}
	return
}
