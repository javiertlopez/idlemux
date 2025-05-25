package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/javiertlopez/awesome/errorcodes"
	"github.com/javiertlopez/awesome/model"
)

func Test_videoController_Create(t *testing.T) {
	completeVideo := model.Video{
		Title:       "Some Might Say",
		Description: "(What's the Story) Morning Glory?",
		SourceURL:   "https://storage.googleapis.com/muxdemofiles/mux-video-intro.mp4",
	}

	tests := []struct {
		name         string
		body         string
		video        model.Video
		wantedError  error
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Complete",
			body:         `{"title":"Some Might Say","description":"(What's the Story) Morning Glory?","source_url":"https://storage.googleapis.com/muxdemofiles/mux-video-intro.mp4"}`,
			video:        completeVideo,
			wantedError:  nil,
			expectedCode: http.StatusCreated,
			expectedBody: `{"title":"Some Might Say","description":"(What's the Story) Morning Glory?","source_url":"https://storage.googleapis.com/muxdemofiles/mux-video-intro.mp4"}`,
		},
		{
			name:         "Video Unprocessable (title)",
			body:         `{"description":"(What's the Story) Morning Glory?"}`,
			video:        completeVideo,
			wantedError:  errorcodes.ErrVideoUnprocessable,
			expectedCode: http.StatusUnprocessableEntity,
			expectedBody: `{"message":"Unprocessable entity","status":422}`,
		},
		{
			name:         "Error",
			body:         `{"description":"(What's the Story) Morning Glory?"}`,
			video:        model.Video{},
			wantedError:  errors.New("failed"),
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"message":"Internal server error","status":500}`,
		},
		{
			name:         "Bad request",
			body:         `{"title":23,"description":?",}`,
			video:        model.Video{},
			wantedError:  nil,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"message":"Bad request","status":400}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ingestion := NewMockIngestion(t)
			controller := &controller{
				ingestion: ingestion,
			}

			r, _ := http.NewRequest("POST", "/videos", bytes.NewBuffer([]byte(tt.body)))
			w := httptest.NewRecorder()

			ctx := r.Context()

			// Only set up mock expectations for valid JSON
			if tt.name != "Bad request" {
				ingestion.On("Create", ctx, mock.Anything).Return(tt.video, tt.wantedError)
			}

			controller.Create(w, r)

			// Check the content type, status code and body
			assert.Equal(t, "application/json; charset=UTF-8", w.Header().Get("Content-Type"), "Should return JSON content type")
			assert.Equal(t, tt.expectedCode, w.Code, "Should return expected status code")
			assert.Equal(t, tt.expectedBody, w.Body.String(), "Response body should match expected")
		})
	}
}

func Test_videoController_GetByID(t *testing.T) {
	uuid := "4e5bf8f2-9c50-4576-b9d4-1d1fd0705885"
	completeVideo := model.Video{
		ID:          uuid,
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
			delivery := NewMockDelivery(t)
			controller := &controller{
				delivery: delivery,
			}

			r, _ := http.NewRequest("GET", "/videos/abcd", nil)
			w := httptest.NewRecorder()

			r = mux.SetURLVars(r, map[string]string{
				"id": tt.id,
			})

			ctx := r.Context()

			// Only set up mock expectations for valid IDs (36 chars)
			if len(tt.id) == 36 {
				delivery.On("GetByID", ctx, tt.id).Return(tt.video, tt.wantedError)
			}

			controller.GetByID(w, r)

			// Check the content type, status code and body
			assert.Equal(t, "application/json; charset=UTF-8", w.Header().Get("Content-Type"), "Should return JSON content type")
			assert.Equal(t, tt.expectedCode, w.Code, "Should return expected status code")
			assert.Equal(t, tt.expectedBody, w.Body.String(), "Response body should match expected")
		})
	}
}

func Test_videoController_List(t *testing.T) {
	videos := []model.Video{
		{ID: "id1", Title: "Video 1"},
		{ID: "id2", Title: "Video 2"},
	}

	videosJSON, _ := json.Marshal(videos)
	errorJSON := `{"message":"Internal server error","status":500}`

	tests := []struct {
		name         string
		url          string
		page         int
		limit        int
		mockReturn   []model.Video
		mockError    error
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Success",
			url:          "/videos",
			page:         1,
			limit:        10,
			mockReturn:   videos,
			mockError:    nil,
			expectedCode: http.StatusOK,
			expectedBody: string(videosJSON),
		},
		{
			name:         "With pagination params",
			url:          "/videos?page=2&limit=5",
			page:         2,
			limit:        5,
			mockReturn:   videos,
			mockError:    nil,
			expectedCode: http.StatusOK,
			expectedBody: string(videosJSON),
		},
		{
			name:         "Internal error",
			url:          "/videos",
			page:         1,
			limit:        10,
			mockReturn:   nil,
			mockError:    assert.AnError,
			expectedCode: http.StatusInternalServerError,
			expectedBody: errorJSON,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			delivery := NewMockDelivery(t)
			delivery.On("List", mock.Anything, tt.page, tt.limit).Return(tt.mockReturn, tt.mockError)
			controller := &controller{delivery: delivery}

			r, _ := http.NewRequest("GET", tt.url, nil)
			w := httptest.NewRecorder()

			controller.List(w, r)

			assert.Equal(t, tt.expectedCode, w.Code, "Should return expected status code")
			if tt.mockError == nil {
				var got []model.Video
				err := json.Unmarshal(w.Body.Bytes(), &got)
				assert.NoError(t, err, "Should unmarshal response body without errors")
				assert.Equal(t, tt.mockReturn, got, "Response should match expected videos")
			} else {
				var resp map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err, "Should unmarshal error response without errors")
				assert.Equal(t, "Internal server error", resp["message"], "Error message should match")
				assert.Equal(t, float64(http.StatusInternalServerError), resp["status"], "Status code should match")
			}
		})
	}
}
