package store

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jay-bhogayata/blogapi/internal/config"
	"github.com/jay-bhogayata/blogapi/internal/database"
	openSearchClient "github.com/jay-bhogayata/blogapi/internal/opensearch"
)

type DB interface {
	Ping(context.Context) error
}

type BlogStoreType struct {
	DB               *pgxpool.Pool
	OpenSearchClient *openSearchClient.OpenSearch
	Query            *database.Queries
	Logger           *slog.Logger
	Config           *config.Config
}

var BlogStore *BlogStoreType

func UseBlogStore(handleConf *BlogStoreType) {
	BlogStore = &BlogStoreType{
		DB:               handleConf.DB,
		OpenSearchClient: handleConf.OpenSearchClient,
		Query:            handleConf.Query,
		Logger:           handleConf.Logger,
		Config:           handleConf.Config,
	}
}

type HealthResponse struct {
	Message string `json:"message"`
}

// func (h *BlogStore) CheckHealth(w http.ResponseWriter, r *http.Request) {

// 	err := h.DB.Ping(r.Context())
// 	if err != nil {
// 		h.Logger.Error("error while pining the db", err)
// 		http.Error(w, "something went wrong", http.StatusInternalServerError)
// 		return
// 	}

// 	hr := &HealthResponse{
// 		Message: "WORKING",
// 	}

// 	res, err := json.Marshal(hr)
// 	if err != nil {
// 		h.Logger.Error("error marshalling health response: ", "error", err)
// 		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	w.Write(res)
// }
