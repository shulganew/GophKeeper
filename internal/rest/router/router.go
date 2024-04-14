package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/shulganew/GophKeeper/internal/app/config"
	"github.com/shulganew/GophKeeper/internal/rest/handler"
)

// Chi Router for application.
func RouteShear(conf *config.Config) (r *chi.Mux) {
	r = chi.NewRouter()
	v := handler.NewVersion(conf)
	r.Route("/", func(r chi.Router) {

		r.Get("/", http.HandlerFunc(v.GetVersion))

	})
	return
}
