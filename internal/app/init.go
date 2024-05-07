package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/jmoiron/sqlx"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/shulganew/GophKeeper/internal/app/config"
	"github.com/shulganew/GophKeeper/internal/storage/pg"
	"github.com/shulganew/GophKeeper/internal/storage/s3"

	"go.uber.org/zap"
)

func InitLog() zap.SugaredLogger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	zap.ReplaceGlobals(logger)
	defer func() {
		_ = logger.Sync()
	}()

	sugar := *logger.Sugar()

	defer func() {
		_ = sugar.Sync()
	}()
	return sugar
}

// Init context from graceful shutdown. Send to all function for return by
//
//	syscall.SIGINT, syscall.SIGTERM
func InitContext() (ctx context.Context, cancel context.CancelFunc) {
	ctx, cancel = signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	return
}

func InitStore(ctx context.Context, conf config.Config) (stor *pg.Repo, err error) {
	// Connection for Keeper Database.
	db, err := sqlx.Connect(config.DataBaseType, conf.DSN)
	if err != nil {
		return nil, err
	}

	// Load storage.
	stor, err = pg.NewRepo(ctx, db)
	if err != nil {
		return nil, err
	}

	zap.S().Infoln("Application init complite")
	return stor, nil
}

func InitMinIO(ctx context.Context, conf config.Config) (fstor *s3.FileRepo, err error) {
	// Connection for Keeper MINIO.
	mio, err := minio.New("localhost:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("GYD6J3FzdOVG49aB6Ycb", "Ms37ZNWDG9CLFQhW92tA36NfgZa1zgy0z76bVmIJ", ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}
	// Load file storage.
	fstor, err = s3.NewFileRepo(ctx, conf.Backet, mio)
	if err != nil {
		return nil, err
	}
	return
}
