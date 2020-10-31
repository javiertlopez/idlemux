package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// Videos interface, for testing purposes
type Videos interface {
	Insert(ctx context.Context, anyVideo *Video) (*Video, error)
}

// Video struct
type Video struct {
	ID          *string  `json:"id,omitempty"`
	Title       string   `json:"title,omitempty"`
	Description string   `json:"description,omitempty"`
	Duration    *float64 `json:"duration,omitempty"`
	CreatedAt   string   `json:"created_at,omitempty"`
	UpdatedAt   string   `json:"updated_at,omitempty"`
}

// addVideoHandler adds the handler to the mux router
func addVideoHandler(r *mux.Router) {
	r.HandleFunc("/videos", CreateVideoHandler).Methods("POST")
}

// CreateVideoHandler handler
func CreateVideoHandler(w http.ResponseWriter, r *http.Request) {
	var video Video
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&video); err != nil {
		JSONResponse(
			w, http.StatusBadRequest,
			Response{
				Message: "Invalid request",
				Status:  http.StatusBadRequest,
			},
		)
		return
	}
	defer r.Body.Close()

	err := fmt.Errorf("handler not implemented")

	if err != nil {
		JSONResponse(
			w, http.StatusInternalServerError,
			Response{
				Message: err.Error(),
				Status:  http.StatusBadRequest,
			},
		)

		return
	}

	JSONResponse(
		w,
		http.StatusCreated,
		video,
	)
}
