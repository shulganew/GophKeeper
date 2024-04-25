package services

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"net/http"

	"github.com/shulganew/GophKeeper/internal/entities"
	"github.com/shulganew/GophKeeper/internal/rest/oapi"
	"go.uber.org/zap"
)

// Add new site credentials: site, login and password.
func (k *Keeper) AddSite(w http.ResponseWriter, r *http.Request) {
	// Check registration.
	userID, isRegistered := CheckUserAuth(r.Context())
	if isRegistered {
		http.Error(w, "JWT not found. Not authorized.", http.StatusUnauthorized)
		return
	}

	var site oapi.NewSite
	if err := json.NewDecoder(r.Body).Decode(&site); err != nil {
		// If can't decode 400
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var data bytes.Buffer
	encoder := gob.NewEncoder(&data)
	err := encoder.Encode(&site)
	if err != nil {
		zap.S().Errorln("Error cover site to data: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dbSite := entities.Secret{UserID: userID, Stype: entities.SITE, Data: data.Bytes()}
	err = k.stor.AddSite(r.Context(), dbSite)
	if err != nil {
		zap.S().Errorln("Error adding site credentials: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "text/plain")

	// set status code 201
	w.WriteHeader(http.StatusCreated)

	_, err = w.Write([]byte("Site credentials added."))
	if err != nil {
		zap.S().Errorln("Can't write to response in LoginUser handler", err)
	}
	zap.S().Infoln("Site credentials added. ", userID)
}
