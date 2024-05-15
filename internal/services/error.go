package services

import (
	"encoding/json"
	"net/http"

	"github.com/shulganew/GophKeeper/internal/api/oapi"
)

// sendKeeperError wraps sending of an error in the Error format, and
// handling the failure to marshal that.
func sendKeeperError(w http.ResponseWriter, code int, message string) {
	petErr := oapi.Error{
		Code:    int32(code),
		Message: message,
	}
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(petErr)
}
