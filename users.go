package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handlerAddUsers(w http.ResponseWriter, r *http.Request) {
	
	type parameters struct {
		Email		string		`json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	dbUser, err := cfg.dbQueries.CreateUser(r.Context(), params.Email)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
        return
    }

	user := User{
        ID:        dbUser.ID,
        CreatedAt: dbUser.CreatedAt,
        UpdatedAt: dbUser.UpdatedAt,
        Email:     dbUser.Email.(string),
    }

	respondWithJSON(w, 201, user)

}