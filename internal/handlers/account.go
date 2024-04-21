package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jay-bhogayata/blogapi/internal/database"
	"github.com/jay-bhogayata/blogapi/internal/helper"
	"github.com/jay-bhogayata/blogapi/internal/mailer"
)

type User struct {
	UserName string      `json:"username"`
	UserID   pgtype.UUID `json:"user_id"`
	Email    string      `json:"email"`
}

type Response struct {
	Message string `json:"message"`
}

func (h *Handlers) RegisterUser(w http.ResponseWriter, r *http.Request) {

	type ReqUser struct {
		UserName string `json:"username" required:"true"`
		Email    string `json:"email" required:"true"`
		Password string `json:"password" required:"true"`
	}
	u := ReqUser{}

	err := helper.DecodeJSONBody(w, r, &u)
	if err != nil {

		var mr *helper.MalformedRequest

		if errors.As(err, &mr) {
			h.Logger.Error("error in decoding json body", "error", err.Error())
			h.respondWithError(w, mr.Status, mr.Msg)
		} else {
			h.Logger.Error(err.Error())
			h.respondWithError(w, mr.Status, http.StatusText(http.StatusInternalServerError))
		}
		return
	}

	hash, err := argon2id.CreateHash(u.Password, argon2id.DefaultParams)
	if err != nil {
		h.Logger.Error("error in hashing the password", "error:", err)
		h.respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	token, err := helper.GenerateToken()
	if err != nil {
		h.Logger.Error("error in generating verification token", "error", err)
		h.respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	usr, err := h.Query.CreateUser(r.Context(), database.CreateUserParams{Username: u.UserName, Email: u.Email, PasswordHash: hash, IsVerified: pgtype.Bool{Bool: false, Valid: true}, VerificationToken: pgtype.Text{String: token, Valid: true}})
	if err != nil {

		if strings.Contains(err.Error(), "duplicate key") && strings.Contains(err.Error(), "users_email_key") {
			h.respondWithError(w, http.StatusBadRequest, "email is already taken try to register with different email")
			return
		}

		h.Logger.Error("error in creating user in database", "error", err)
		h.respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	activationLink := fmt.Sprintf("http://%s/api/v1/users/verify?token=%s", r.Host, token)
	body, err := mailer.SetupVerificationTemplate(usr.Username, activationLink)
	if err != nil {
		h.Logger.Error("error in setting verification template", err)
		h.respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	err = mailer.SendEmail(usr.Email, "Account Verification", body, h.Config.EmailSender)
	if err != nil {
		h.Logger.Error("error in sending email to", usr.Email, err)
		h.respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	msg := fmt.Sprintf("account created successfully , verification link has been sent to %s", usr.Email)
	res := Response{Message: msg}

	h.respondWithJSON(w, http.StatusCreated, res)
}

func (h *Handlers) VerifyUser(w http.ResponseWriter, r *http.Request) {

	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	u, err := h.Query.GetUserByVerificationToken(r.Context(), pgtype.Text{String: token, Valid: true})
	if err != nil {
		h.Logger.Error("error getting user by verification token", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = h.Query.VerifyUser(r.Context(), u.VerificationToken)
	if err != nil {
		h.Logger.Error("error verifying user", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(`<html><body><h1>Account Verified</h1></body></html>`))
}

func (h *Handlers) LoginUser(w http.ResponseWriter, r *http.Request) {

	type user struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	u := user{}

	err := helper.DecodeJSONBody(w, r, &u)
	if err != nil {
		h.Logger.Error("error in decoding json body", "error", err)
		h.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	usr, err := h.Query.GetUserByEmail(r.Context(), u.Email)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			h.respondWithError(w, http.StatusUnauthorized, "User with this email does not exist")
			return
		}
		h.Logger.Error("error in getting user by email", "error", err)
		h.respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if !usr.IsVerified.Bool {
		h.respondWithError(w, http.StatusUnauthorized, "User is not verified")
		return
	}

	match, err := argon2id.ComparePasswordAndHash(u.Password, usr.PasswordHash)
	if err != nil {
		fmt.Println(err)
		h.respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	if !match {
		h.respondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  usr.UserID,
		"email":    usr.Email,
		"username": usr.Username,
		"expiry":   time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte(h.Config.JWTSecret))
	if err != nil {
		h.Logger.Error("error in signing JWT token", "error", err)
		h.respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
	})

	h.respondWithJSON(w, http.StatusOK, &Response{Message: "Logged in successfully"})
}

func (h *Handlers) LogoutUser(w http.ResponseWriter, r *http.Request) {

	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
	})

	h.respondWithJSON(w, http.StatusOK, &Response{Message: "Logged out successfully"})
}

func (h *Handlers) GetUserInfoByUserEmail(w http.ResponseWriter, r *http.Request) {
	mail := chi.URLParam(r, "mail")

	usr, err := h.Query.GetUserByEmail(r.Context(), mail)
	if err != nil {
		h.Logger.Error("error in getting user by email", "error", err)
		h.respondWithError(w, http.StatusInternalServerError, "error in getting user by email")
		return
	}

	user := User{UserName: usr.Username, UserID: usr.UserID, Email: usr.Email}

	h.respondWithJSON(w, http.StatusOK, user)
}

func (h *Handlers) GetAllUsers(w http.ResponseWriter, r *http.Request) {

	user, err := h.Query.GetAllUsers(r.Context())
	if err != nil {
		h.Logger.Error("error in getting all users", "error", err)
		h.respondWithError(w, http.StatusInternalServerError, "error in getting all users")
		return
	}

	res, err := json.Marshal(user)
	if err != nil {
		h.Logger.Error("error in marshalling users", "error", err)
		h.respondWithError(w, http.StatusInternalServerError, "error in marshalling users")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func (h *Handlers) DeleteUser(w http.ResponseWriter, r *http.Request) {

	userID, err := h.ExtractUserIDFromJWT(r)

	if err != nil {
		h.Logger.Error("Error extracting user ID from JWT", err)
		h.respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	_, err = h.Query.DeleteUser(r.Context(), userID)

	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			h.respondWithError(w, http.StatusNotFound, "User with this ID does not exist")
			return
		}
		h.Logger.Error(err.Error())
		h.respondWithError(w, http.StatusInternalServerError, "error while deleting user")
		return
	}

	h.respondWithJSON(w, http.StatusOK, "User has been deleted")
}
