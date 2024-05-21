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

func TestText(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		requestPath    string
		requestPathDel string
		statusAdd      int
		statusList     int
		statusPut      int
		statusDel      int
		ntexts         []oapi.NewGtext
	}{
		{
			name:           "Check text add and text list methods",
			requestPath:    "/user/text",
			requestPathDel: "/user/",
			statusAdd:      http.StatusCreated,
			statusList:     http.StatusOK,
			statusPut:      http.StatusOK,
			statusDel:      http.StatusOK,
			ntexts:         []oapi.NewGtext{{Definition: "mytext1", Note: "Long story"}, {Definition: "mytext2", Note: "Long tail"}},
		},
	}

	ctx := context.Background()

	conf := &config.Config{}
	// Init application
	conf.Address = "localhost:8088"
	conf.MasterKey = "MasterPw"
	conf.DSN = ""
	conf.PathJWT = "cert/jwtpkey.pem"
	conf.ZapLevel = "debug"
	conf.RunLocal = true
	app.InitLog(*conf)
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

			// Test mem storage.
			storage := make(map[string]entities.SecretEncoded)

			_ = repo.EXPECT().
				AddSecretStor(gomock.Any(), gomock.Any(), gomock.Any()).
				DoAndReturn(func(_ any, e entities.NewSecretEncoded, types entities.SecretType) (*uuid.UUID, error) {
					secretID, err := uuid.NewV7()
					require.NoError(t, err)
					storage[secretID.String()] = entities.SecretEncoded{SecretID: secretID, NewSecret: e.NewSecret, DataCr: e.DataCr}
					return &secretID, nil
				}).
				AnyTimes()

			allowAll, err := auth.CreateJWSWithClaims(userID.String(), []string{})
			require.NoError(t, err)

			for _, text := range tt.ntexts {
				jsontext, err := json.Marshal(text)
				require.NoError(t, err)

				// Create request.
				rr := testutil.NewRequest().Post(tt.requestPath).WithContentType("application/json").WithBody(jsontext).WithHeader("Authorization", config.AuthPrefix+string(allowAll)).GoWithHTTPHandler(t, rt).Recorder
				require.Equal(t, tt.statusAdd, rr.Code)

				var resulttext oapi.Gtext
				err = json.NewDecoder(rr.Body).Decode(&resulttext)
				require.NoError(t, err, "error unmarshaling response")
				t.Log("Result: ", resulttext)
			}

			// List all texts data.
			t.Log("List data: ", len(storage))
			_ = repo.EXPECT().
				GetSecretsStor(gomock.Any(), gomock.Any(), gomock.Any()).
				DoAndReturn(func(_ any, _ string, _ entities.SecretType) (texts []*entities.SecretEncoded, err error) {
					s := make([]*entities.SecretEncoded, 0, len(storage))
					for _, value := range storage {
						s = append(s, &value)
					}
					return s, nil
				}).
				AnyTimes()

			// List all.
			rr := testutil.NewRequest().Get(tt.requestPath).WithHeader("Authorization", config.AuthPrefix+string(allowAll)).GoWithHTTPHandler(t, rt).Recorder
			require.Equal(t, tt.statusList, rr.Code)

			var secrets map[string]oapi.Gtext
			err = json.NewDecoder(rr.Body).Decode(&secrets)

			require.NoError(t, err, "error unmarshaling response")
			t.Log("Result: ", secrets)

			// Update checking.
			_ = repo.EXPECT().
				UpdateSecretStor(gomock.Any(), gomock.Any(), gomock.Any()).
				DoAndReturn(func(_ any, e entities.NewSecretEncoded, secretID string) (err error) {
					u, err := uuid.FromString(secretID)
					require.NoError(t, err)
					storage[secretID] = entities.SecretEncoded{SecretID: u, NewSecret: e.NewSecret, DataCr: e.DataCr}
					return nil
				}).
				AnyTimes()
			// Get secret id for test.
			var textIDs []string
			for k := range secrets {
				textIDs = append(textIDs, k)
			}

			// Cange first element.
			updated := oapi.Gtext{GtextID: textIDs[0], Definition: "Correct mytext1", Note: "New text"}
			jsontext, err := json.Marshal(updated)
			require.NoError(t, err)

			rr = testutil.NewRequest().Put(tt.requestPath).WithContentType("application/json").WithBody(jsontext).WithHeader("Authorization", config.AuthPrefix+string(allowAll)).GoWithHTTPHandler(t, rt).Recorder
			require.Equal(t, tt.statusPut, rr.Code)

			_ = repo.EXPECT().
				DeleteSecretStor(gomock.Any(), gomock.Any()).
				DoAndReturn(func(_ any, secretID string) (err error) {
					delete(storage, secretID)
					t.Log("Delete from storage, ", secretID, " current len: ", len(storage))
					return nil
				}).
				AnyTimes()

			// Delete second element.
			rr = testutil.NewRequest().Delete(tt.requestPathDel+textIDs[1]).WithHeader("Authorization", config.AuthPrefix+string(allowAll)).GoWithHTTPHandler(t, rt).Recorder
			require.Equal(t, tt.statusDel, rr.Code)

			// List elements and check existense of first siteIDs[0]
			rr = testutil.NewRequest().Get(tt.requestPath).WithHeader("Authorization", config.AuthPrefix+string(allowAll)).GoWithHTTPHandler(t, rt).Recorder
			require.Equal(t, tt.statusList, rr.Code)

			secrets = make(map[string]oapi.Gtext)
			err = json.NewDecoder(rr.Body).Decode(&secrets)
			require.NoError(t, err)
			// Updated element exist and equal.
			require.Equal(t, updated, secrets[textIDs[0]])
			// Deleted second element.
			require.Equal(t, 1, len(secrets))
		})
	}
}
