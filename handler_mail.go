package main

import (
	"encoding/json"
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

	type parameters struct {
		Summarize *bool `json:"summarize"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	if params.Summarize != nil && *params.Summarize {
		_, err = cfg.DB.SummarizePostsOn(r.Context(), user.ID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't turn on posts summarization")
			return
		}
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

	if !user.Summarize {
		_, err := cfg.DB.SummarizePostsOff(r.Context(), user.ID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't turn off posts summarization")
			return
		}
	}

	_, err := cfg.DB.UnsubscribeUser(r.Context(), user.ID)
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

func (cfg *apiConfig) handlerToggleSummarizePosts(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Summarize *bool `json:"summarize"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	if params.Summarize == nil {
		respondWithError(w, http.StatusBadRequest, "Missing summarize parameter")
		return
	}

	if *params.Summarize && user.Summarize {
		respondWithError(w, http.StatusBadRequest, "Posts summarization is already on")
		return
	} else if (!*params.Summarize) && (!user.Summarize) {
		respondWithError(w, http.StatusBadRequest, "Posts summarization is already off")
		return
	}

	if *params.Summarize {
		_, err = cfg.DB.SummarizePostsOn(r.Context(), user.ID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't turn on posts summarization")
			return
		}
	} else {
		_, err = cfg.DB.SummarizePostsOff(r.Context(), user.ID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't turn off posts summarization")
			return
		}
	}

	respondWithJSON(w, http.StatusOK, "Posts summarization toggled successfully")
}
