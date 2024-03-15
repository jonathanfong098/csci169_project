package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) recommendFeeds(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Interest string
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	response, err := cfg.client.recommendFeed(params.Interest)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't recommend feed")
		return
	}
	respondWithJSON(w, http.StatusOK, response)
}
