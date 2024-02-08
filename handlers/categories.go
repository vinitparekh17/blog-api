package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jay-bhogayata/blogapi/database"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

type Category struct {
	Name string `json:"name"`
}

func (h *Handlers) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.query.GetAllCategories(r.Context())
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusOK, categories)
}

func (h *Handlers) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	id, err := h.ParseID(chi.URLParam(r, "id"))
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	category, err := h.query.GetCategoryById(r.Context(), id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			h.respondWithJSON(w, http.StatusNotFound, &ErrorResponse{Message: "No category found with given id"})
			return
		}
		h.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusOK, category)
}

func (h *Handlers) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var category Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		h.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.query.CreateCategory(r.Context(), category.Name)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusCreated, id)
}

func (h *Handlers) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	var category database.UpdateCategoryParams
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		h.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	updatedCategory, err := h.query.UpdateCategory(r.Context(), category)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusOK, updatedCategory)
}

func (h *Handlers) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id, err := h.ParseID(chi.URLParam(r, "id"))
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	type res struct {
		Message string `json:"message"`
	}

	if _, err := h.query.GetCategoryById(r.Context(), id); err != nil {
		if err.Error() == "no rows in result set" {
			h.respondWithJSON(w, http.StatusNotFound, &ErrorResponse{Message: "No category found with given id"})
			return
		}
		h.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := h.query.DeleteCategory(r.Context(), id); err != nil {
		h.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusOK, &res{Message: "Category deleted successfully"})
}

func (h *Handlers) respondWithError(w http.ResponseWriter, code int, message string) {
	h.logger.Error(message)
	h.respondWithJSON(w, code, &ErrorResponse{Message: message})
}

func (h *Handlers) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		h.logger.Error(err.Error())
		h.respondWithError(w, http.StatusInternalServerError, "error while marshalling response")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (h *Handlers) ParseID(idStr string) (int32, error) {
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		h.logger.Error(err.Error())
		return 0, err
	}
	return int32(id), nil
}
