package controller

import (
	"net/http"

	"github.com/javiertlopez/awesome/pkg/usecase"
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

// NewEventController returns an EventController
func NewEventController(videos usecase.Videos) VideoController {
	return &videoController{
		videos,
	}
}
