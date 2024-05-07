package services

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"net/http"

	"github.com/shulganew/GophKeeper/internal/api/jwt"
	"github.com/shulganew/GophKeeper/internal/api/oapi"
	"github.com/shulganew/GophKeeper/internal/entities"
	"go.uber.org/zap"
)

// Add new Gtext.
func (k *Keeper) AddGtext(w http.ResponseWriter, r *http.Request) {
	// Get userID from jwt.
	userID, err := jwt.GetUserID(k.ua, r)
	if err != nil {
		zap.S().Errorln("Error getting userID: ", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
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
	err = gob.NewEncoder(&db).Encode(&newGtext)
	if err != nil {
		zap.S().Errorln("Error coding Gtext to data: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	secretID, _, err := k.AddSecret(r.Context(), userID, entities.TEXT, db.Bytes())
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
	// Get userID from jwt.
	userID, err := jwt.GetUserID(k.ua, r)
	if err != nil {
		zap.S().Errorln("Error getting userID: ", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
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
	gtexts := make(map[string]oapi.Gtext, len(secretDecoded))
	for _, secret := range secretDecoded {
		var gtext oapi.Gtext
		err = gob.NewDecoder(bytes.NewReader(secret.Data)).Decode(&gtext)
		if err != nil {
			zap.S().Errorln("Error decode Gtext to data: ", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		gtext.GtextID = secret.SecretID.String()
		gtexts[gtext.GtextID] = gtext
	}

	w.Header().Add("Content-Type", "application/json")
	if len(gtexts) == 0 {
		zap.S().Infoln("No content.")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	// Set status code 200.
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(gtexts)
	if err != nil {
		zap.S().Errorln("Can't write to response in ListGtexts handler", err)
	}
}

func (k *Keeper) UpdateGtext(w http.ResponseWriter, r *http.Request) {
	// Get userID from jwt.
	userID, err := jwt.GetUserID(k.ua, r)
	if err != nil {
		zap.S().Errorln("Error getting userID: ", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Decode gtext credentials from JSON.
	var gtext oapi.Gtext
	if err := json.NewDecoder(r.Body).Decode(&gtext); err != nil {
		sendKeeperError(w, http.StatusBadRequest, "Invalid format for NewSite")
		return
	}
	// Write data to storage.
	var db bytes.Buffer
	err = gob.NewEncoder(&db).Encode(&gtext)
	if err != nil {
		zap.S().Errorln("Error coding site to data: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = k.UpdateSecret(r.Context(), userID, entities.SITE, db.Bytes(), gtext.GtextID)
	if err != nil {
		zap.S().Errorln("Error adding site to DB: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// set status code 200
	w.WriteHeader(http.StatusOK)
}
