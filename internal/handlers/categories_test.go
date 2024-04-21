package handlers_test

// import (
// 	"log/slog"
// 	"testing"

// 	"github.com/jay-bhogayata/blogapi/database"
// 	"github.com/jay-bhogayata/blogapi/handlers"
// )

// func TestParseID(t *testing.T) {

// 	db := new(MockDB)
// 	query := &database.Queries{}
// 	logger := slog.Default()

// 	h := handlers.NewHandlers(nil, db, query, logger)

// 	tests := []struct {
// 		name    string
// 		idStr   string
// 		want    int32
// 		wantErr bool
// 	}{
// 		{
// 			name:    "Valid ID",
// 			idStr:   "123",
// 			want:    123,
// 			wantErr: false,
// 		},
// 		{
// 			name:    "Invalid ID",
// 			idStr:   "abc",
// 			want:    0,
// 			wantErr: true,
// 		},
// 		{
// 			name:    "Empty ID",
// 			idStr:   "",
// 			want:    0,
// 			wantErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := h.ParseID(tt.idStr)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("Handlers.parseID() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if got != tt.want {
// 				t.Errorf("Handlers.parseID() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
