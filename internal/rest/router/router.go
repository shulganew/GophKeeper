package router

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/shulganew/GophKeeper/internal/app"
	"github.com/shulganew/GophKeeper/internal/entities"
	"github.com/shulganew/GophKeeper/internal/rest/handler"
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
		userReg := handler.NewHandlerRegister(conf, application.UserService())
		r.Post("/api/user/register", http.HandlerFunc(userReg.AddUser))

		userLogin := handler.NewHandlerLogin(conf, application.UserService())
		r.Post("/api/user/login", http.HandlerFunc(userLogin.LoginUser))

	})
	return
}
