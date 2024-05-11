package services

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"
	"github.com/oapi-codegen/testutil"
	"github.com/shulganew/GophKeeper/internal/api/jwt"
	"github.com/shulganew/GophKeeper/internal/api/oapi"
	"github.com/shulganew/GophKeeper/internal/api/router"
	"github.com/shulganew/GophKeeper/internal/app"
	"github.com/shulganew/GophKeeper/internal/app/config"
	"github.com/shulganew/GophKeeper/internal/entities"
	"github.com/shulganew/GophKeeper/internal/services/mocks"
	"go.uber.org/zap"

	"github.com/stretchr/testify/require"
)

func TestUser(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		register       string
		login          string
		userRegister   oapi.NewUser
		userLogin      oapi.NewUser
		statusRegister int
		statusLogin    int
	}{
		{
			name:           "Check user registration and login OK",
			register:       "/auth/register",
			login:          "/auth/login",
			userRegister:   oapi.NewUser{Email: "me@ya.ru", Login: "user", Password: "123"},
			userLogin:      oapi.NewUser{Email: "me@ya.ru", Login: "user", Password: "123"}, // same user: login and registered
			statusRegister: http.StatusCreated,
			statusLogin:    http.StatusOK,
		},
		{
			name:           "Check user registration and login bad auth 401",
			register:       "/auth/register",
			login:          "/auth/login",
			userRegister:   oapi.NewUser{Email: "me@ya.ru", Login: "user", Password: "123"},
			userLogin:      oapi.NewUser{Email: "me@ya.ru", Login: "user", Password: "456"}, // bad pass
			statusRegister: http.StatusCreated,
			statusLogin:    http.StatusUnauthorized,
		},
	}

	app.InitLog()
	ctx := context.Background()

	conf := &config.Config{}

	// Init application
	conf.Address = "localhost:8088"
	conf.PathJWT = "JWTsecret"
	conf.MasterKey = "MasterPw"
	conf.DSN = ""
	conf.PathJWT = "cert/jwtpkey.pem"
	// init mock.

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// crete mock storege
	repo := mocks.NewMockKeeperer(ctrl)

	// Load keys on startup.
	_ = repo.EXPECT().
		LoadEKeysc(ctx).
		AnyTimes().
		Return([]entities.EKeyDB{}, nil)

	_ = repo.EXPECT().
		SaveEKeyc(ctx, gomock.Any()).
		AnyTimes().
		Return(nil)

	// Create JWT authenticator.
	auth, err := jwt.NewUserAuthenticator([]byte(JWTPemKey))
	if err != nil {
		zap.S().Fatalln(err)
	}

	// Init web.
	// Get the swagger description of our API
	swagger, err := oapi.GetSwagger()
	require.NoError(t, err)

	// Clear out the servers array in the swagger spec, that skips validating
	// that server names match. We don't know how this thing will be run.
	swagger.Servers = nil

	// Create router.
	rt := router.RouteShear(*conf, swagger, auth)

	keeper := NewKeeper(ctx, repo, nil, *conf, auth)

	// We now register our GophKeeper above as the handler for the interface
	oapi.HandlerFromMux(keeper, rt)

	userStorage := make(map[string]entities.User, 0)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log("Test name: ", tt.name)
			userID, err := uuid.NewV7()
			require.NoError(t, err)

			_ = repo.EXPECT().
				AddUser(gomock.Any(), tt.userRegister.Login, gomock.Any(), tt.userRegister.Email).
				DoAndReturn(func(_ any, login string, hash string, email string) (*uuid.UUID, error) {
					userStorage[login] = entities.User{UUID: userID, Login: login, PassHash: hash, Email: email}
					return &userID, nil
				}).
				AnyTimes()

			jsonSite, err := json.Marshal(tt.userRegister)
			require.NoError(t, err)

			// Create request.
			rr := testutil.NewRequest().Post(tt.register).WithContentType("application/json").WithBody(jsonSite).GoWithHTTPHandler(t, rt).Recorder
			require.Equal(t, tt.statusRegister, rr.Code)

			// Check user login.
			_ = repo.EXPECT().
				GetByLogin(gomock.Any(), tt.userRegister.Login).
				DoAndReturn(func(ctx any, login string) (*entities.User, error) {
					user, ok := userStorage[login]
					if ok {
						t.Log("User found in db.", login)
						return &user, nil
					}

					return nil, errors.New("User not found")
				}).
				AnyTimes()

			// Create request.
			jsonSite, err = json.Marshal(tt.userLogin)
			require.NoError(t, err)
			rr = testutil.NewRequest().Post(tt.login).WithContentType("application/json").WithBody(jsonSite).GoWithHTTPHandler(t, rt).Recorder
			require.Equal(t, tt.statusLogin, rr.Code)
		})
	}
}
