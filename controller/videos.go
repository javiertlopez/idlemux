package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/javiertlopez/awesome/errorcodes"
	"github.com/javiertlopez/awesome/model"
)

// Create controller
func (c controller) Create(w http.ResponseWriter, r *http.Request) {
	var video model.Video
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&video); err != nil {
		JSONResponse(
			w, http.StatusBadRequest,
			Response{
				Message: "Bad request",
				Status:  http.StatusBadRequest,
			},
		)
		return
	}
	defer r.Body.Close()

	response, err := c.ingestion.Create(r.Context(), video)

	if err != nil {
		// Look for Custom Error
		if err == errorcodes.ErrVideoUnprocessable {
			JSONResponse(
				w, http.StatusUnprocessableEntity,
				Response{
					Message: "Unprocessable entity",
					Status:  http.StatusUnprocessableEntity,
				},
			)
			return
		}

		JSONResponse(
			w, http.StatusInternalServerError,
			Response{
				Message: "Internal server error",
				Status:  http.StatusInternalServerError,
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

// GetByID controller
func (c controller) GetByID(w http.ResponseWriter, r *http.Request) {
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

	response, err := c.delivery.GetByID(r.Context(), id)

	if err != nil {
		// Look for Custom Error
		if err == errorcodes.ErrVideoNotFound {
			JSONResponse(
				w, http.StatusNotFound,
				Response{
					Message: "Not found",
					Status:  http.StatusNotFound,
				},
			)
			return
		}

		// Anything besides Not Found should be return as an internal error
		JSONResponse(
			w, http.StatusInternalServerError,
			Response{
				Message: "Internal server error",
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

// List controller
func (c controller) List(w http.ResponseWriter, r *http.Request) {
	page := 1
	limit := 10
	if p := r.URL.Query().Get("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}
	if l := r.URL.Query().Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = v
		}
	}
	videos, err := c.delivery.List(r.Context(), page, limit)
	if err != nil {
		JSONResponse(
			w, http.StatusInternalServerError,
			Response{
				Message: "Internal server error",
				Status:  http.StatusInternalServerError,
			},
		)
		return
	}
	JSONResponse(w, http.StatusOK, videos)
}
