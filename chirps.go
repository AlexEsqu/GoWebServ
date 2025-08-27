package main

import (
	"net/http"
	"time"
	"encoding/json"
	"log"
)



func writeResponse(w http.ResponseWriter, r *http.Request, httpReturn int, errorMesg string, isValid bool){
    type returnVals struct {
        // the key will be the name of struct field unless you give it an explicit JSON tag
        CreatedAt time.Time `json:"created_at"`
        ID int `json:"id"`
		Error string `json:"error, omitiempty"`
		Valid bool `json:"valid"`
    }
    respBody := returnVals{
        CreatedAt: time.Now(),
        ID: 123,
		Error: errorMesg,
		Valid: isValid,
    }
    responseJSON, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
    w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpReturn)
    w.Write(responseJSON)
}

func (cfg *apiConfig) handlerValidateChirp(w http.ResponseWriter, r *http.Request){

	const genericError = "Something went wrong"
	const tooLongError = "Chirp is too long"

	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		writeResponse(w, r, 500, genericError, false)
	} else if len(params.Body) > 140 {
		writeResponse(w, r, 400, tooLongError, false)
	} else {
		writeResponse(w, r, 200, "", true)
	}
}

