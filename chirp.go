package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
	"database/sql"
	
	"github.com/google/uuid"
    "github.com/mkling/GoWebServ/internal/database"
)

func checkProfanity(mess string) string {
	split := strings.Split(mess, " ")
	for i, word := range split {
		lowWord := strings.ToLower(word)
		if lowWord == "kerfuffle" || lowWord == "sharbert" || lowWord == "fornax" {
			split[i] = "****"
		}
	}
	return strings.Join(split, " ")
}

func (cfg *apiConfig) handlerAddChirp(w http.ResponseWriter, r *http.Request) {
    type parameters struct {
        Body   string `json:"body"`
        UserID string `json:"user_id"`
    }

    decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
        return
    }

    const maxChirpLength = 140
    if len(params.Body) > maxChirpLength {
        respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
        return
    }

    // Parse and validate user ID
    userUUID, err := uuid.Parse(params.UserID)
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid user ID format", err)
        return
    }

    // Clean profanity from body
    cleanedBody := checkProfanity(params.Body)

    // Create chirp using the correct parameter structure
    chirpParams := database.CreateChirpParams{
        Body:   cleanedBody,
        UserID: userUUID,
    }

    dbChirp, err := cfg.dbQueries.CreateChirp(r.Context(), chirpParams)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
        return
    }

    type chirpResponse struct {
        ID        uuid.UUID `json:"id"`
        CreatedAt time.Time `json:"created_at"`
        UpdatedAt time.Time `json:"updated_at"`
        Body      string    `json:"body"`
        UserID    uuid.UUID `json:"user_id"`
    }

    response := chirpResponse{
        ID:        dbChirp.ID,
        CreatedAt: dbChirp.CreatedAt,
        UpdatedAt: dbChirp.UpdatedAt,
        Body:      dbChirp.Body,
        UserID:    dbChirp.UserID,
    }

    respondWithJSON(w, http.StatusCreated, response)
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
    dbChirps, err := cfg.dbQueries.GetAllChirps(r.Context())
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
        return
    }

    type chirpResponse struct {
        ID        uuid.UUID `json:"id"`
        CreatedAt time.Time `json:"created_at"`
        UpdatedAt time.Time `json:"updated_at"`
        Body      string    `json:"body"`
        UserID    uuid.UUID `json:"user_id"`
    }

    chirps := make([]chirpResponse, len(dbChirps))
    for i, dbChirp := range dbChirps {
        chirps[i] = chirpResponse{
            ID:        dbChirp.ID,
            CreatedAt: dbChirp.CreatedAt,
            UpdatedAt: dbChirp.UpdatedAt,
            Body:      dbChirp.Body,
            UserID:    dbChirp.UserID,
        }
    }

    respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerGetChirpByID(w http.ResponseWriter, r *http.Request) {
	chirpIDStr := r.PathValue("chirpID")
	if chirpIDStr == "" {
        respondWithError(w, http.StatusBadRequest, "Missing chirp ID", nil)
        return
    }

    chirpID, err := uuid.Parse(chirpIDStr)
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid chirp ID format", err)
        return
    }

    dbChirp, err := cfg.dbQueries.GetChirp(r.Context(), chirpID)
    if err != nil {
        if err == sql.ErrNoRows {
            respondWithError(w, http.StatusNotFound, "Chirp not found", nil)
            return
        }
        respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirp", err)
        return
    }

    type chirpResponse struct {
        ID        uuid.UUID `json:"id"`
        CreatedAt time.Time `json:"created_at"`
        UpdatedAt time.Time `json:"updated_at"`
        Body      string    `json:"body"`
        UserID    uuid.UUID `json:"user_id"`
    }

    response := chirpResponse{
        ID:        dbChirp.ID,
        CreatedAt: dbChirp.CreatedAt,
        UpdatedAt: dbChirp.UpdatedAt,
        Body:      dbChirp.Body,
        UserID:    dbChirp.UserID,
    }

    respondWithJSON(w, http.StatusOK, response)
}