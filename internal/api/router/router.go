package router

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/go-chi/chi/v5"
	middleware "github.com/oapi-codegen/nethttp-middleware"
	"github.com/shulganew/GophKeeper/internal/api/jwt"
	"github.com/shulganew/GophKeeper/internal/app/config"
)

// Chi Router for application.
func RouteShear(conf config.Config, swagger *openapi3.T, v jwt.JWSValidator) (r *chi.Mux) {
	r = chi.NewMux()
	// Use our validation middleware to check all requests against the
	// OpenAPI schema.
	r.Use(middleware.OapiRequestValidatorWithOptions(swagger,
		&middleware.Options{
			Options: openapi3filter.Options{
				AuthenticationFunc: jwt.NewAuthenticator(v),
			},
		}))
	return
}
