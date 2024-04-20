package router

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/shulganew/GophKeeper/internal/app"
	"github.com/shulganew/GophKeeper/internal/entities"
	"github.com/shulganew/GophKeeper/internal/rest/handler"
	"github.com/shulganew/GophKeeper/internal/rest/middlewares"
)

// Chi Router for application.
func RouteShear(application *app.UseCases) (r *chi.Mux) {
	conf := application.Config()
	r = chi.NewRouter()
	// Send password for enctription to middlewares.
	r.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), entities.CtxPassKey{}, conf.PassJWT)
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	r.Route("/", func(r chi.Router) {
		// User registration.
		r.Route("/api/user/auth", func(r chi.Router) {
			userReg := handler.NewHandlerRegister(conf, application.UserService())
			r.Post("/register", http.HandlerFunc(userReg.AddUser))
			userLogin := handler.NewHandlerLogin(conf, application.UserService())
			r.Post("/login", http.HandlerFunc(userLogin.LoginUser))

		})
		// Add auth/

		r.Route("/api/user/site", func(r chi.Router) {
			r.Use(middlewares.Auth)
			siteAdd := handler.NewSiteAdd(conf, application.SiteService())
			r.Post("/add", http.HandlerFunc(siteAdd.SiteAdd))
		})

	})
	return
}
