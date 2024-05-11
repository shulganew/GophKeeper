package services

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/gob"
	"encoding/json"
	"io"
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

// Implementation from gophkeeper client.
// Files add with two steps:
// 1. Uplod file and return created file id in minio storage.
// 2. Create file metadata as sectet in db with users description (definition field and file_id)
type UploadReader struct {
	file      io.Reader
	preambule []byte
	metadata  []byte
	index     int64
	metaLen   int64
}

// Constructor for Upload files.
// byte structute: |8-byte preambule with meta length| N-bytes metadata newGfile | File bytes |
func NewUploadReader(file io.Reader, preambule []byte, metadata []byte) *UploadReader {
	r := new(UploadReader)
	r.file = file
	r.preambule = preambule
	r.metadata = metadata
	r.metaLen = int64(len(r.metadata))
	return r
}

// Read to b []byte preambule, then metadata, then original file.
func (r *UploadReader) Read(b []byte) (totlal int, err error) {
	// Add preambule bytes (PreambleLeth), witch contains lenth of metadata (newGfile)
	if r.index < PreambleLeth {
		n := copy(b, r.preambule[r.index:PreambleLeth])
		r.index += int64(n)
		totlal += n
	}

	// Add metadata bytes - newGfiles object.
	if r.index >= PreambleLeth && r.index < PreambleLeth+r.metaLen {
		n := copy(b[PreambleLeth:], r.metadata[r.index-PreambleLeth:r.metaLen])
		r.index += int64(n)
		totlal += n
	}
	// Add file bytes
	if r.index >= PreambleLeth+r.metaLen {
		bf := make([]byte, len(b)-totlal)
		fn, err := r.file.Read(bf)
		if err != nil {
			return totlal, err
		}
		n := copy(b[totlal:totlal+fn], bf)
		r.index += int64(n)
		totlal += n
		return totlal, nil

	}
	return
}

func TestGfile(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		requestPath    string
		requestPathDel string
		statusAdd      int
		statusList     int
		statusPut      int
		statusDel      int
		nGfiles        []oapi.NewGfile
		files          [][]byte
	}{
		{
			name:           "Check Gfile add and Gfile list methods",
			requestPath:    "/user/file",
			requestPathDel: "/user/",
			statusAdd:      http.StatusCreated,
			statusList:     http.StatusOK,
			statusPut:      http.StatusOK,
			statusDel:      http.StatusOK,
			nGfiles:        []oapi.NewGfile{{Definition: "myGfile1", Fname: "secret.png"}, {Definition: "myGfile1", Fname: "secret.img"}},
			files:          [][]byte{[]byte(("file data 1")), []byte(("file data 1"))},
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
	fstor := mocks.NewMockFileKeeper(ctrl)

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

	keeper := NewKeeper(ctx, repo, fstor, *conf, auth)

	// We now register our GophKeeper above as the handler for the interface
	oapi.HandlerFromMux(keeper, rt)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log("Test name: ", tt.name)
			// Generate userID.
			userID, err := uuid.NewV7()
			require.NoError(t, err)

			// Test mem storage. map[secretID]
			storage := make(map[string]entities.SecretEncoded)

			// Test s3 storage, map[fileID]filebytes
			fileStorage := make(map[string][]byte)
			_ = repo.EXPECT().
				AddSecretStor(gomock.Any(), gomock.Any(), gomock.Any()).
				DoAndReturn(func(_ any, e entities.NewSecretEncoded, types entities.SecretType) (*uuid.UUID, error) {
					fileID, err := uuid.NewV7()
					require.NoError(t, err)
					storage[fileID.String()] = entities.SecretEncoded{SecretID: fileID, NewSecret: e.NewSecret, DataCr: e.DataCr}
					return &fileID, nil
				}).
				AnyTimes()

			_ = fstor.EXPECT().
				UploadFile(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				DoAndReturn(func(_ any, _ string, fileID string, fr io.Reader) error {
					file, err := io.ReadAll(fr)
					require.NoError(t, err)
					fileStorage[fileID] = file
					return nil
				}).
				AnyTimes()

				// JWT Auth.
			allowAll, err := auth.CreateJWSWithClaims(userID.String(), []string{})
			require.NoError(t, err)

			for i, nfile := range tt.nGfiles {
				// Create file
				// Encode nfile as metadata to binary
				var md bytes.Buffer
				err = gob.NewEncoder(&md).Encode(&nfile)
				require.NoError(t, err)

				metadata := md.Bytes()
				// Write preambule - size of newGfile object (metadata).
				p := make([]byte, PreambleLeth)
				mLen := uint64(len(metadata))
				zap.S().Debugln("Metadata length: ", mLen)
				binary.LittleEndian.PutUint64(p, mLen)

				ur := NewUploadReader(bytes.NewReader(tt.files[i]), p, metadata)
				body, err := io.ReadAll(ur)
				require.NoError(t, err)
				// Create request.
				rr := testutil.NewRequest().Post(tt.requestPath).WithContentType("application/octet-stream").WithBody(body).WithHeader("Authorization", config.AuthPrefix+string(allowAll)).GoWithHTTPHandler(t, rt).Recorder
				require.Equal(t, tt.statusAdd, rr.Code)

				var resultGfile oapi.Gfile
				err = json.NewDecoder(rr.Body).Decode(&resultGfile)
				require.NoError(t, err, "error unmarshaling response")
				t.Log("Result: ", resultGfile)
			}

			// List all Gfiles data.
			_ = repo.EXPECT().
				GetSecretsStor(gomock.Any(), gomock.Any(), gomock.Any()).
				DoAndReturn(func(_ any, _ string, _ entities.SecretType) (gfiles []entities.SecretEncoded, err error) {
					s := make([]entities.SecretEncoded, 0, len(storage))
					for _, value := range storage {
						s = append(s, value)
					}
					return s, nil
				}).
				AnyTimes()

			// Get one.
			_ = repo.EXPECT().
				GetSecretStor(gomock.Any(), gomock.Any()).
				DoAndReturn(func(_ any, fileID string) (*entities.SecretEncoded, error) {
					gfile := storage[fileID]
					return &gfile, nil
				}).
				AnyTimes()
			// List all
			rr := testutil.NewRequest().Get(tt.requestPath).WithHeader("Authorization", config.AuthPrefix+string(allowAll)).GoWithHTTPHandler(t, rt).Recorder
			require.Equal(t, tt.statusList, rr.Code)

			var gfiles map[string]oapi.Gfile
			err = json.NewDecoder(rr.Body).Decode(&gfiles)

			require.NoError(t, err, "error unmarshaling response")
			t.Log("Result: ", gfiles)

			// Get secret id for test.
			var GfileIDs []string
			for k := range gfiles {
				GfileIDs = append(GfileIDs, k)
			}
			// Total listed 2 files.
			require.Equal(t, 2, len(gfiles))

			// Get file by id.
			_ = fstor.EXPECT().
				DownloadFile(gomock.Any(), gomock.Any(), gomock.Any()).
				DoAndReturn(func(_ any, _ string, fileID string) (fr io.Reader, err error) {
					t.Log("Download file ", fileID)
					f := fileStorage[fileID]
					r := bytes.NewReader(f)
					rc := io.NopCloser(r) // Wrapper for io.readcloser.
					return rc, nil
				}).
				AnyTimes()

			path := tt.requestPath + "/" + GfileIDs[0]
			rr = testutil.NewRequest().Get(path).WithHeader("Authorization", config.AuthPrefix+string(allowAll)).GoWithHTTPHandler(t, rt).Recorder
			require.Equal(t, tt.statusList, rr.Code)
			file, err := io.ReadAll(rr.Body)
			require.NoError(t, err)
			require.Equal(t, tt.files[0], file)
		})
	}
}
