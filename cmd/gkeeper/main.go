package main

import (
	"github.com/shulganew/GophKeeper/internal/app"
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

	// Error channel.
	componentsErrs := make(chan error, 1)

	a, err := app.InitApp(ctx)
	if err != nil {
		panic(err)
	}
	// Create router.
	rt := router.RouteShear(a)

	// Start web server.
	restDone := app.StartREST(ctx, a.Config(), componentsErrs, rt)
	//Graceful shutdown.
	app.Graceful(ctx, cancel, componentsErrs)

	// Waiting http server shuting down.
	<-restDone
	zap.S().Infoln("App done.")
}
