package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jay-bhogayata/blogapi/internal/database"
	"github.com/jay-bhogayata/blogapi/internal/pb/blogservice"
	"github.com/jay-bhogayata/blogapi/internal/store"
)

func (srv *Server) CreateBlog(ctx context.Context, blogPayload *blogservice.BlogPayload) (*blogservice.BlogId, error) {

	var UserID pgtype.UUID

	err := UserID.Scan(blogPayload.Blog.AuthorId)
	if err != nil {
		store.BlogStore.Logger.Error("invalid uuid found", "error", err)
		return nil, err
	}

	var TagID pgtype.Int4

	err = TagID.Scan(blogPayload.Blog.TagId)
	if err != nil {
		store.BlogStore.Logger.Error("invalid tag id found", "error", err)
		return nil, err
	}

	article := database.CreateArticleParams{
		Title:   blogPayload.Blog.Title,
		Content: blogPayload.Blog.Content,
		UserID:  UserID,
		TagID:   TagID,
	}

	res, err := store.BlogStore.Query.CreateArticle(ctx, article)
	if err != nil {
		store.BlogStore.Logger.Error("Error creating article", err)
		return nil, err
	}

	var stringUUID string

	if err := res.Scan(stringUUID); err != nil {
		return nil, err
	}
	return &blogservice.BlogId{
		Id: stringUUID,
	}, nil
}

// func PublishArticle(w http.ResponseWriter, r *http.Request) {

// 	articleID := chi.URLParam(r, "id")
// 	if articleID == "" {
// 		// h.respondWithError(w, http.StatusBadRequest, "Invalid request payload article id is missing")
// 		return
// 	}

// 	var articleId pgtype.UUID
// 	if err := articleId.Scan(articleID); err != nil {
// 		// h.respondWithError(w, http.StatusBadRequest, "Invalid article ID")
// 		return
// 	}

// 	// err := h.CheckUserOwnsArticle(r, articleId)
// 	// if err != nil {
// 	// h.respondWithError(w, http.StatusUnauthorized, err.Error())
// 	// return
// 	// }

// 	err := store.BlogStore.Query.PublishArticle(r.Context(), articleId)
// 	if err != nil {
// 		// h.respondWithError(w, http.StatusInternalServerError, "error while publishing article")
// 		return
// 	}

// 	// h.respondWithJSON(w, http.StatusOK, &Response{Message: "Article published successfully"})

// }

func (srv *Server) ListBlog(ctx context.Context, req *blogservice.ListBlogRequest) (*blogservice.ListBlogResponse, error) {
	pageNo := req.PageNumber
	if pageNo == 0 || pageNo < 0 {
		pageNo = 1
	}

	offset := (pageNo - 1) * 10

	articles, err := store.BlogStore.Query.GetAllArticles(ctx, offset)
	if err != nil {
		return nil, err
	}

	if len(articles) == 0 {
		// h.respondWithJSON(w, http.StatusOK, []database.Article{})
		return &blogservice.ListBlogResponse{}, nil
	}

	var blogList []*blogservice.Blog

	for _, article := range articles {
		var Id string

		err := article.ArticleID.Scan(&Id)
		if err != nil {
			store.BlogStore.Logger.Error("invalid uuid found", "error", err)
			return nil, err
		}

		var AuthorID string

		err = article.UserID.Scan(&AuthorID)
		if err != nil {
			store.BlogStore.Logger.Error("invalid uuid found", "error", err)
			return nil, err
		}

		blogList = append(blogList, &blogservice.Blog{
			Id:       Id,
			Title:    article.Title,
			Content:  article.Content,
			AuthorId: AuthorID,
			// TagId:    int32(article.TagID.Int),
		})
	}

	return &blogservice.ListBlogResponse{
		Blog: blogList,
	}, nil
}

func GetAllArticles(w http.ResponseWriter, r *http.Request) {
	pageNo, err := strconv.ParseInt(r.URL.Query().Get("page"), 10, 32)

	if pageNo == 0 || err != nil {
		pageNo = 1
	} else if pageNo < 0 {
		// h.respondWithError(w, http.StatusBadRequest, "Invalid page number")
		return
	}

	offset := (pageNo - 1) * 10

	articles, err := store.BlogStore.Query.GetAllArticles(r.Context(), int32(offset))
	if err != nil {
		return
	}

	if len(articles) == 0 {
		return
	}
}

func GetAllArticlesByUser(w http.ResponseWriter, r *http.Request) {
	// userID, err := h.ExtractUserIDFromJWT(r)
	// if err != nil {
	// h.Logger.Error("Error extracting user ID from JWT", err)
	// h.respondWithError(w, http.StatusUnauthorized, "Unauthorized")
	// return
	// }

	pageNo, err := strconv.ParseInt(r.URL.Query().Get("page"), 10, 32)
	if pageNo == 0 || err != nil {
		pageNo = 1
	} else if pageNo < 0 {
		// h.respondWithError(w, http.StatusBadRequest, "Invalid page number")
		return
	}

	// offSet := (pageNo - 1) * 10

	// articles, err := h.Query.GetAllArticleByUser(r.Context(), database.GetAllArticleByUserParams{
	// UserID: userID,
	// Offset: int32(offSet),
	// })

	if err != nil {
		// h.respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// if len(articles) == 0 {
	// h.respondWithJSON(w, http.StatusOK, []database.Article{})
	// return
	// }
	// h.respondWithJSON(w, http.StatusOK, articles)
}

func CheckUserOwnsArticle(r *http.Request, articleId pgtype.UUID) error {
	_, err := store.BlogStore.Query.GetUserIdByArticleId(r.Context(), articleId) //userid

	if err != nil {
		store.BlogStore.Logger.Error("Error fetching user id for article", err)
		return fmt.Errorf("error while fetching user id for article: %w", err)
	}

	// userID, err := h.ExtractUserIDFromJWT(r)
	// if err != nil {
	// h.Logger.Error("Error extracting user ID from JWT", err)
	// return errors.New("Unauthorized")
	// }

	// if userId != userID {
	// return errors.New("unauthorized for this article")
	// }

	return nil
}

func SearchArticle(w http.ResponseWriter, r *http.Request) {

	matchType := r.URL.Query().Get("match-type")
	if matchType == "" {
		matchType = "single"
	}

	term := r.URL.Query().Get("term")
	if term == "" {
		// h.respondWithError(w, http.StatusBadRequest, "Invalid request payload term is missing")
		return
	}

	fields := r.URL.Query().Get("fields")
	var fieldsArr []string

	if fields == "" {
		fields = "title"
	}

	fieldsArr = strings.Split(fields, ",")

	openSearchQuery, err := store.BlogStore.OpenSearchClient.QueryBuilder(matchType, fieldsArr, term)
	if err != nil {
		store.BlogStore.Logger.Error("Error building search query", err)
		// h.respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	result, err := store.BlogStore.OpenSearchClient.SearchQuery("blog-index", openSearchQuery, r.Context())

	if err != nil {
		store.BlogStore.Logger.Error("Error searching article", err)
		// h.respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	var searchResult map[string]interface{}
	decoder := json.NewDecoder(result.Body)
	err = decoder.Decode(&searchResult)

	if err != nil {
		store.BlogStore.Logger.Error("Error decoding search result", err)
		// h.respondWithError(w, http.StatusInternalServerError, "Internal Server Error: "+err.Error())
		return
	}

	// h.respondWithJSON(w, http.StatusOK, searchResult)

	err = result.Body.Close()
	if err != nil {
		store.BlogStore.Logger.Error("Error closing response body", err)
		// h.respondWithError(w, http.StatusInternalServerError, "Internal Server Error: "+err.Error())
		return
	}
}

func ExtractUserIDFromJWT(r *http.Request) (pgtype.UUID, error) {
	cookie, err := r.Cookie("jwt")
	if err != nil {
		return pgtype.UUID{}, errors.New("missing or invalid JWT")
	}

	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}
		return []byte(store.BlogStore.Config.JWTSecret), nil
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
