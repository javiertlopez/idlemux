package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// Videos interface, for testing purposes
type Videos interface {
	Insert(ctx context.Context, anyVideo *Video) (*Video, error)
	GetByID(ctx context.Context, id string) (*Video, error)
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
func (a *App) addVideoHandler(r *mux.Router) {
	r.HandleFunc("/videos", a.CreateVideoHandler).Methods("POST")
	r.HandleFunc("/videos/{id}", a.ReadVideoHandler).Methods("GET")
}

// CreateVideoHandler handler
func (a *App) CreateVideoHandler(w http.ResponseWriter, r *http.Request) {
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

	// Title and Description are mandatory fields
	if len(video.Title) == 0 || len(video.Description) == 0 {
		JSONResponse(
			w, http.StatusUnprocessableEntity,
			Response{
				Message: "Unprocessable Entity",
				Status:  http.StatusUnprocessableEntity,
			},
		)
		return
	}

	response, err := a.videos.Insert(r.Context(), &video)

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
		response,
	)
}

// ReadVideoHandler handler
func (a *App) ReadVideoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Wrong type of ID should return 422 error?
	if len(id) != 36 {
		JSONResponse(
			w, http.StatusUnprocessableEntity,
			Response{
				Message: "Unprocessable Entity",
				Status:  http.StatusUnprocessableEntity,
			},
		)
		return
	}

	response, err := a.videos.GetByID(r.Context(), id)

	if err != nil {
		// Look for Custom Error
		if err == ErrVideoNotFound {
			JSONResponse(
				w, http.StatusNotFound,
				Response{
					Message: err.Error(),
					Status:  http.StatusNotFound,
				},
			)
			return
		}

		// Anything besides Not Found should be return as an internal error
		JSONResponse(
			w, http.StatusInternalServerError,
			Response{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			},
		)
		return
	}

	JSONResponse(
		w,
		http.StatusOK,
		response,
	)
}
