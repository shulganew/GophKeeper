package services

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"net/http"

	"github.com/shulganew/GophKeeper/internal/api/oapi"
	"github.com/shulganew/GophKeeper/internal/entities"
	"go.uber.org/zap"
)

// Add new Gtext.
func (k *Keeper) AddGtext(w http.ResponseWriter, r *http.Request) {
	// Check registration.
	userID, isRegistered := CheckUserAuth(r.Context())
	if !isRegistered {
		http.Error(w, "JWT not found. Not authorized.", http.StatusUnauthorized)
		return
	}
	// Decode Gtext credentials from JSON.
	var newGtext oapi.NewGtext
	if err := json.NewDecoder(r.Body).Decode(&newGtext); err != nil {
		sendKeeperError(w, http.StatusBadRequest, "Invalid format for NewGtexts")
		return
	}
	// Write data to storage.
	var db bytes.Buffer
	err := gob.NewEncoder(&db).Encode(&newGtext)
	if err != nil {
		zap.S().Errorln("Error coding Gtext to data: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	secretID, err := k.AddSecret(r.Context(), userID, entities.TEXT, db.Bytes())
	if err != nil {
		zap.S().Errorln("Error adding Gtext to DB: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return created gtext to client in responce (client add it to client's mem storage)
	gtext := oapi.Gtext{GtextID: secretID.String(), Definition: newGtext.Definition, Note: newGtext.Note}
	w.Header().Add("Content-Type", "application/json")

	// set status code 201
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(gtext)
	if err != nil {
		zap.S().Errorln("Can't write to response in AddGtext handler", err)
	}
	zap.S().Debugln("Gtext credentials added. ", gtext.GtextID, " ", gtext.Definition)
}

// List all created Gtexts.
func (k *Keeper) ListGtexts(w http.ResponseWriter, r *http.Request) {
	// Check registration.
	userID, isRegistered := CheckUserAuth(r.Context())
	if !isRegistered {
		http.Error(w, "JWT not found. Not authorized.", http.StatusUnauthorized)
		return
	}

	// Load all user's Gtexts from database.
	secretDecoded, err := k.GetSecrets(r.Context(), userID, entities.TEXT)
	if err != nil {
		zap.S().Errorln("Error getting Gtext credentials: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Load decoded data and decode binary data to oapi.Gtext.
	var Gtexts []oapi.Gtext
	for _, secret := range secretDecoded {
		var newGtext oapi.NewGtext
		err = gob.NewDecoder(bytes.NewReader(secret.Data)).Decode(&newGtext)
		if err != nil {
			zap.S().Errorln("Error decode Gtext to data: ", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		Gtext := oapi.Gtext{GtextID: secret.UUID.String(), Definition: newGtext.Definition, Note: newGtext.Note}
		Gtexts = append(Gtexts, Gtext)
	}

	w.Header().Add("Content-Type", "application/json")
	if len(Gtexts) == 0 {
		zap.S().Infoln("No content.")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	// Set status code 200.
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(Gtexts)
	if err != nil {
		zap.S().Errorln("Can't write to response in ListGtexts handler", err)
	}
}
