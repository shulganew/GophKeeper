package services

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/shulganew/GophKeeper/internal/api/oapi"
	"go.uber.org/zap"
)

const PreambleLeth = 8

// Files add with two steps:
// 1. Uplod file and return created file id in minio storage.
// 2. Create file metadata as sectet in db with users description (definition field and file_id)
//

func (k *Keeper) AddGfile(w http.ResponseWriter, r *http.Request) {

	// Generate file id for minio.
	fileID, err := uuid.NewV7()
	if err != nil {
		http.Error(w, "Can't generate uuid.", http.StatusInternalServerError)
	}
	// Get file reader from body.
	fr := r.Body

	preamble := make([]byte, PreambleLeth)
	_, err = fr.Read(preamble)
	if err != nil {
		http.Error(w, "Can't Read preambule.", http.StatusInternalServerError)
	}
	zap.S().Infoln("Peambule: ", preamble)
	data := binary.LittleEndian.Uint64(preamble)
	zap.S().Infoln("Metadata size: ", data)

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

	zap.S().Infoln("Metadata: ", nfile)

	err = k.fstor.UploadFile(r.Context(), fileID.String(), fr)
	if err != nil {
		zap.S().Errorln("Can't upload: ", err)
		http.Error(w, "Can't upload file fo s3.", http.StatusInternalServerError)
	}

}
