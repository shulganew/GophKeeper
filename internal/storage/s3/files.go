package s3

import (
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
)

func (r *FileRepo) UploadFile(ctx context.Context, fileID string, fr io.ReadCloser) (err error) {

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
