package app

import (
	"context"
	"fmt"
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
	"go.uber.org/zap/zapcore"
)

func InitLog(conf config.Config) zap.SugaredLogger {
	lvl, err := zap.ParseAtomicLevel(conf.ZapLevel)
	if err != nil {
		fmt.Println("Can't set log level: ", err, conf.ZapLevel)
		panic(err)
	}

	var op []string
	var ep []string
	if conf.RunLocal {
		op = []string{"stdout"}
		ep = []string{"stderr"}
	} else {
		op = []string{conf.ZapPath}
		ep = []string{conf.ZapPath}
	}

	cfg := zap.Config{
		Encoding:         "console",
		Level:            lvl,
		OutputPaths:      op,
		ErrorOutputPaths: ep,
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			TimeKey:    "time",
			EncodeTime: zapcore.RFC3339TimeEncoder,

			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	zapLogger := zap.Must(cfg.Build())
	zapLogger.Info("logger construction succeeded")
	zap.ReplaceGlobals(zapLogger)
	defer func() {
		_ = zapLogger.Sync()
	}()

	sugar := *zapLogger.Sugar()

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
		Creds:  credentials.NewStaticV4(conf.IDmi, conf.Secretmi, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}
	// Load file storage.
	fstor, err = s3.NewFileRepo(ctx, conf.Backetmi, mio)
	if err != nil {
		return nil, err
	}
	return
}
