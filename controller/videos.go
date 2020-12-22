package controller

import (
	"encoding/json"
	"net/http"

	"github.com/javiertlopez/awesome/errorcodes"
	"github.com/javiertlopez/awesome/model"
	"github.com/javiertlopez/awesome/usecase"

	"github.com/gorilla/mux"
)

// VideoController handles the HTTP requests
type VideoController interface {
	Create(w http.ResponseWriter, r *http.Request)
	GetByID(w http.ResponseWriter, r *http.Request)
}

// videoController struct holds the usecase
type videoController struct {
	videos usecase.Videos
}

// NewVideoController returns a VideoController
func NewVideoController(videos usecase.Videos) VideoController {
	return &videoController{
		videos: videos,
	}
}

// Create controller
func (vc *videoController) Create(w http.ResponseWriter, r *http.Request) {
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

	response, err := vc.videos.Create(r.Context(), video)

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
