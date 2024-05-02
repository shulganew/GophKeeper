package services

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"encoding/json"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/shulganew/GophKeeper/internal/api/oapi"
	"github.com/shulganew/GophKeeper/internal/entities"
	"go.uber.org/zap"
)

const PreambleLeth = 8

// Files add with two steps:
// 1. Uplod file and return created file id in minio storage.
// 2. Create file metadata as sectet in db with users description (definition field and file_id)
func (k *Keeper) AddGfile(w http.ResponseWriter, r *http.Request) {
	// Check registration.
	userID, isRegistered := CheckUserAuth(r.Context())
	if !isRegistered {
		http.Error(w, "JWT not found. Not authorized.", http.StatusUnauthorized)
		return
	}

	// Get file reader from body.
	fr := r.Body
	// Get preablule with lenth of metadata
	preamble := make([]byte, PreambleLeth)
	_, err := fr.Read(preamble)
	if err != nil {
		http.Error(w, "Can't Read preambule.", http.StatusInternalServerError)
	}

	// Read as uint64  value for metadatadata length.
	data := binary.LittleEndian.Uint64(preamble)
	meta := make([]byte, data)
	_, err = fr.Read(meta)
	if err != nil {
		http.Error(w, "Can't Read preambule.", http.StatusInternalServerError)
	}

	// Encode nfile as metadata to binary
	var nfile oapi.NewGfile
	err = gob.NewDecoder(bytes.NewReader(meta)).Decode(&nfile)
	if err != nil {
		http.Error(w, "Can't Read metadata.", http.StatusInternalServerError)
	}

	// Generate file id for minio.
	fileID, err := uuid.NewV7()
	if err != nil {
		http.Error(w, "Can't generate uuid.", http.StatusInternalServerError)
	}

	err = k.fstor.UploadFile(r.Context(), fileID.String(), fr)
	if err != nil {
		zap.S().Errorln("Can't upload: ", err)
		http.Error(w, "Can't upload file fo s3.", http.StatusInternalServerError)
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
	secretID, err := k.AddSecret(r.Context(), userID, entities.FILE, db.Bytes())
	if err != nil {
		zap.S().Errorln("Error adding gfile to DB: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return created gfile to client in responce (client add it to client's mem storage)
	gfile := oapi.Gfile{GfileID: secretID.String(), StorageID: fileID.String(), Definition: nfile.Definition, Fname: nfile.Fname}
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

	// Check registration.
	userID, isRegistered := CheckUserAuth(r.Context())
	if !isRegistered {
		http.Error(w, "JWT not found. Not authorized.", http.StatusUnauthorized)
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
	var gfiles []oapi.Gfile
	for _, secret := range secretDecoded {
		var gfile oapi.Gfile
		err = gob.NewDecoder(bytes.NewReader(secret.Data)).Decode(&gfile)
		if err != nil {
			zap.S().Errorln("Error decode gfile to data: ", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		gfile.GfileID = secret.UUID.String()
		gfiles = append(gfiles, gfile)
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
