package router

import (
	"context"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
	middleware "github.com/oapi-codegen/nethttp-middleware"
	"github.com/shulganew/GophKeeper/internal/app/config"
	"github.com/shulganew/GophKeeper/internal/entities"
	"github.com/shulganew/GophKeeper/internal/rest/middlewares"
)

// Chi Router for application.
func RouteShear(conf config.Config, swagger *openapi3.T) (r *chi.Mux) {
	r = chi.NewRouter()
	// Use our validation middleware to check all requests against the
	// OpenAPI schema.

	// Send password for enctription to middlewares.
	r.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), entities.CtxPassKey{}, conf.PassJWT)
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	})
	r.Use(middlewares.MidlewZip)
	r.Use(middlewares.Auth)
	r.Use(middleware.OapiRequestValidator(swagger))

	return
}
