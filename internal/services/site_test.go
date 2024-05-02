package services

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"
	middleware "github.com/oapi-codegen/nethttp-middleware"
	"github.com/oapi-codegen/testutil"
	"github.com/shulganew/GophKeeper/internal/api/middlewares"
	"github.com/shulganew/GophKeeper/internal/api/oapi"
	"github.com/shulganew/GophKeeper/internal/app"
	"github.com/shulganew/GophKeeper/internal/app/config"
	"github.com/shulganew/GophKeeper/internal/entities"
	"github.com/shulganew/GophKeeper/internal/services/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSite(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		method     string
		requestAdd string
		hasJWT     bool
		status     int
	}{
		{
			name:       "Check site add and site list methods",
			method:     http.MethodPost,
			requestAdd: "/api/user/site/add",
			hasJWT:     true,
			status:     http.StatusCreated,
		},
		{
			name:       "Check user jwt",
			method:     http.MethodPost,
			requestAdd: "/api/user/site/add",
			hasJWT:     false,
			status:     http.StatusUnauthorized,
		},
	}

	app.InitLog()
	ctx := context.Background()

	conf := &config.Config{}

	// Init application
	conf.Address = "localhost:8088"
	conf.PassJWT = "JWTsecret"
	conf.MasterKey = "MasterPw"
	conf.DSN = ""
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

	keeper := NewKeeper(ctx, repo, nil, *conf)

	// Init web.
	// Get the swagger description of our API
	swagger, err := oapi.GetSwagger()
	require.NoError(t, err)

	// Clear out the servers array in the swagger spec, that skips validating
	// that server names match. We don't know how this thing will be run.
	swagger.Servers = nil

	// Create router.
	r := chi.NewRouter()

	// Send password for enctription to middlewares.
	r.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), entities.CtxPassKey{}, conf.PassJWT)
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	})
	r.Use(middlewares.Auth)
	r.Use(middleware.OapiRequestValidator(swagger))

	// We now register our GophKeeper above as the handler for the interface
	oapi.HandlerFromMux(keeper, r)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log("Test name: ", tt.name)
			userID, err := uuid.NewV7()
			require.NoError(t, err)

			secret_id, err := uuid.NewV7()
			require.NoError(t, err)
			//
			nsite := oapi.NewSite{Definition: "mysite", Site: "www.ru", Slogin: "igor", Spw: "123"}
			_ = repo.EXPECT().
				AddSite(gomock.Any(), gomock.Any()).
				AnyTimes().
				Return(&secret_id, nil)

			jsonSite, err := json.Marshal(nsite)
			require.NoError(t, err)

			//body := strings.NewReader(jsonSite)
			assert.NoError(t, err)

			// Add jwt to header.
			jwt := ""
			if tt.hasJWT {
				jwt, _ = middlewares.BuildJWTString(userID, conf.PassJWT)
			}
			// Create request.
			rr := testutil.NewRequest().Post(tt.requestAdd).WithContentType("application/json").WithBody(jsonSite).WithHeader("Authorization", config.AuthPrefix+jwt).GoWithHTTPHandler(t, r).Recorder
			assert.Equal(t, tt.status, rr.Code)

			// Not check answer if jwt not existed.
			if tt.hasJWT {
				var resultSite oapi.Site
				err = json.NewDecoder(rr.Body).Decode(&resultSite)

				require.NoError(t, err, "error unmarshaling response")
				t.Log("Result: ", resultSite)
			}

		})
	}
}
