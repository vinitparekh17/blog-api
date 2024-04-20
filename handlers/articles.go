package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

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

	userID, err := h.ExtractUserIDFromJWT(r)
	if err != nil {
		h.Logger.Error("Error extracting user ID from JWT", err)
		h.respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	article.UserID = userID

	fmt.Println(article)

	res, err := h.Query.CreateArticle(ctx, article)
	if err != nil {
		h.Logger.Error("Error creating article", err)
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

	err := h.CheckUserOwnsArticle(r, articleId)
	if err != nil {
		h.respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	err = h.Query.PublishArticle(r.Context(), articleId)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "error while publishing article")
		return
	}

	h.respondWithJSON(w, http.StatusOK, &Response{Message: "Article published successfully"})

}

func (h *Handlers) GetAllArticles(w http.ResponseWriter, r *http.Request) {
	pageNo, err := strconv.ParseInt(r.URL.Query().Get("page"), 10, 32)

	if pageNo == 0 || err != nil {
		pageNo = 1
	} else if pageNo < 0 {
		h.respondWithError(w, http.StatusBadRequest, "Invalid page number")
		return
	}

	offset := (pageNo - 1) * 10

	articles, err := h.Query.GetAllArticles(r.Context(), int32(offset))
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if len(articles) == 0 {
		h.respondWithJSON(w, http.StatusOK, []database.Article{})
		return
	}

	h.respondWithJSON(w, http.StatusOK, articles)
}

func (h *Handlers) GetAllArticlesByUser(w http.ResponseWriter, r *http.Request) {
	userID, err := h.ExtractUserIDFromJWT(r)
	if err != nil {
		h.Logger.Error("Error extracting user ID from JWT", err)
		h.respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	pageNo, err := strconv.ParseInt(r.URL.Query().Get("page"), 10, 32)
	if pageNo == 0 || err != nil {
		pageNo = 1
	} else if pageNo < 0 {
		h.respondWithError(w, http.StatusBadRequest, "Invalid page number")
		return
	}

	offSet := (pageNo - 1) * 10

	articles, err := h.Query.GetAllArticleByUser(r.Context(), database.GetAllArticleByUserParams{
		UserID: userID,
		Offset: int32(offSet),
	})

	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if len(articles) == 0 {
		h.respondWithJSON(w, http.StatusOK, []database.Article{})
		return
	}
	h.respondWithJSON(w, http.StatusOK, articles)
}

func (h *Handlers) CheckUserOwnsArticle(r *http.Request, articleId pgtype.UUID) error {
	userId, err := h.Query.GetUserIdByArticleId(r.Context(), articleId)

	if err != nil {
		h.Logger.Error("Error fetching user id for article", err)
		return fmt.Errorf("error while fetching user id for article: %w", err)
	}

	userID, err := h.ExtractUserIDFromJWT(r)
	if err != nil {
		h.Logger.Error("Error extracting user ID from JWT", err)
		return errors.New("Unauthorized")
	}

	if userId != userID {
		return errors.New("unauthorized for this article")
	}

	return nil
}

func (h *Handlers) SearchArticle(w http.ResponseWriter, r *http.Request) {

	matchType := r.URL.Query().Get("match-type")
	if matchType == "" {
		matchType = "single"
	}

	term := r.URL.Query().Get("term")
	if term == "" {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request payload term is missing")
		return
	}

	fields := r.URL.Query().Get("fields")
	var fieldsArr []string

	if fields == "" {
		fields = "title"
	}

	fieldsArr = strings.Split(fields, ",")

	openSearchQuery, err := h.OpenSearchClient.QueryBuilder(matchType, fieldsArr, term)
	if err != nil {
		h.Logger.Error("Error building search query", err)
		h.respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	result, err := h.OpenSearchClient.SearchQuery("blog-index", openSearchQuery, r.Context())

	if err != nil {
		h.Logger.Error("Error searching article", err)
		h.respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	var searchResult map[string]interface{}
	decoder := json.NewDecoder(result.Body)
	err = decoder.Decode(&searchResult)

	if err != nil {
		h.Logger.Error("Error decoding search result", err)
		h.respondWithError(w, http.StatusInternalServerError, "Internal Server Error: "+err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusOK, searchResult)

	err = result.Body.Close()
	if err != nil {
		h.Logger.Error("Error closing response body", err)
		h.respondWithError(w, http.StatusInternalServerError, "Internal Server Error: "+err.Error())
		return
	}
}

func (h *Handlers) ExtractUserIDFromJWT(r *http.Request) (pgtype.UUID, error) {
	cookie, err := r.Cookie("jwt")
	if err != nil {
		return pgtype.UUID{}, errors.New("missing or invalid JWT")
	}

	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}
		return []byte(h.Config.JWTSecret), nil
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
