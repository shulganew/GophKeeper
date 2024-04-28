package services

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"net/http"
	"time"

	"github.com/shulganew/GophKeeper/internal/entities"
	"github.com/shulganew/GophKeeper/internal/rest/oapi"
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

	data := db.Bytes()
	// Get data key.
	dKey, _, err := CreateDataKey()
	if err != nil {
		zap.S().Errorln("Error create data key: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Encode date before store.
	datac, err := EncodeData(dKey, data)
	if err != nil {
		zap.S().Errorln("Error encode data: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Get Ephemeral current key.
	eKey := k.GetActualEKey()
	dKeyc, err := EncodeKey(dKey, eKey.EKey)
	if err != nil {
		zap.S().Errorln("Error getting ephemeral key: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dbSite := entities.Secret{UserID: userID, Stype: entities.SITE, Data: datac, EKeyVer: eKey.TS, DKey: dKeyc, Uploaded: time.Now()}
	siteID, err := k.stor.AddSite(r.Context(), dbSite)
	if err != nil {
		zap.S().Errorln("Error adding site credentials: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Return created site to client in responce (client add it to client's mem storage)
	site := oapi.Site{SiteID: siteID.String(), Site: newSite.Site, Slogin: newSite.Slogin, Spw: newSite.Spw}

	w.Header().Add("Content-Type", "application/json")

	// set status code 201
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(site)
	if err != nil {
		zap.S().Errorln("Can't write to response in AddSite handler", err)
	}
	zap.S().Infoln("Site credentials added. ", userID)
}

// List all users sites with credentials.
func (k *Keeper) ListSite(w http.ResponseWriter, r *http.Request) {
	// Check registration.
	userID, isRegistered := CheckUserAuth(r.Context())
	if !isRegistered {
		http.Error(w, "JWT not found. Not authorized.", http.StatusUnauthorized)
		return
	}

	// Load all user's sites credentials from database.
	dbSites, err := k.stor.GetSites(r.Context(), userID, entities.SITE)
	if err != nil {
		zap.S().Errorln("Error getting site credentials: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Decode secret data from storage.
	var sites []oapi.Site
	for _, dbSite := range dbSites {

		// Get ephemeral key (version from ts in db) for decode data key.
		eKey := k.GetEKey(dbSite.EKeyVer)
		// Decode dKeyc
		dKey, err := DecodeKey(dbSite.DKey, eKey.EKey)
		if err != nil {
			zap.S().Errorln("Error decode data key: ", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Decode data.
		data, err := DecodeData(dKey, dbSite.Data)
		if err != nil {
			zap.S().Errorln("Error decode stored data: ", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		var newSite oapi.NewSite
		err = gob.NewDecoder(bytes.NewReader(data)).Decode(&newSite)
		if err != nil {
			zap.S().Errorln("Error decode site to data: ", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		site := oapi.Site{Definition: newSite.Definition, SiteID: dbSite.UUID.String(), Site: newSite.Site, Slogin: newSite.Slogin, Spw: newSite.Spw}
		sites = append(sites, site)
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
	zap.S().Infoln("Site credentials added. ", userID)
}
