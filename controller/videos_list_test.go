package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/javiertlopez/awesome/controller/mocks"
	"github.com/javiertlopez/awesome/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_videoController_List(t *testing.T) {
	videos := []model.Video{
		{ID: "id1", Title: "Video 1"},
		{ID: "id2", Title: "Video 2"},
	}
	t.Run("Success", func(t *testing.T) {
		delivery := &mocks.Delivery{}
		delivery.On("List", mock.Anything, 1, 10).Return(videos, nil)
		controller := &controller{delivery: delivery}

		r, _ := http.NewRequest("GET", "/videos", nil)
		w := httptest.NewRecorder()

		controller.List(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		var got []model.Video
		err := json.Unmarshal(w.Body.Bytes(), &got)
		assert.NoError(t, err)
		assert.Equal(t, videos, got)
	})
	t.Run("With pagination params", func(t *testing.T) {
		delivery := &mocks.Delivery{}
		delivery.On("List", mock.Anything, 2, 5).Return(videos, nil)
		controller := &controller{delivery: delivery}

		r, _ := http.NewRequest("GET", "/videos?page=2&limit=5", nil)
		w := httptest.NewRecorder()

		controller.List(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		var got []model.Video
		err := json.Unmarshal(w.Body.Bytes(), &got)
		assert.NoError(t, err)
		assert.Equal(t, videos, got)
	})
	t.Run("Internal error", func(t *testing.T) {
		delivery := &mocks.Delivery{}
		delivery.On("List", mock.Anything, 1, 10).Return(nil, assert.AnError)
		controller := &controller{delivery: delivery}

		r, _ := http.NewRequest("GET", "/videos", nil)
		w := httptest.NewRecorder()

		controller.List(w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "Internal server error", resp["message"])
		assert.Equal(t, float64(http.StatusInternalServerError), resp["status"])
	})
}
