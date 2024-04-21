package handlers_test

// import (
// 	"context"
// 	"errors"
// 	"log/slog"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/jay-bhogayata/blogapi/database"
// 	"github.com/jay-bhogayata/blogapi/handlers"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// )

// type DB interface {
// 	Ping(context.Context) error
// }

// type MockDB struct {
// 	mock.Mock
// }

// func (db *MockDB) Ping(ctx context.Context) error {
// 	args := db.Called(ctx)
// 	return args.Error(0)
// }

// func TestCheckHealth(t *testing.T) {
// 	db := new(MockDB)
// 	query := &database.Queries{}
// 	logger := slog.Default()

// 	h := handlers.NewHandlers(nil, db, query, logger)

// 	req, err := http.NewRequest("GET", "/health", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	rr := httptest.NewRecorder()

// 	handler := http.HandlerFunc(h.CheckHealth)

// 	db.On("Ping", mock.Anything).Return(nil)

// 	handler.ServeHTTP(rr, req)

// 	assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")

// 	expected := `{"message":"WORKING"}`
// 	assert.Contains(t, rr.Body.String(), expected, "handler returned unexpected body")
// }

// func TestCheckHealth_Error(t *testing.T) {
// 	db := new(MockDB)
// 	query := &database.Queries{}
// 	logger := slog.Default()

// 	h := handlers.NewHandlers(nil, db, query, logger)

// 	req, err := http.NewRequest("GET", "/health", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	rr := httptest.NewRecorder()

// 	handler := http.HandlerFunc(h.CheckHealth)

// 	db.On("Ping", mock.Anything).Return(errors.New("database error"))

// 	handler.ServeHTTP(rr, req)

// 	assert.Equal(t, http.StatusInternalServerError, rr.Code, "handler returned wrong status code")

// 	expected := "something went wrong"
// 	assert.Contains(t, rr.Body.String(), expected, "handler returned unexpected body")

// }
