package middlewares

import (
	"context"
	"net/http"

	"github.com/gofrs/uuid"

	"github.com/shulganew/GophKeeper/internal/services"

	"go.uber.org/zap"
)

// send pass to midleware.
type CtxPassKey struct{}

// Send values through middleware in context.

type AuthContext struct {
	userID       uuid.UUID
	isRegistered bool
}

func NewAuthContext(userID uuid.UUID, isRegistered bool) AuthContext {
	return AuthContext{userID: userID, isRegistered: isRegistered}
}

func (c AuthContext) GetUserID() uuid.UUID {
	return c.userID
}

func (c AuthContext) IsRegistered() bool {
	return c.isRegistered
}

// Check JWT and read userID form it.
func Auth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// Check user JWT in Header.
		passVal := req.Context().Value(CtxPassKey{})
		pass, ok := passVal.(string)
		if !ok {
			zap.S().Errorln("Can't git pass key from context.")
			h.ServeHTTP(res, req)
			return
		}
		jwt, hasJWT := services.GetHeaderJWT(req.Header)
		var userID uuid.UUID
		var err error
		if hasJWT {
			userID, err = services.GetUserIDJWT(jwt, pass)
			if err != nil {
				zap.S().Errorln("Can't get user UUID form JWT.", err)
				hasJWT = false
			}
		}
		ctx := context.WithValue(req.Context(), AuthContext{}, NewAuthContext(userID, hasJWT))
		h.ServeHTTP(res, req.WithContext(ctx))
	})
}
