package main

import (
	"encoding/json"
	"net/http"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
		"github.com/winddrifter/basic_server/internal/database"

	"os"
)

type User struct {
	ID        pgtype.UUID `json:"id"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
	Email     pgtype.Text    `json:"email"`
}

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email pgtype.Text `json:"email"`
	}
	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), params.Email)
	fmt.Println(err)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
	})
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	currentMode := os.Getenv("MODE")

	if currentMode != "Development" {
	
		respondWithError(w, http.StatusForbidden, "Not allow", nil)
		return
	
	}
	err := cfg.db.DeleteAllUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error reseting", err)
		return
	}
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

func (cfg *apiConfig) handlerAllUsers(rw http.ResponseWriter, r *http.Request) {
	type response struct {
		Users []database.User
	}
	users, err := cfg.db.GetAllUsers(r.Context())
	
	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, "Error reseting", err)
		return
	}
	respondWithJSON(rw, http.StatusOK, response{
		Users: users,
	})
	rw.Header().Add("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
}