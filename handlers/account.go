package handlers

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jay-bhogayata/blogapi/database"
	"github.com/jay-bhogayata/blogapi/mailer"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handlers) RegisterUser(w http.ResponseWriter, r *http.Request) {
	u := database.CreateUserParams{}
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	hashedPassword, err := h.HashPassword(u.PasswordHash)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	token, _ := h.GenerateToken()

	user, err := h.query.CreateUser(r.Context(), database.CreateUserParams{Username: u.Username, Email: u.Email, PasswordHash: hashedPassword, IsVerified: pgtype.Bool{Bool: false, Valid: true}, VerificationToken: pgtype.Text{String: token, Valid: true}})
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	activationLink := fmt.Sprintf("http://localhost:8080/api/v1/accounts/verify?token=%s", token)

	body, err := mailer.SetupVerificationTemplate(u.Username, activationLink)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	err = mailer.SendEmail(u.Email, "Account Verification", body)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	h.respondWithJSON(w, http.StatusCreated, user)
}

func (h *Handlers) HashPassword(password string) (hashOfPassword string, err error) {

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		h.logger.Error("error hashing password", "error", err)
		return "", err
	}

	return string(bytes), err

}

func (h *Handlers) MatchPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false
	}
	return err == nil
}

func (h *Handlers) GenerateToken() (token string, err error) {
	randBytes := make([]byte, 16)
	_, err = rand.Read(randBytes)
	if err != nil {
		h.logger.Error("error generating token", "error", err)
		return "", err
	}
	return fmt.Sprintf("%x", randBytes), nil

}

func (h *Handlers) VerifyUser(w http.ResponseWriter, r *http.Request) {

	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	u, err := h.query.GetUserByVerificationToken(r.Context(), pgtype.Text{String: token, Valid: true})
	if err != nil {
		h.logger.Error("error getting user by verification token", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = h.query.VerifyUser(r.Context(), u.VerificationToken)
	if err != nil {
		h.logger.Error("error verifying user", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(`<html><body><h1>Account Verified</h1></body></html>`))
}
