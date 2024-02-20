package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jay-bhogayata/blogapi/database"
)

func (h *Handlers) CreateArticle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var article database.CreateArticleParams
	if err := json.NewDecoder(r.Body).Decode(&article); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	userID, err := h.extractUserIDFromJWT(r)
	if err != nil {
		h.logger.Error("Error extracting user ID from JWT", err)
		h.respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	article.UserID = userID

	res, err := h.query.CreateArticle(ctx, article)
	if err != nil {
		h.logger.Error("Error creating article", err)
		h.respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	response := database.Article{
		ArticleID:   res.ArticleID,
		Title:       res.Title,
		Content:     res.Content,
		UserID:      res.UserID,
		CategoryID:  res.CategoryID,
		CreatedAt:   res.CreatedAt,
		UpdatedAt:   res.UpdatedAt,
		IsPublished: res.IsPublished}

	h.respondWithJSON(w, http.StatusCreated, response)
}

func (h *Handlers) extractUserIDFromJWT(r *http.Request) (pgtype.UUID, error) {
	cookie, err := r.Cookie("jwt")
	if err != nil {
		return pgtype.UUID{}, errors.New("missing or invalid JWT")
	}

	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}
		return []byte(h.config.JWTSecret), nil
	})
	if err != nil {
		return pgtype.UUID{}, fmt.Errorf("error parsing token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return pgtype.UUID{}, errors.New("invalid JWT claims")
	}
	fmt.Println("claims: ", claims)

	userID, ok := claims["user_id"].(string)
	if !ok {
		return pgtype.UUID{}, errors.New("missing user_id in JWT claims")
	}

	var userIDUUID pgtype.UUID
	if err := userIDUUID.Scan(userID); err != nil {
		return pgtype.UUID{}, fmt.Errorf("error scanning user ID: %w", err)
	}

	return userIDUUID, nil
}

// package handlers

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"

// 	"github.com/golang-jwt/jwt/v5"
// 	"github.com/jackc/pgx/v5/pgtype"
// 	"github.com/jay-bhogayata/blogapi/database"
// )

// func (h *Handlers) CreateArticle(w http.ResponseWriter, r *http.Request) {

// 	article := database.CreateArticleParams{}

// 	err := json.NewDecoder(r.Body).Decode(&article)
// 	if err != nil {
// 		h.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
// 		return
// 	}

// 	cookie, err := r.Cookie("jwt")
// 	if err != nil {
// 		h.logger.Error("Error while getting cookie", err)
// 	}
// 	id, err := extractUserIDFromJWT(cookie.Value, h.config.JWTSecret)
// 	if err != nil {
// 		h.logger.Error("Error while extracting user id from jwt", err)
// 	}

// 	var userid pgtype.UUID
// 	userid.Scan(id)
// 	article.UserID = userid

// 	res, err := h.query.CreateArticle(r.Context(), article)
// 	if err != nil {
// 		h.respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
// 		return
// 	}

// 	var re database.Article
// 	re.ArticleID = res.ArticleID
// 	re.Title = res.Title
// 	re.Content = res.Content
// 	re.UserID = res.UserID
// 	re.CategoryID = res.CategoryID
// 	re.CreatedAt = res.CreatedAt
// 	re.IsPublished = res.IsPublished

// 	h.respondWithJSON(w, http.StatusCreated, re)
// }

// func extractUserIDFromJWT(tokenString string, secret string) (string, error) {
// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, fmt.Errorf("invalid signing method")
// 		}
// 		return []byte(secret), nil
// 	})
// 	if err != nil {
// 		return "", fmt.Errorf("error parsing token: %w", err)
// 	}

// 	claims, ok := token.Claims.(jwt.MapClaims)
// 	if !ok {
// 		return "", fmt.Errorf("error getting claims: %w", err)
// 	}

// 	userID, ok := claims["user_id"].(string)
// 	if !ok {
// 		return "", fmt.Errorf("error getting user id from claims: %w", err)
// 	}

// 	return userID, nil
// }
