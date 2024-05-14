package services

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/shulganew/GophKeeper/internal/api/jwt"
	"github.com/shulganew/GophKeeper/internal/api/oapi"
	"github.com/shulganew/GophKeeper/internal/entities"
	"go.uber.org/zap"
)

// Add new site credentials: site, login and password.
func (k *Keeper) AddSite(w http.ResponseWriter, r *http.Request) {
	// Get userID from jwt.
	userID, err := jwt.GetUserID(k.ua, r)
	if err != nil {
		zap.S().Errorln("Error getting userID: ", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	// Read all data from body for unmarshal and saving to sectert srorage.
	body, err := io.ReadAll(r.Body)
	if err != nil {
		zap.S().Errorln("Error reading body: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check json is correct.
	var newSite oapi.NewSite
	err = json.Unmarshal(body, &newSite)
	if err != nil {
		zap.S().Errorln("Can't Read json: ", err)
		http.Error(w, "Can't Read json.", http.StatusInternalServerError)
	}

	secretID, err := k.AddSecret(r.Context(), userID, entities.SITE, body)
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
	// Get userID from jwt.
	userID, err := jwt.GetUserID(k.ua, r)
	if err != nil {
		zap.S().Errorln("Error getting userID: ", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
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
		err = json.NewDecoder(bytes.NewReader(secret.Data)).Decode(&site)
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
	// Get userID from jwt.
	userID, err := jwt.GetUserID(k.ua, r)
	if err != nil {
		zap.S().Errorln("Error getting userID: ", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	// Read all data from body for unmarshal and saving to sectert srorage.
	body, err := io.ReadAll(r.Body)
	if err != nil {
		zap.S().Errorln("Error reading body: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check json is correct.
	var site oapi.Site
	err = json.Unmarshal(body, &site)
	if err != nil {
		zap.S().Errorln("Can't Read json site data: ", err)
		http.Error(w, "Can't Read json site data.", http.StatusInternalServerError)
	}

	err = k.UpdateSecret(r.Context(), userID, entities.SITE, body, site.SiteID)
	if err != nil {
		zap.S().Errorln("Error adding site to DB: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// set status code 200
	w.WriteHeader(http.StatusOK)
}
