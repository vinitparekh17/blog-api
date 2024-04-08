package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
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

	// userID, err := h.extractUserIDFromJWT(r)
	// if err != nil {
	// 	h.logger.Error("Error extracting user ID from JWT", err)
	// 	h.respondWithError(w, http.StatusUnauthorized, "Unauthorized")
	// 	return
	// }

	// article.UserID = userID

	res, err := h.query.CreateArticle(ctx, article)
	if err != nil {
		h.logger.Error("Error creating article", err)
		h.respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	type response struct {
		Message   string      `json:"message"`
		ArticleID pgtype.UUID `json:"article_id"`
	}

	h.respondWithJSON(w, http.StatusCreated, response{
		Message:   "Article created successfully",
		ArticleID: res,
	})
}

func (h *Handlers) PublishArticle(w http.ResponseWriter, r *http.Request) {

	articleID := chi.URLParam(r, "id")
	if articleID == "" {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request payload article id is missing")
		return
	}

	var articleId pgtype.UUID
	if err := articleId.Scan(articleID); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid article ID")
		return
	}

	err := h.checkUserOwnsArticle(r, articleId)
	if err != nil {
		h.respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	err = h.query.PublishArticle(r.Context(), articleId)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "error while publishing article")
		return
	}

	h.respondWithJSON(w, http.StatusOK, &Response{Message: "Article published successfully"})

}

func (h *Handlers) GetAllArticles(w http.ResponseWriter, r *http.Request) {
	articles, err := h.query.GetAllArticles(r.Context())
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	h.respondWithJSON(w, http.StatusOK, articles)
}

func (h *Handlers) GetAllArticlesByUser(w http.ResponseWriter, r *http.Request) {
	userID, err := h.extractUserIDFromJWT(r)
	if err != nil {
		h.logger.Error("Error extracting user ID from JWT", err)
		h.respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	articles, err := h.query.GetAllArticleByUser(r.Context(), userID)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	h.respondWithJSON(w, http.StatusOK, articles)
}

// utils

func (h *Handlers) checkUserOwnsArticle(r *http.Request, articleId pgtype.UUID) error {
	userId, err := h.query.GetUserIdByArticleId(r.Context(), articleId)

	if err != nil {
		h.logger.Error("Error fetching user id for article", err)
		return fmt.Errorf("error while fetching user id for article: %w", err)
	}

	userID, err := h.extractUserIDFromJWT(r)
	if err != nil {
		h.logger.Error("Error extracting user ID from JWT", err)
		return errors.New("Unauthorized")
	}

	if userId != userID {
		return errors.New("unauthorized for this article")
	}

	return nil
}

func (h *Handlers) extractUserIDFromJWT(r *http.Request) (pgtype.UUID, error) {
	cookie, err := r.Cookie("jwt")
	if err != nil {
		return pgtype.UUID{}, errors.New("missing or invalid JWT")
	}

	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (any, error) {
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
