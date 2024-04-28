package main

import (
	"github.com/shulganew/GophKeeper/internal/app"
	"github.com/shulganew/GophKeeper/internal/app/config"
	"github.com/shulganew/GophKeeper/internal/rest/oapi"
	"github.com/shulganew/GophKeeper/internal/rest/router"
	"github.com/shulganew/GophKeeper/internal/services"
	"go.uber.org/zap"
)

func main() {
	zap.S().Infoln("Starting GophKeeper...")

	// Get application config.
	conf := config.InitConfig()

	// Init application logging.
	app.InitLog()

	// Root app context.
	ctx, cancel := app.InitContext()

	// Error channel.
	componentsErrs := make(chan error, 1)

	// Init Repo
	stor, err := app.InitStore(ctx, conf)
	if err != nil {
		zap.S().Fatalln(err)
	}

	swagger, err := oapi.GetSwagger()
	if err != nil {
		zap.S().Fatalln(err)
	}
	// Clear out the servers array in the swagger spec, that skips validating
	// that server names match. We don't know how this thing will be run.
	swagger.Servers = nil

	// Create router.
	rt := router.RouteShear(conf, swagger)

	keeper := services.NewKeeper(ctx, stor, conf)

	// We now register our GophKeeper above as the handler for the interface
	oapi.HandlerFromMux(keeper, rt)

	// Start web server.
	restDone := app.StartREST(ctx, &conf, componentsErrs, rt)
	//Graceful shutdown.
	app.Graceful(ctx, cancel, componentsErrs)

	// Waiting http server shuting down.
	<-restDone

	zap.S().Infoln("App done.")
}
