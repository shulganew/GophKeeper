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

func TestAdminPremissions(t *testing.T) {
	tests := []struct {
		name             string
		permission       string
		path             string
		requestPathAdmin string
		statusPut        int
	}{
		{
			name:             "Ephemeral key chande",
			permission:       "admin",
			requestPathAdmin: "/admin/key",
			statusPut:        http.StatusCreated,
		},
		{
			name:             "Ephemeral key chande: wrong permissions",
			permission:       "",
			requestPathAdmin: "/admin/key",
			statusPut:        http.StatusUnauthorized,
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

			allowAll, err := auth.CreateJWSWithClaims(userID.String(), []string{tt.permission})
			require.NoError(t, err)

			_ = repo.EXPECT().
				SaveEKeyc(gomock.Any(), gomock.Any()).
				AnyTimes().
				Return(nil)

			// Make new ephemeral key
			rr := testutil.NewRequest().Put(tt.requestPathAdmin).WithHeader("Authorization", config.AuthPrefix+string(allowAll)).GoWithHTTPHandler(t, rt).Recorder
			t.Log("Response: ", rr.Code)
			require.Equal(t, tt.statusPut, rr.Code)

		})
	}
}

func TestAdminEKey(t *testing.T) {
	tests := []struct {
		name             string
		path             string
		requestPath      string
		requestPathAdmin string
		statusAdd        int
		statusList       int
		statusPut        int
		ntext            oapi.NewGtext
	}{
		{
			name:             "Ephemeral key chande",
			requestPath:      "/user/text",
			requestPathAdmin: "/admin/key",
			statusAdd:        http.StatusCreated,
			statusList:       http.StatusOK,
			statusPut:        http.StatusCreated,
			ntext:            oapi.NewGtext{Definition: "mytext1", Note: "Long story"},
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

	// Test eKey mem storage.
	eKeyStorage := make([]entities.EKeyDB, 0)
	// Create new eKey.
	_ = repo.EXPECT().
		SaveEKeyc(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ any, eKeyc entities.EKeyDB) error {
			eKeyStorage = append(eKeyStorage, eKeyc)
			return nil
		}).
		AnyTimes()

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

			userPerm, err := auth.CreateJWSWithClaims(userID.String(), []string{})
			require.NoError(t, err)

			jsontext, err := json.Marshal(tt.ntext)
			require.NoError(t, err)

			// Create request with current ephemeral key.
			rr := testutil.NewRequest().Post(tt.requestPath).WithContentType("application/json").WithBody(jsontext).WithHeader("Authorization", config.AuthPrefix+string(userPerm)).GoWithHTTPHandler(t, rt).Recorder
			require.Equal(t, tt.statusAdd, rr.Code)

			var resulttext oapi.Gtext
			err = json.NewDecoder(rr.Body).Decode(&resulttext)
			require.NoError(t, err, "error unmarshaling response")
			t.Log("Result: ", resulttext)

			// JWT with admin rigth.
			adminPerm, err := auth.CreateJWSWithClaims(userID.String(), []string{"admin"})
			require.NoError(t, err)

			rr = testutil.NewRequest().Put(tt.requestPathAdmin).WithHeader("Authorization", config.AuthPrefix+string(adminPerm)).GoWithHTTPHandler(t, rt).Recorder
			t.Log("Response: ", rr.Code)
			require.Equal(t, tt.statusPut, rr.Code)

			// Create request with new ephemeral key.
			rr = testutil.NewRequest().Post(tt.requestPath).WithContentType("application/json").WithBody(jsontext).WithHeader("Authorization", config.AuthPrefix+string(userPerm)).GoWithHTTPHandler(t, rt).Recorder
			require.Equal(t, tt.statusAdd, rr.Code)

			resulttext = oapi.Gtext{}
			err = json.NewDecoder(rr.Body).Decode(&resulttext)
			require.NoError(t, err, "error unmarshaling response")
			t.Log("Result: ", resulttext)

			// Two key created.
			require.Equal(t, 2, len(eKeyStorage))

			// All keys existed
			keys := make([]entities.EKeyMem, 0)
			for _, v := range storage {
				key, err := keeper.GetEKey(v.EKeyVer)
				require.NoError(t, err)
				keys = append(keys, *key)
			}

			// Two key created.
			require.Equal(t, 2, len(keys))

			// Keys not equal.
			require.NotEqual(t, keys[0], keys[1])

		})
	}
}

func TestAdminMaster(t *testing.T) {
	tests := []struct {
		name             string
		path             string
		requestPath      string
		requestPathAdmin string
		oldMaster        string
		newMaster        string
		statusAdd        int
		statusList       int
		statusPut        int
		ntext            oapi.NewGtext
	}{
		{
			name:             "MAster key chande",
			requestPath:      "/user/text",
			requestPathAdmin: "/admin/master",
			oldMaster:        "OldMaster",
			newMaster:        "NewMaster",
			statusAdd:        http.StatusCreated,
			statusList:       http.StatusOK,
			statusPut:        http.StatusCreated,
			ntext:            oapi.NewGtext{Definition: "mytext1", Note: "Long story"},
		},
	}

	app.InitLog()
	ctx := context.Background()

	conf := &config.Config{}

	// Init application
	conf.Address = "localhost:8088"
	conf.MasterKey = tests[0].oldMaster
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

	// Test eKey mem storage.
	eKeyStorage := make([]entities.EKeyDB, 0)
	// Create new eKey.
	_ = repo.EXPECT().
		SaveEKeyc(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ any, eKeyc entities.EKeyDB) error {
			eKeyStorage = append(eKeyStorage, eKeyc)
			return nil
		}).
		AnyTimes()

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

			// Save old eKey coded from DB, chande Master and compare new coded eKey.
			oldEKeyCoded := eKeyStorage[len(eKeyStorage)-1]

			// JWT with admin rigth.
			adminPerm, err := auth.CreateJWSWithClaims(userID.String(), []string{"admin"})
			require.NoError(t, err)

			key := oapi.Key{Old: tt.oldMaster, New: tt.newMaster}
			jsontext, err := json.Marshal(key)
			require.NoError(t, err)

			// Change master
			_ = repo.EXPECT().
				DropKeys(gomock.Any()).
				DoAndReturn(func(_ any) error {
					// Clean storage
					eKeyStorage = make([]entities.EKeyDB, 0)
					return nil
				}).
				AnyTimes()

			rr := testutil.NewRequest().Post(tt.requestPathAdmin).WithContentType("application/json").WithBody(jsontext).WithHeader("Authorization", config.AuthPrefix+string(adminPerm)).GoWithHTTPHandler(t, rt).Recorder
			t.Log("Response: ", rr.Code)
			require.Equal(t, tt.statusPut, rr.Code)

			newEKeyCoded := eKeyStorage[len(eKeyStorage)-1]

			// Key was chaned.
			require.NotEqual(t, oldEKeyCoded, newEKeyCoded)

			oldKey, err := DecodeKey(oldEKeyCoded.EKeyc, []byte(tt.oldMaster))
			require.NoError(t, err)
			newKey, err := DecodeKey(newEKeyCoded.EKeyc, []byte(tt.newMaster))
			require.NoError(t, err)

			// Old and new are same.
			require.Equal(t, oldKey, newKey)

		})
	}
}
