package main

import (
	"net/http"

	"github.com/jonathanfong098/csci169project/internal/database"
)

func (cfg *apiConfig) handlerSubscribeUser(w http.ResponseWriter, r *http.Request, user database.User) {
	if user.Subscribed {
		respondWithError(w, http.StatusBadRequest, "User is already subscribed")
		return
	}

	_, err := cfg.DB.SubscribeUser(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't subscribe user")
		return
	}

	err = cfg.SmtpServer.sendSubscribeEmail(user)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to send email")
		return
	}

	respondWithJSON(w, http.StatusOK, "User subscribed successfully")
}

func (cfg *apiConfig) handlerUnsubscribeUser(w http.ResponseWriter, r *http.Request, user database.User) {
	if !user.Subscribed {
		respondWithError(w, http.StatusBadRequest, "User is already unsubscribed")
		return
	}

	_, err := cfg.DB.SubscribeUser(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't unsubscribe user")
		return
	}

	err = cfg.SmtpServer.sendUnsubscribeEmail(user)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to send email")
		return
	}

	respondWithJSON(w, http.StatusOK, "User unsubscribed successfully")
}
