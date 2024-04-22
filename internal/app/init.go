package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/jmoiron/sqlx"
	"github.com/shulganew/GophKeeper/internal/app/config"
	"github.com/shulganew/GophKeeper/internal/storage"
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

func InitApp(ctx context.Context, conf config.Config) (application *UseCases, err error) {

	// Connection for Keeper Database.
	db, err := sqlx.Connect(config.DataBaseType, conf.DSN)
	if err != nil {
		return nil, err
	}

	// Load storage.
	stor, err := storage.NewRepo(ctx, db)
	if err != nil {
		return nil, err
	}

	// Create config Container
	application = NewUseCases(conf, stor)

	zap.S().Infoln("Application init complite")
	return application, nil
}
