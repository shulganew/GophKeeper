package main

import (
	"github.com/shulganew/GophKeeper/internal/app"
	"github.com/shulganew/GophKeeper/internal/app/config"
	"github.com/shulganew/GophKeeper/internal/rest/router"
	"go.uber.org/zap"
)

func main() {
	zap.S().Infoln("Starting GophKeeper...")
	// Init application logging.
	app.InitLog()

	// Root app context.
	ctx, cancel := app.InitContext()

	zap.S().Infoln("Hello passwor master!")

	conf := config.InitConfig()

	// Error channel.
	componentsErrs := make(chan error, 1)

	// Create router.
	rt := router.RouteShear(conf)

	// Start web server.
	restDone := app.StartREST(ctx, conf, componentsErrs, rt)

	//Graceful shutdown.
	app.Graceful(ctx, cancel, componentsErrs)

	// Waiting http server shuting down.
	<-restDone
	zap.S().Infoln("App done.")
}
