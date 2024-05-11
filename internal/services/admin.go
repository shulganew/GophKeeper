package services

import (
	"encoding/json"
	"net/http"

	"github.com/shulganew/GophKeeper/internal/api/oapi"
	"github.com/shulganew/GophKeeper/internal/entities"
	"go.uber.org/zap"
)

// Generate new ephemeral key by admin request.
func (k *Keeper) EKeyNew(w http.ResponseWriter, r *http.Request) {
	eKey, ts, err := CreateEphemeralKey()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		zap.S().Errorln("Error create  eKeys: ", err)
	}
	// Generate new key.
	ekeyMem := entities.EKeyMem{EKey: eKey, TS: ts}
	// Add to mem keyring.
	k.eKeys = append(k.eKeys, ekeyMem)

	// Save to storage.
	k.SaveKeyRing(r.Context(), ekeyMem)

	// Set status code 201.
	w.WriteHeader(http.StatusCreated)
}

// Chane old master on new master key.
func (k *Keeper) NewMaster(w http.ResponseWriter, r *http.Request) {
	// Decode keys credentials from JSON.
	var key oapi.Key
	if err := json.NewDecoder(r.Body).Decode(&key); err != nil {
		zap.S().Errorln("Invalid format for Master keys", err)
		sendKeeperError(w, http.StatusBadRequest, "Invalid format for Master keys")
		return
	}
	zap.S().Debugf("Cerate new master key, %+v", key)
	if key.Old != k.conf.MasterKey {
		zap.S().Errorln("Old master mismatch")
		sendKeeperError(w, http.StatusBadRequest, "Old master mismatch")
		return
	}

	// Set new master key.
	k.conf.MasterKey = key.New

	// Drop eKeys encoded by old master.
	k.DropKeyRing(r.Context())

	// Save eKey from memory to database coded with new Master keys.
	k.SaveKeysRing(r.Context())
	// Set status code 201.
	w.WriteHeader(http.StatusCreated)
}
