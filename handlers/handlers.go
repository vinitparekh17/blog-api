package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jay-bhogayata/blogapi/database"
	"github.com/jay-bhogayata/blogapi/logger"
)

type Handlers struct {
	DB    *pgxpool.Pool
	query *database.Queries
}

func NewHandlers(db *pgxpool.Pool, query *database.Queries) *Handlers {
	return &Handlers{
		DB:    db,
		query: query,
	}
}

type HealthResponse struct {
	Message string `json:"message"`
}

func (h *Handlers) CheckHealth(w http.ResponseWriter, r *http.Request) {

	err := h.DB.Ping(r.Context())
	if err != nil {
		logger.Log.Error("error while pining the db", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	hr := &HealthResponse{
		Message: "WORKING",
	}

	res, err := json.Marshal(hr)
	if err != nil {
		slog.Error("error marshalling health response: ", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
