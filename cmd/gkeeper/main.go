package main

import (
	"github.com/shulganew/GophKeeper/internal/api/jwt"
	"github.com/shulganew/GophKeeper/internal/api/oapi"
	"github.com/shulganew/GophKeeper/internal/api/router"
	"github.com/shulganew/GophKeeper/internal/app"
	"github.com/shulganew/GophKeeper/internal/app/config"
	"github.com/shulganew/GophKeeper/internal/services"
	"go.uber.org/zap"
)

func main() {

	// Get application config.
	conf := config.InitConfig()

	// Init application logging.
	app.InitLog()

	zap.S().Infoln("Starting GophKeeper...")

	// Root app context.
	ctx, cancel := app.InitContext()

	// Error channel.
	componentsErrs := make(chan error, 1)

	// Init Repo
	stor, err := app.InitStore(ctx, conf)
	if err != nil {
		zap.S().Fatalln(err)
	}

	// Init file storage
	fstor, err := app.InitMinIO(ctx, conf)
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

	key, err := jwt.GetPrivateKey(conf)
	if err != nil {
		zap.S().Fatalln(err)
	}

	// Create JWT authenticator.
	auth, err := jwt.NewUserAuthenticator(key)
	if err != nil {
		zap.S().Fatalln(err)
	}

	// Create router.
	rt := router.RouteShear(conf, swagger, auth)

	keeper := services.NewKeeper(ctx, stor, fstor, conf, auth)

	// We now register our GophKeeper above as the handler for the interface
	oapi.HandlerFromMux(keeper, rt)

	// Start web server.
	restDone := app.StartAPI(ctx, &conf, componentsErrs, rt)
	//Graceful shutdown.
	app.Graceful(ctx, cancel, componentsErrs)

	// Waiting http server shuting down.
	<-restDone

	zap.S().Infoln("App done.")
}
