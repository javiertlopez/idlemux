package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	awesome "github.com/javiertlopez/awesome/pkg"
)

// addVideoHandler adds the handler to the mux router
func (a *App) addVideoHandler(r *mux.Router) {
	r.HandleFunc("/videos", a.CreateVideoHandler).Methods("POST")
	r.HandleFunc("/videos/{id}", a.ReadVideoHandler).Methods("GET")
}

// CreateVideoHandler handler
func (a *App) CreateVideoHandler(w http.ResponseWriter, r *http.Request) {
	var video awesome.Video
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

	// If body contains a Source File URL, send it to Ingestion
	if len(video.SourceURL) > 0 {
		assetID, err := a.assets.Ingest(r.Context(), video.SourceURL, true)
		if err == nil {
			video.Asset = &awesome.Asset{
				ID: assetID,
			}
		}
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
		if err == awesome.ErrVideoNotFound {
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

	// If video document contains an Asset ID, retrieve the information
	if response.Asset != nil {
		asset, err := a.assets.GetByID(r.Context(), response.Asset.ID)
		if err == nil {
			response.Asset = asset
		}
	}

	JSONResponse(
		w,
		http.StatusOK,
		response,
	)
}
