package services

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"
	"github.com/oapi-codegen/testutil"
	"github.com/shulganew/GophKeeper/internal/api/oapi"
	"github.com/shulganew/GophKeeper/internal/api/router"
	"github.com/shulganew/GophKeeper/internal/app"
	"github.com/shulganew/GophKeeper/internal/app/config"
	"github.com/shulganew/GophKeeper/internal/entities"
	"github.com/shulganew/GophKeeper/internal/services/mocks"
	"go.uber.org/zap"

	"github.com/shulganew/GophKeeper/internal/api/jwt"
	"github.com/stretchr/testify/require"
)

const JWTPemKey = `
-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIP/dvfMGvKrM79LZuO9yfc3/HQGAvFoVzxYu2F1xkGKEoAoGCCqGSM49
AwEHoUQDQgAEV/0PntMTRVNu/ZZ8/mUdWZVCOevNaXlqUHSKR+YaC7X24Slj8HH1
cYJis1ufjejX19xk8XbFT8M1zyh4h0jwrw==
-----END EC PRIVATE KEY-----
`

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
			requestAdd: "/user/site",
			hasJWT:     true,
			status:     http.StatusCreated,
		},
		{
			name:       "Check user jwt",
			method:     http.MethodPost,
			requestAdd: "/user/site",
			hasJWT:     false,
			status:     http.StatusUnauthorized,
		},
	}

	app.InitLog()
	ctx := context.Background()

	conf := &config.Config{}

	// Init application
	conf.Address = "localhost:8088"
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
				AddSecretStor(gomock.Any(), gomock.Any(), gomock.Any()).
				AnyTimes().
				Return(&secret_id, nil)

			jsonSite, err := json.Marshal(nsite)
			require.NoError(t, err)

			//body := strings.NewReader(jsonSite)
			require.NoError(t, err)

			// Add jwt to header.
			var allowAll []byte
			if tt.hasJWT {
				allowAll, err = auth.CreateJWSWithClaims(userID.String(), []string{})
				require.NoError(t, err)

			}
			t.Log(string(allowAll))
			// Create request.
			rr := testutil.NewRequest().Post(tt.requestAdd).WithContentType("application/json").WithBody(jsonSite).WithHeader("Authorization", config.AuthPrefix+string(allowAll)).GoWithHTTPHandler(t, rt).Recorder
			require.Equal(t, tt.status, rr.Code)

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
