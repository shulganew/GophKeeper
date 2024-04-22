package handler

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"
	"github.com/shulganew/GophKeeper/internal/app"
	"github.com/shulganew/GophKeeper/internal/app/config"
	"github.com/shulganew/GophKeeper/internal/entities"
	"github.com/shulganew/GophKeeper/internal/services"
	"github.com/shulganew/GophKeeper/internal/services/mocks"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestUserlogin(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		login      string
		passLogin  string
		passDB     string
		statusCode int
	}{
		{
			name:       "Login user success",
			method:     http.MethodPost,
			login:      "user",
			passLogin:  "qwerty",
			passDB:     "qwerty",
			statusCode: http.StatusOK,
		},
		{
			name:       "Login user fail",
			method:     http.MethodPost,
			login:      "user",
			passLogin:  "qwerty",
			passDB:     "asdfg",
			statusCode: http.StatusUnauthorized,
		},
	}

	app.InitLog()
	conf := config.Config{}

	// Init application.
	conf.Address = "localhost:8080"
	conf.PassJWT = "JWTsecret"
	requestURL := "localhost:8080" + "/api/user/login"

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log("Test name: ", tt.name)
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repoUser := mocks.NewMockuserRepo(ctrl)
			userServ := services.NewUserService(repoUser)

			uuid, err := uuid.NewV7()
			assert.NoError(t, err)

			cPass, err := bcrypt.GenerateFromPassword([]byte(tt.passDB), bcrypt.DefaultCost)
			assert.NoError(t, err)

			dbUser := entities.User{UUID: uuid, Login: tt.login, PassHash: string(cPass)}

			_ = repoUser.EXPECT().
				GetByLogin(gomock.Any(), gomock.Any()).
				AnyTimes().
				Return(&dbUser, nil)

			loginUser := entities.User{UUID: uuid, Login: tt.login, Password: tt.passLogin}

			jsonWs, err := json.Marshal(loginUser)
			if err != nil {
				log.Fatalln(err)
			}

			body := strings.NewReader(string(jsonWs))
			assert.NoError(t, err)

			// Add chi context.
			rctx := chi.NewRouteContext()
			req := httptest.NewRequest(http.MethodPost, requestURL, body)

			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			req.Header.Add("Content-Type", "application/json")

			// Create status recorder.
			resRecord := httptest.NewRecorder()

			// Make request.
			userLogin := NewHandlerLogin(conf, userServ)
			userLogin.LoginUser(resRecord, req)

			// Get result.
			res := resRecord.Result()

			b, _ := io.ReadAll(res.Body)

			t.Log(string(b))

			err = res.Body.Close()
			assert.NoError(t, err)

			// Check answer code.
			t.Log("StatusCode test: ", tt.statusCode, " server: ", res.StatusCode)

			assert.Equal(t, tt.statusCode, res.StatusCode)
		})
	}
}
