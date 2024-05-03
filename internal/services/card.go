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

// Add new card.
func (k *Keeper) AddCard(w http.ResponseWriter, r *http.Request) {
	// Check registration.
	userID, isRegistered := CheckUserAuth(r.Context())
	if !isRegistered {
		http.Error(w, "JWT not found. Not authorized.", http.StatusUnauthorized)
		return
	}
	// Decode Card credentials from JSON.
	var newCard oapi.NewCard
	if err := json.NewDecoder(r.Body).Decode(&newCard); err != nil {
		sendKeeperError(w, http.StatusBadRequest, "Invalid format for NewCards")
		return
	}
	// Write data to storage.
	var db bytes.Buffer
	err := gob.NewEncoder(&db).Encode(&newCard)
	if err != nil {
		zap.S().Errorln("Error coding card to data: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	secretID, err := k.AddSecret(r.Context(), userID, entities.CARD, db.Bytes())
	if err != nil {
		zap.S().Errorln("Error adding card to DB: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return created card to client in responce (client add it to client's mem storage)
	card := oapi.Card{CardID: secretID.String(), Definition: newCard.Definition, Ccn: newCard.Ccn, Exp: newCard.Exp, Cvv: newCard.Cvv, Hld: newCard.Hld}
	w.Header().Add("Content-Type", "application/json")

	// set status code 201
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(card)
	if err != nil {
		zap.S().Errorln("Can't write to response in AddCard handler", err)
	}
	zap.S().Debugln("Card credentials added. ", card.CardID, " ", card.Definition)
}

// List all created cards.
func (k *Keeper) ListCards(w http.ResponseWriter, r *http.Request) {
	// Check registration.
	userID, isRegistered := CheckUserAuth(r.Context())
	if !isRegistered {
		http.Error(w, "JWT not found. Not authorized.", http.StatusUnauthorized)
		return
	}

	// Load all user's cards from database.
	secretDecoded, err := k.GetSecrets(r.Context(), userID, entities.CARD)
	if err != nil {
		zap.S().Errorln("Error getting card credentials: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Load decoded data and decode binary data to oapi.Card.
	var cards []oapi.Card
	for _, secret := range secretDecoded {
		var card oapi.Card
		err = gob.NewDecoder(bytes.NewReader(secret.Data)).Decode(&card)
		if err != nil {
			zap.S().Errorln("Error decode card to data: ", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		card.CardID = secret.SecretID.String()
		cards = append(cards, card)
	}

	w.Header().Add("Content-Type", "application/json")
	if len(cards) == 0 {
		zap.S().Infoln("No content.")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	// Set status code 200.
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(cards)
	if err != nil {
		zap.S().Errorln("Can't write to response in ListCards handler", err)
	}
}
