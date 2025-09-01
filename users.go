package main

import (
    "encoding/json"
    "net/http"

    "github.com/mkling/GoWebServ/internal/auth"
    "github.com/mkling/GoWebServ/internal/database"
)

func (cfg *apiConfig) handlerAddUsers(w http.ResponseWriter, r *http.Request) {
    type parameters struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
        return
    }

    hashedPassword, err := auth.HashPassword(params.Password)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
        return
    }

    createParams := database.CreateUserParams{
        Email:          params.Email,
        HashedPassword: hashedPassword,
    }

    dbUser, err := cfg.dbQueries.CreateUser(r.Context(), createParams)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
        return
    }

    user := User{
        ID:        dbUser.ID,
        CreatedAt: dbUser.CreatedAt,
        UpdatedAt: dbUser.UpdatedAt,
        Email:     dbUser.Email,
    }

    respondWithJSON(w, http.StatusCreated, user)
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
    type parameters struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
        return
    }

    dbUser, err := cfg.dbQueries.GetUserByEmail(r.Context(), params.Email)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
        return
    }

    err = auth.CheckPasswordHash(params.Password, dbUser.HashedPassword)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
        return
    }

    user:= User{
        ID:        dbUser.ID,
        CreatedAt: dbUser.CreatedAt,
        UpdatedAt: dbUser.UpdatedAt,
        Email:     dbUser.Email,
    }

    respondWithJSON(w, http.StatusOK, user)
}