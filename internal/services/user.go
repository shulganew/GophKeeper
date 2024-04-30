package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
	"github.com/shulganew/GophKeeper/internal/app/config"
	"github.com/shulganew/GophKeeper/internal/entities"
	"github.com/shulganew/GophKeeper/internal/rest/oapi"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// Register new user in Keeper.
func (k *Keeper) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user oapi.NewUser
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		zap.S().Errorln("Can't decode json")
		// If can't decode 400
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Set hash as user password.
	hash, err := k.HashPassword(user.Password)
	if err != nil {
		zap.S().Errorln("Error creating hash from password")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Add user to database.
	userID, err := k.stor.AddUser(r.Context(), user.Login, hash, user.Email)
	if err != nil {
		var pgErr *pq.Error
		// If URL exist in DataBase
		if errors.As(err, &pgErr) && pgerrcode.UniqueViolation == pgErr.Code {
			errt := "User's login is used"
			zap.S().Infoln(errt, err)
			http.Error(w, errt, http.StatusConflict)
			return
		}
		return
	}

	jwt, _ := BuildJWTString(*userID, k.conf.PassJWT)
	w.Header().Add("Content-Type", "text/plain")
	w.Header().Add("Authorization", config.AuthPrefix+jwt)

	// set status code 201
	w.WriteHeader(http.StatusCreated)

	_, err = w.Write([]byte("User added."))
	if err != nil {
		zap.S().Errorln("Can't write to response in SetUser handler", err)
	}

}

// Validate user in Keeper, if sucsess it return user's id.
func (k *Keeper) Login(w http.ResponseWriter, r *http.Request) {
	var oapiUser oapi.NewUser
	if err := json.NewDecoder(r.Body).Decode(&oapiUser); err != nil {
		// If can't decode 400
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get User from storage
	dbUser, err := k.stor.GetByLogin(r.Context(), oapiUser.Login)
	zap.S().Infof("User form db: %v \n", dbUser)
	if err != nil {
		zap.S().Infoln("User not found by login. ", err)
		http.Error(w, "Wrong login or password", http.StatusUnauthorized)
		return
	}

	// Check pass is correct
	err = k.CheckPassword(oapiUser.Password, dbUser.PassHash)
	if err != nil {
		http.Error(w, "Wrong login or password", http.StatusUnauthorized)
	}

	zap.S().Debug("Login sucsess, user id is: ", dbUser.UUID)
	jwt, _ := BuildJWTString(dbUser.UUID, k.conf.PassJWT)

	w.Header().Add("Content-Type", "text/plain")
	w.Header().Add("Authorization", config.AuthPrefix+jwt)

	// set status code 200
	w.WriteHeader(http.StatusOK)

	_, err = w.Write([]byte("User loged in."))
	if err != nil {
		zap.S().Errorln("Can't write to response in LoginUser handler", err)
	}
}

// HashPassword returns the bcrypt hash of the password.
func (k Keeper) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

// CheckPassword checks if the provided password is correct or not.
func (k Keeper) CheckPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
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

// Check if contex has JWT valid user token from auth middleware.
func CheckUserAuth(ctx context.Context) (userID string, isRegistered bool) {
	auth := ctx.Value(entities.AuthContext{}).(entities.AuthContext)
	zap.S().Debugf("UserID: %s, JWT: %s, is registered: %t \n", auth.GetUserID(), auth.GetUserJWT(), auth.IsRegistered)
	return auth.GetUserID().String(), auth.IsRegistered()
}
