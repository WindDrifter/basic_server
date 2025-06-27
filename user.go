package main

import (
	"encoding/json"
	"net/http"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/winddrifter/basic_server/internal/database"
	"github.com/winddrifter/basic_server/internal/auth"
	"os"
		"github.com/google/uuid"
		"time"

)

type User struct {
	ID        pgtype.UUID `json:"id"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
	Email     pgtype.Text    `json:"email"`
}
type parameters struct {
		Email pgtype.Text `json:"email"`
		Password string `json:"password"`
}

type loginParameters struct {
		Email pgtype.Text `json:"email"`
		Password string `json:"password"`
		ExpiresInSeconds  int32 `json:"expires_in_seconds"`
}

type userResponse struct {
	UserReturn
}

type UserReturn struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
	Email     pgtype.Text    `json:"email"`
}

type LoginReturn struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
	Email     pgtype.Text    `json:"email"`
	Token 	  string `json:"token"`
}

type LoginResponse struct {
	LoginReturn
}

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	
	decoder := json.NewDecoder(r.Body)
	type convertedParameters struct {
		Email pgtype.Text `json:"email"`
		PasswordHash string `json:"password"`
	}

	params := convertedParameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	hashed_password, pass_err := auth.HashPassword(params.PasswordHash)
	if pass_err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}
	params.PasswordHash = string(hashed_password)
	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams(params))
	
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, userResponse{
		UserReturn: UserReturn{
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
	err = cfg.db.DeleteAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error reseting chirps", err)
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

func (cfg *apiConfig) handlerLoginUser(rw http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := loginParameters{}
	err := decoder.Decode(&params)
	var duration time.Duration
	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)

	// prevent retrying password
	if err != nil {
		respondWithError(rw, http.StatusUnauthorized, "User not found", err)
		return
	}
	err = auth.CheckPasswordHash(user.PasswordHash, params.Password)
		// prevent retrying password
	if err != nil {
		respondWithError(rw, http.StatusUnauthorized, "Bad password", err)
		return
	}
	// if params.ExpiresInSeconds != nil {
	// 	timeInNanoSec := (params.ExpiresInSeconds * 1000 * 1000000)
	// 	duration =  time.Duration(timeInNanoSec)
	// } else {
		timeInNanoSec := (60 * 3600 * 1000 * 1000000)
		duration =  time.Duration(timeInNanoSec)
	// }
	token, tokenErr := auth.MakeJWT(user.ID, cfg.jwtSecret, duration)

	if tokenErr != nil {
		respondWithError(rw, http.StatusInternalServerError, "Something wrong when creating token", tokenErr)
		return
	}
	respondWithJSON(rw, http.StatusOK, LoginResponse{
		LoginReturn: LoginReturn{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
			Token: token,
		},
	})
}

