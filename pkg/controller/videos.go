package controller

import (
	"encoding/json"
	"net/http"

	"github.com/javiertlopez/awesome/pkg/errorcodes"
	"github.com/javiertlopez/awesome/pkg/model"

	"github.com/gorilla/mux"
)

// Create controller
func (vc *videoController) Create(w http.ResponseWriter, r *http.Request) {
	var video model.Video
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

	response, err := vc.videos.Create(r.Context(), video)

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

// GetByID controller
func (vc *videoController) GetByID(w http.ResponseWriter, r *http.Request) {
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

	response, err := vc.videos.GetByID(r.Context(), id)

	if err != nil {
		// Look for Custom Error
		if err == errorcodes.ErrVideoNotFound {
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
