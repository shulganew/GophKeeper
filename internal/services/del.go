package services

import (
	"net/http"

	"go.uber.org/zap"
)

// Delelete site data.
func (k *Keeper) DelAny(w http.ResponseWriter, r *http.Request, secretID string) {
	err := k.DeleteSecret(r.Context(), secretID)
	if err != nil {
		zap.S().Errorln("Error adding site to DB: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// set status code 200
	w.WriteHeader(http.StatusOK)
}
