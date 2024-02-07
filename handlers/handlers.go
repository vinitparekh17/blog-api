package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/jay-bhogayata/blogapi/database"
)

type DB interface {
	Ping(context.Context) error
}

type Handlers struct {
	DB     DB
	query  *database.Queries
	logger *slog.Logger
}

func NewHandlers(db DB, query *database.Queries, logger *slog.Logger) *Handlers {
	return &Handlers{
		DB:     db,
		query:  query,
		logger: logger,
	}
}

type HealthResponse struct {
	Message string `json:"message"`
}

func (h *Handlers) CheckHealth(w http.ResponseWriter, r *http.Request) {

	err := h.DB.Ping(r.Context())
	if err != nil {
		h.logger.Error("error while pining the db", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	hr := &HealthResponse{
		Message: "WORKING",
	}

	res, err := json.Marshal(hr)
	if err != nil {
		h.logger.Error("error marshalling health response: ", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
