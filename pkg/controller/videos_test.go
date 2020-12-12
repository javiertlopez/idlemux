package controller

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/javiertlopez/awesome/pkg/errorcodes"
	mocks "github.com/javiertlopez/awesome/pkg/mocks/usecase"
	"github.com/javiertlopez/awesome/pkg/model"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
)

func Test_videoController_Create(t *testing.T) {
	completeVideo := model.Video{
		Title:       "Some Might Say",
		Description: "(What's the Story) Morning Glory?",
		SourceURL:   "https://storage.googleapis.com/muxdemofiles/mux-video-intro.mp4",
	}

	tests := []struct {
		name         string
		expectedCode int
		expectedBody string
		body         string
		video        model.Video
		wantedError  error
	}{
		{
			"Complete",
			201,
			`{"title":"Some Might Say","description":"(What's the Story) Morning Glory?","source_url":"https://storage.googleapis.com/muxdemofiles/mux-video-intro.mp4"}`,
			`{"title":"Some Might Say","description":"(What's the Story) Morning Glory?","source_url":"https://storage.googleapis.com/muxdemofiles/mux-video-intro.mp4"}`,
			completeVideo,
			nil,
		},
		{
			"Video Unprocessable (title)",
			422,
			`{"message":"Unprocessable entity","status":422}`,
			`{"description":"(What's the Story) Morning Glory?"}`,
			completeVideo,
			errorcodes.ErrVideoUnprocessable,
		},
		{
			"Error",
			500,
			`{"message":"Internal server error","status":500}`,
			`{"description":"(What's the Story) Morning Glory?"}`,
			model.Video{},
			errors.New("failed"),
		},
		{
			"Bad request",
			400,
			`{"message":"Bad request","status":400}`,
			`{"title":23,"description":?",}`,
			model.Video{},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			videos := &mocks.Videos{}
			vc := &videoController{
				videos,
			}

			r, _ := http.NewRequest("POST", "/videos", bytes.NewBuffer([]byte(tt.body)))
			w := httptest.NewRecorder()

			ctx := r.Context()

			videos.On("Create", ctx, mock.Anything).Return(tt.video, tt.wantedError)

			vc.Create(w, r)

			// Check the content type is what we expect.
			expected := "application/json; charset=UTF-8"
			m := w.Header()
			if contentType := m.Get("Content-Type"); contentType != expected {
				t.Errorf(
					"handler returned wrong content type: got %v want %v",
					contentType,
					expected,
				)
			}

			// Check the status code is what we expect.
			if status := w.Code; status != tt.expectedCode {
				t.Errorf(
					"handler returned wrong status code: got %v want %v",
					status,
					tt.expectedCode,
				)
			}

			// Check the response body is what we expect.
			if w.Body.String() != tt.expectedBody {
				t.Errorf(
					"handler returned unexpected body: got %v want %v",
					w.Body.String(),
					tt.expectedBody,
				)
			}
		})
	}
}

func Test_videoController_GetByID(t *testing.T) {
	uuid := "4e5bf8f2-9c50-4576-b9d4-1d1fd0705885"
	completeVideo := model.Video{
		ID:          &uuid,
		Title:       "Some Might Say",
		Description: "Oasis song from (What's the Story) Morning Glory? album.",
	}

	tests := []struct {
		name         string
		id           string
		expectedCode int
		expectedBody string
		video        model.Video
		wantedError  error
	}{
		{
			"Success",
			"4e5bf8f2-9c50-4576-b9d4-1d1fd0705885",
			200,
			`{"id":"4e5bf8f2-9c50-4576-b9d4-1d1fd0705885","title":"Some Might Say","description":"Oasis song from (What's the Story) Morning Glory? album."}`,
			completeVideo,
			nil,
		},
		{
			"Bad ID",
			"123",
			422,
			`{"message":"Unprocessable Entity","status":422}`,
			model.Video{},
			nil,
		},
		{
			"Not found",
			"4e5bf8f2-9c50-4576-b9d4-1d1fd0705885",
			404,
			`{"message":"Not found","status":404}`,
			model.Video{},
			errorcodes.ErrVideoNotFound,
		},
		{
			"Error",
			"4e5bf8f2-9c50-4576-b9d4-1d1fd0705885",
			500,
			`{"message":"Internal server error","status":500}`,
			model.Video{},
			errors.New("failed"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			videos := &mocks.Videos{}
			vc := &videoController{
				videos,
			}

			r, _ := http.NewRequest("GET", "/videos/abcd", nil)
			w := httptest.NewRecorder()

			r = mux.SetURLVars(r, map[string]string{
				"id": tt.id,
			})

			ctx := r.Context()

			videos.On("GetByID", ctx, tt.id).Return(tt.video, tt.wantedError)

			vc.GetByID(w, r)
			// Check the content type is what we expect.
			expected := "application/json; charset=UTF-8"
			m := w.Header()
			if contentType := m.Get("Content-Type"); contentType != expected {
				t.Errorf(
					"handler returned wrong content type: got %v want %v",
					contentType,
					expected,
				)
			}

			// Check the status code is what we expect.
			if status := w.Code; status != tt.expectedCode {
				t.Errorf(
					"handler returned wrong status code: got %v want %v",
					status,
					tt.expectedCode,
				)
			}

			// Check the response body is what we expect.
			if w.Body.String() != tt.expectedBody {
				t.Errorf(
					"handler returned unexpected body: got %v want %v",
					w.Body.String(),
					tt.expectedBody,
				)
			}
		})
	}
}
