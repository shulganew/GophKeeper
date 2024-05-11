package services

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/shulganew/GophKeeper/internal/api/jwt"
	"github.com/shulganew/GophKeeper/internal/api/oapi"
	"github.com/shulganew/GophKeeper/internal/entities"
	"go.uber.org/zap"
)

const PreambleLeth = 8

// Files add with two steps:
// 1. Uplod file and return created file id in minio storage.
// 2. Create file metadata as sectet in db with users description (definition field and file_id)
func (k *Keeper) AddGfile(w http.ResponseWriter, r *http.Request) {
	// Get userID from jwt.
	userID, err := jwt.GetUserID(k.ua, r)
	if err != nil {
		zap.S().Errorln("Error getting userID: ", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Get file reader from body.
	fr := r.Body
	// Get preablule with lenth of metadata
	preamble := make([]byte, PreambleLeth)
	_, err = fr.Read(preamble)
	if err != nil {
		http.Error(w, "Can't Read preambule.", http.StatusInternalServerError)
	}

	// Read as uint64  value for metadatadata length.
	data := binary.LittleEndian.Uint64(preamble)
	meta := make([]byte, data)
	_, err = fr.Read(meta)
	if err != nil {
		zap.S().Errorln("Can't Read preambule: ", err)
		http.Error(w, "Can't Read preambule.", http.StatusInternalServerError)
	}

	// Encode nfile as metadata to binary
	var nfile oapi.NewGfile
	err = gob.NewDecoder(bytes.NewReader(meta)).Decode(&nfile)
	if err != nil {
		zap.S().Errorln("Can't Read metadata: ", err)
		http.Error(w, "Can't Read metadata.", http.StatusInternalServerError)
	}

	// Encode data.
	var db bytes.Buffer
	err = gob.NewEncoder(&db).Encode(&nfile)
	if err != nil {
		zap.S().Errorln("Error coding Gfile to data: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Save file metadata newGfile to DB.
	gfileID, dKey, err := k.AddSecret(r.Context(), userID, entities.FILE, db.Bytes())
	if err != nil {
		zap.S().Errorln("Error adding gfile to DB: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Encode data before updloading.
	zap.S().Debugln("data key add: ", hex.EncodeToString(dKey))

	dataF, err := io.ReadAll(fr)
	if err != nil {
		zap.S().Errorln("Error read file data: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fr.Close()

	dataFc, err := EncodeData(dKey, dataF)
	if err != nil {
		zap.S().Errorln("Error decode data: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Upload file to minio
	err = k.fstor.UploadFile(r.Context(), k.conf.Backetmi, gfileID.String(), bytes.NewBuffer(dataFc))
	if err != nil {
		zap.S().Errorln("Can't upload: ", err)
		http.Error(w, "Can't upload file fo s3.", http.StatusInternalServerError)
	}

	// Return created gfile to client in responce (client add it to client's mem storage)
	gfile := oapi.Gfile{GfileID: gfileID.String(), Definition: nfile.Definition, Fname: nfile.Fname}

	w.Header().Add("Content-Type", "application/json")

	// set status code 201
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(gfile)
	if err != nil {
		zap.S().Errorln("Can't write to response in Addgfile handler", err)
	}
	zap.S().Debugln("gfile credentials added. ", gfile.GfileID, " ", gfile.Definition)
}

// Return add gfiles metadata from DB.
func (k *Keeper) ListGfiles(w http.ResponseWriter, r *http.Request) {
	// Get userID from jwt.
	userID, err := jwt.GetUserID(k.ua, r)
	if err != nil {
		zap.S().Errorln("Error getting userID: ", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Load all user's gfiles from database.
	secretDecoded, err := k.GetSecrets(r.Context(), userID, entities.FILE)
	if err != nil {
		zap.S().Errorln("Error getting Gfile credentials: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Load decoded data and decode binary data to oapi.gfile.
	gfiles := make(map[string]oapi.Gfile, len(secretDecoded))
	for _, secret := range secretDecoded {
		var gfile oapi.Gfile
		err = gob.NewDecoder(bytes.NewReader(secret.Data)).Decode(&gfile)
		if err != nil {
			zap.S().Errorln("Error decode gfile to data: ", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		gfile.GfileID = secret.SecretID.String()
		gfiles[gfile.GfileID] = gfile
	}

	w.Header().Add("Content-Type", "application/json")
	if len(gfiles) == 0 {
		zap.S().Infoln("No content.")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	// Set status code 200.
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(gfiles)
	if err != nil {
		zap.S().Errorln("Can't write to response in Listgfiles handler", err)
	}

}

// Return gfile from storage. fileID == secretID (+) entities.FILE
func (k *Keeper) GetGfile(w http.ResponseWriter, r *http.Request, fileID string) {
	// Get userID from jwt.
	userID, err := jwt.GetUserID(k.ua, r)
	if err != nil {
		zap.S().Errorln("Error getting userID: ", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	zap.S().Debugln("File id:", fileID)
	_, err = uuid.FromString(fileID)
	if err != nil {
		zap.S().Errorln("file ID not correct: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Get Gfile from DB.
	// Load all user's gfiles from database.
	secretDecoded, err := k.GetSecret(r.Context(), userID, entities.FILE, fileID)
	if err != nil {
		zap.S().Errorln("Error getting Gfile credentials: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if secretDecoded == nil {
		zap.S().Infoln("File not found")
		http.Error(w, "file not found", http.StatusInternalServerError)
		return
	}

	// Load decoded data and decode binary data to oapi.gfile.
	var gfile oapi.Gfile
	err = gob.NewDecoder(bytes.NewReader(secretDecoded.Data)).Decode(&gfile)
	if err != nil {
		zap.S().Errorln("Error decode gfile to data: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	gfile.GfileID = secretDecoded.SecretID.String()
	fr, err := k.fstor.DownloadFile(r.Context(), k.conf.Backetmi, gfile.GfileID)
	if err != nil {
		zap.S().Errorln("Can't Download: ", err)
		http.Error(w, "Can't Download.", http.StatusInternalServerError)
	}
	defer func() {
		err := fr.Close()
		if err != nil {
			zap.S().Errorln("Can't close minio: ", err)
		}
	}()

	// set status code 200
	w.WriteHeader(http.StatusOK)

	// Decod data.
	dataFc, err := io.ReadAll(fr)
	if err != nil {
		zap.S().Errorln("Error read file data: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dataF, err := DecodeData(secretDecoded.DKeyCr, dataFc)
	if err != nil {
		zap.S().Errorln("Error decode data: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Copy and Decode file reader.
	if _, err := io.CopyN(w, bytes.NewBuffer(dataF), int64(len(dataF))); err != nil {
		zap.S().Errorln("Can't copy to resp: ", err)
		http.Error(w, "Can't copy to resp.", http.StatusInternalServerError)
	}
}

// Return gfile from storage. fileID == secretID (+) entities.FILE
func (k *Keeper) DelGfile(w http.ResponseWriter, r *http.Request, fileID string) {

	zap.S().Debugln("File id:", fileID)
	_, err := uuid.FromString(fileID)
	if err != nil {
		zap.S().Errorln("file ID not correct: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Del Gfile from DB.
	err = k.DeleteSecret(r.Context(), fileID)
	if err != nil {
		zap.S().Errorln("Error deleting file in DB: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = k.fstor.DeleteFile(r.Context(), k.conf.Backetmi, fileID)
	if err != nil {
		zap.S().Errorln("Can't Delete in File Storage: ", err)
		http.Error(w, "Can't Delete in File Storage.", http.StatusInternalServerError)
	}

	// set status code 200
	w.WriteHeader(http.StatusOK)
}
