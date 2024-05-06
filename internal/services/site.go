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

// Add new site credentials: site, login and password.
func (k *Keeper) AddSite(w http.ResponseWriter, r *http.Request) {
	// Check registration.
	userID, isRegistered := CheckUserAuth(r.Context())
	if !isRegistered {
		http.Error(w, "JWT not found. Not authorized.", http.StatusUnauthorized)
		return
	}
	// Decode Site credentials from JSON.
	var newSite oapi.NewSite
	if err := json.NewDecoder(r.Body).Decode(&newSite); err != nil {
		sendKeeperError(w, http.StatusBadRequest, "Invalid format for NewSite")
		return
	}
	// Write data to storage.
	var db bytes.Buffer
	err := gob.NewEncoder(&db).Encode(&newSite)
	if err != nil {
		zap.S().Errorln("Error coding site to data: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	secretID, _, err := k.AddSecret(r.Context(), userID, entities.SITE, db.Bytes())
	if err != nil {
		zap.S().Errorln("Error adding site to DB: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Return created site to client in responce (client add it to client's mem storage)
	site := oapi.Site{SiteID: secretID.String(), Definition: newSite.Definition, Site: newSite.Site, Slogin: newSite.Slogin, Spw: newSite.Spw}
	w.Header().Add("Content-Type", "application/json")

	// set status code 201
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(site)
	if err != nil {
		zap.S().Errorln("Can't write to response in AddSite handler", err)
	}
	zap.S().Debugln("Site credentials added. ", site.SiteID, " ", site.Site)
}

// List all users sites with credentials.
func (k *Keeper) ListSites(w http.ResponseWriter, r *http.Request) {
	// Check registration.
	userID, isRegistered := CheckUserAuth(r.Context())
	if !isRegistered {
		http.Error(w, "JWT not found. Not authorized.", http.StatusUnauthorized)
		return
	}

	// Load all user's sites credentials from database.
	secretDecoded, err := k.GetSecrets(r.Context(), userID, entities.SITE)
	if err != nil {
		zap.S().Errorln("Error getting site credentials: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Load decoded data and decode binary data to oapi.Site.
	sites := make(map[string]oapi.Site, len(secretDecoded))
	for _, secret := range secretDecoded {
		var site oapi.Site
		err = gob.NewDecoder(bytes.NewReader(secret.Data)).Decode(&site)
		if err != nil {
			zap.S().Errorln("Error decode site to data: ", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		site.SiteID = secret.SecretID.String()
		sites[site.SiteID] = site
	}

	w.Header().Add("Content-Type", "application/json")
	if len(sites) == 0 {
		zap.S().Infoln("No content.")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	// Set status code 200.
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(sites)
	if err != nil {
		zap.S().Errorln("Can't write to response in ListSite handler", err)
	}
}

// Site data update.
func (k *Keeper) UpdateSite(w http.ResponseWriter, r *http.Request) {
	// Check registration.
	userID, isRegistered := CheckUserAuth(r.Context())
	if !isRegistered {
		http.Error(w, "JWT not found. Not authorized.", http.StatusUnauthorized)
		return
	}
	// Decode site credentials from JSON.
	var site oapi.Site
	if err := json.NewDecoder(r.Body).Decode(&site); err != nil {
		sendKeeperError(w, http.StatusBadRequest, "Invalid format for NewSite")
		return
	}
	// Write data to storage.
	var db bytes.Buffer
	err := gob.NewEncoder(&db).Encode(&site)
	if err != nil {
		zap.S().Errorln("Error coding site to data: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = k.UpdateSecret(r.Context(), userID, entities.SITE, db.Bytes(), site.SiteID)
	if err != nil {
		zap.S().Errorln("Error adding site to DB: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// set status code 200
	w.WriteHeader(http.StatusOK)
}
