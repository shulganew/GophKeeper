package middlewares

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v4"

	"github.com/shulganew/GophKeeper/internal/app/config"
	"github.com/shulganew/GophKeeper/internal/entities"

	"go.uber.org/zap"
)

// Check JWT and read userID form it.
func Auth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// Check user JWT in Header.
		passVal := req.Context().Value(entities.CtxPassKey{})
		pass, ok := passVal.(string)
		if !ok {
			zap.S().Errorln("Can't get pass key from context.")
			h.ServeHTTP(res, req)
			return
		}
		jwt, hasJWT := GetHeaderJWT(req.Header)
		var userID uuid.UUID
		var err error
		if hasJWT {
			userID, err = GetUserIDJWT(jwt, pass)
			if err != nil {
				zap.S().Errorln("Can't get user UUID form JWT.", err)
				hasJWT = false
			}
		}
		ctx := context.WithValue(req.Context(), entities.AuthContext{}, entities.NewAuthContext(userID, jwt, hasJWT))
		h.ServeHTTP(res, req.WithContext(ctx))
	})
}

// Create JWT token.
func BuildJWTString(userID uuid.UUID, pass string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, entities.JWT{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.TokenExp)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(pass))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Retrive user's UUID from JWT string.
func GetUserIDJWT(tokenString string, pass string) (userID uuid.UUID, err error) {
	claims := &entities.JWT{}
	_, err = jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(pass), nil
	})

	return claims.UserID, err
}

// Create jwt token from string.
func GetJWT(tokenString string, pass string) (token *jwt.Token, err error) {
	claims := &entities.JWT{}
	token, err = jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(pass), nil
	})

	return token, err
}

// Check JWT is Set to Headek.
func GetHeaderJWT(header http.Header) (jwt string, isSet bool) {
	authHeader := header.Get("Authorization")
	if strings.HasPrefix(authHeader, config.AuthPrefix) {
		return authHeader[len(config.AuthPrefix):], true
	}
	return "", false

}
