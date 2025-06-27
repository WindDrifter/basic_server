package main

import (
	"encoding/json"
	"net/http"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/winddrifter/basic_server/internal/database"
	"github.com/google/uuid"
	"github.com/winddrifter/basic_server/internal/auth"
	"fmt"
)

type Chirp struct {
	ID        pgtype.UUID   `json:"id"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	UpdatedAt pgtype.Timestamp  `json:"updated_at"`
	Body      pgtype.Text `json:"body"`
	UserID    pgtype.UUID `json:"user_id"`
}

type ReturnChirp struct {
	ID        uuid.UUID   `json:"id"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	UpdatedAt pgtype.Timestamp  `json:"updated_at"`
	Body      pgtype.Text `json:"body"`
	UserID    pgtype.UUID `json:"user_id"`
}



func (cfg *apiConfig) handlerChirpCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   pgtype.Text `json:"body"`
		UserID pgtype.UUID `json:"user_id"`
	}
	type response struct {
		ReturnChirp
	}
	headers := r.Header
	token, tokenErr := auth.GetBearerToken(headers)

	if tokenErr != nil {
		respondWithError(w, http.StatusUnauthorized, "Cannot access", tokenErr)
		return
	}
	userId, validErr := auth.ValidateJWT(token, cfg.jwtSecret)

	var validPGUUID pgtype.UUID;
	validPGUUID.Scan(userId.String())
	if validErr != nil {
		respondWithError(w, http.StatusUnauthorized, "Cannot validate token", validErr)
		return
	}
	fmt.Printf("%v", validPGUUID)

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	fmt.Printf("%v", userId)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	params.UserID = validPGUUID

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams(params))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		ReturnChirp: ReturnChirp{
			ID: chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			UserID: chirp.UserID,
			Body: chirp.Body,
		},
	})
}

func (cfg *apiConfig) handlerAllChirps(rw http.ResponseWriter, r *http.Request) {
	// type response struct {
	// 	Chirps []Chirp
	// }

	chirps, err := cfg.db.GetAllChirps(r.Context())
	
	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, "Error reseting", err)
		return
	}
	returnVal := []ReturnChirp{};
	 for _, v := range chirps {
		returnVal = append(returnVal,  convertToResponseStruct(v))
	 }
	respondWithJSON(rw, http.StatusOK, returnVal)
	rw.Header().Add("Content-Type", "text/plain; charset=utf-8")
}

func (cfg *apiConfig) handlerGetChirp(rw http.ResponseWriter, r *http.Request) {
	type response struct {
		ReturnChirp
	}
	chrip_id := r.PathValue("chirpID")
	var input uuid.UUID;
	input.Scan(chrip_id)
	chirp, err := cfg.db.GetChirpById(r.Context(), input)
	
	if err != nil {
		respondWithError(rw, http.StatusNotFound, "Not Found", err)
		return
	}
	returnVal := convertToResponseStruct(chirp)
	respondWithJSON(rw, http.StatusOK, response{ReturnChirp: returnVal})
	rw.Header().Add("Content-Type", "text/plain; charset=utf-8")
}

func convertToResponseStruct(input database.Chirp) ReturnChirp {
	return ReturnChirp{
		ID: input.ID,
		CreatedAt: input.CreatedAt,
		UpdatedAt: input.UpdatedAt,
		Body: input.Body,
		UserID: input.UserID,
	}

}