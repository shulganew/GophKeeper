package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/shulganew/GophKeeper/internal/app/config"
	"go.uber.org/zap"
)

const timeoutServerShutdown = time.Second * 5
const timeoutShutdown = time.Second * 20

// Manage web server.
func StartAPI(ctx context.Context, conf *config.Config, componentsErrs chan error, r *chi.Mux) (restDone chan struct{}) {
	// Start web server.
	var srv = http.Server{Addr: conf.Address, Handler: r}
	go func() {
		if err := srv.ListenAndServeTLS(conf.TLSPub, conf.TLSPK); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}
			componentsErrs <- fmt.Errorf("listen and server has failed: %w", err)
		}
	}()

	// Graceful shutdown.
	restDone = make(chan struct{})
	go func() {
		defer zap.S().Infoln("Server web has been graceful shutdown.")
		defer close(restDone)
		<-ctx.Done()
		shutdownTimeoutCtx, cancelShutdownTimeoutCtx := context.WithTimeout(context.Background(), timeoutServerShutdown)
		defer cancelShutdownTimeoutCtx()
		if err := srv.Shutdown(shutdownTimeoutCtx); err != nil {
			zap.S().Infoln("an error occurred during server shutdown: %v", err)
		}
	}()
	return
}

// Block gorutine until context done or get errer from componentsErrs channel.
func Graceful(ctx context.Context, cancel context.CancelFunc, componentsErrs chan error) {
	// Timer hardreset shutdown.
	context.AfterFunc(ctx, func() {
		ctx, cancelCtx := context.WithTimeout(context.Background(), timeoutShutdown)
		defer cancelCtx()
		<-ctx.Done()
		zap.S().Fatalln("failed to gracefully shutdown the service")
	})

	// Graceful shutdown.
	select {
	// Exit on root context done.
	case <-ctx.Done():
	// Exit on errors.
	case err := <-componentsErrs:
		zap.S().Errorln("Get server error: ", err)
		cancel()
	}
}
