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

// Add new card.
func (k *Keeper) AddCard(w http.ResponseWriter, r *http.Request) {
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
	var newCard oapi.NewCard
	err = json.Unmarshal(body, &newCard)
	if err != nil {
		zap.S().Errorln("Can't Read json: ", err)
		http.Error(w, "Can't Read json.", http.StatusInternalServerError)
	}

	secretID, err := k.AddSecret(r.Context(), userID, entities.CARD, body)
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
	// Get userID from jwt.
	userID, err := jwt.GetUserID(k.ua, r)
	if err != nil {
		zap.S().Errorln("Error getting userID: ", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
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
	cards := make(map[string]oapi.Card, len(secretDecoded))
	for _, secret := range secretDecoded {
		var card oapi.Card
		err = json.NewDecoder(bytes.NewReader(secret.Data)).Decode(&card)
		if err != nil {
			zap.S().Errorln("Error decode card to data: ", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		card.CardID = secret.SecretID.String()
		cards[card.CardID] = card
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

func (k *Keeper) UpdateCard(w http.ResponseWriter, r *http.Request) {
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
	var card oapi.Card
	err = json.Unmarshal(body, &card)
	if err != nil {
		zap.S().Errorln("Can't Read json: ", err)
		http.Error(w, "Can't Read json.", http.StatusInternalServerError)
	}

	err = k.UpdateSecret(r.Context(), userID, entities.SITE, body, card.CardID)
	if err != nil {
		zap.S().Errorln("Error adding site to DB: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// set status code 200
	w.WriteHeader(http.StatusOK)
}

// Delelete card data.
func (k *Keeper) DelCard(w http.ResponseWriter, r *http.Request, cardID string) {
	err := k.DeleteSecret(r.Context(), cardID)
	if err != nil {
		zap.S().Errorln("Error adding site to DB: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// set status code 200
	w.WriteHeader(http.StatusOK)
}
