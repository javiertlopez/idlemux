package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/javiertlopez/awesome/pkg/model"
	"github.com/javiertlopez/awesome/pkg/repository/axiom"
	"github.com/javiertlopez/awesome/pkg/repository/muxinc"
	"github.com/javiertlopez/awesome/pkg/usecase"
	"github.com/sirupsen/logrus"
)

func TestCreate(t *testing.T) {
	logger := logrus.New()
	logger.Out = ioutil.Discard

	completeVideo := &model.Video{
		Title:       "Some Might Say",
		Description: "Oasis song from (What's the Story) Morning Glory? album.",
		SourceURL:   "https://storage.googleapis.com/muxdemofiles/mux-video-intro.mp4",
	}

	incompleteVideo := &model.Video{
		Title:       "Some Might Say",
		Description: "Oasis song from (What's the Story) Morning Glory? album.",
	}

	invalidVideo := &model.Video{
		Title:       "Some Might Say",
		Description: "Oasis song from (What's the Story) Morning Glory? album.",
		SourceURL:   "invalidURL",
	}

	titleVideo := &model.Video{
		Title: "Some Might Say",
	}

	descriptionVideo := &model.Video{
		Description: "Oasis song from (What's the Story) Morning Glory? album.",
	}

	emptyVideo := &model.Video{}

	mockedAssets := new(muxinc.MockedAssets)
	mockedVideos := new(axiom.MockedVideos)

	videos := usecase.NewVideoUseCase(mockedAssets, mockedVideos)

	controller := NewVideoController(videos)

	expectedComplete := `{"id":"fcdf5f4e-b086-4b52-8714-bf3623186185","title":"Some Might Say","description":"Oasis song from (What's the Story) Morning Glory? album.","asset":{"id":"5iNFJg9dIww2AgUryhgghbP00Dc4ogoxn00gzitOdjICg"}}`
	expectedIncomplete := `{"id":"fcdf5f4e-b086-4b52-8714-bf3623186185","title":"Some Might Say","description":"Oasis song from (What's the Story) Morning Glory? album."}`
	expectedUnprocessable := `{"message":"Unprocessable Entity","status":422}`
	expectedBadRequest := `{"message":"Invalid request","status":400}`

	tests := []struct {
		name         string
		expectedCode int
		expectedBody string
		body         *model.Video
	}{
		{"Valid", http.StatusCreated, expectedComplete, completeVideo},
		{"Valid (with source file)", http.StatusCreated, expectedIncomplete, incompleteVideo},
		{"Invalid Source File URL", http.StatusCreated, expectedIncomplete, invalidVideo},
		{"Empty description", http.StatusUnprocessableEntity, expectedUnprocessable, titleVideo},
		{"Empty title", http.StatusUnprocessableEntity, expectedUnprocessable, descriptionVideo},
		{"Empty body", http.StatusUnprocessableEntity, expectedUnprocessable, emptyVideo},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, _ := json.Marshal(tt.body)

			// Create a request to pass to our handler.
			req, err := http.NewRequest("POST", "/videos", bytes.NewBuffer(output))
			if err != nil {
				t.Fatal(err)
			}

			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(controller.Create)

			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			handler.ServeHTTP(rr, req)

			// Check the content type is what we expect.
			expected := "application/json; charset=UTF-8"
			m := rr.Header()
			if contentType := m.Get("Content-Type"); contentType != expected {
				t.Errorf(
					"handler returned wrong content type: got %v want %v",
					contentType,
					expected,
				)
			}

			// Check the status code is what we expect.
			if status := rr.Code; status != tt.expectedCode {
				t.Errorf(
					"handler returned wrong status code: got %v want %v",
					status,
					tt.expectedCode,
				)
			}

			// Check the response body is what we expect.
			if rr.Body.String() != tt.expectedBody {
				t.Errorf(
					"handler returned unexpected body: got %v want %v",
					rr.Body.String(),
					tt.expectedBody,
				)
			}
		})
	}

	t.Run("Invalid JSON", func(t *testing.T) {
		invalidJSON := `{"title:Unprocessable Entity","description":422}`

		output, _ := json.Marshal(invalidJSON)

		// Create a request to pass to our handler.
		req, err := http.NewRequest("POST", "/videos", bytes.NewBuffer(output))
		if err != nil {
			t.Fatal(err)
		}

		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(controller.Create)

		// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
		// directly and pass in our Request and ResponseRecorder.
		handler.ServeHTTP(rr, req)

		// Check the content type is what we expect.
		expected := "application/json; charset=UTF-8"
		m := rr.Header()
		if contentType := m.Get("Content-Type"); contentType != expected {
			t.Errorf(
				"handler returned wrong content type: got %v want %v",
				contentType,
				expected,
			)
		}

		// Check the status code is what we expect.
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf(
				"handler returned wrong status code: got %v want %v",
				status,
				http.StatusBadRequest,
			)
		}

		// Check the response body is what we expect.
		if rr.Body.String() != expectedBadRequest {
			t.Errorf(
				"handler returned unexpected body: got %v want %v",
				rr.Body.String(),
				expectedBadRequest,
			)
		}
	})
}

func TestGetByID(t *testing.T) {
	logger := logrus.New()
	logger.Out = ioutil.Discard

	ID := "fcdf5f4e-b086-4b52-8714-bf3623186185"
	IDWithSourceFile := "a9200233-9b62-489c-9cbc-bb37f2922804"

	mockedAssets := new(muxinc.MockedAssets)
	mockedVideos := new(axiom.MockedVideos)

	videos := usecase.NewVideoUseCase(mockedAssets, mockedVideos)

	controller := NewVideoController(videos)

	expectedComplete := `{"id":"fcdf5f4e-b086-4b52-8714-bf3623186185","title":"Some Might Say","description":"Oasis song from (What's the Story) Morning Glory? album."}`
	expectedWithSourceFile := `{"id":"a9200233-9b62-489c-9cbc-bb37f2922804","title":"Some Might Say","description":"Oasis song from (What's the Story) Morning Glory? album.","poster":"https://image.mux.com/5iNFJg9dIww2AgUryhgghbP00Dc4ogoxn00gzitOdjICg/thumbnail.png?width=1920\u0026height=1080\u0026smart_crop=true\u0026time=7","thumbnail":"https://image.mux.com/5iNFJg9dIww2AgUryhgghbP00Dc4ogoxn00gzitOdjICg/thumbnail.png?width=640\u0026height=360\u0026smart_crop=true\u0026time=7","sources":[{"src":"https://stream.mux.com/5iNFJg9dIww2AgUryhgghbP00Dc4ogoxn00gzitOdjICg.m3u8","type":"application/x-mpegURL"}]}`
	expectedNotFound := `{"message":"video not found","status":404}`
	expectedUnprocessable := `{"message":"Unprocessable Entity","status":422}`

	tests := []struct {
		name         string
		expectedCode int
		expectedBody string
		ID           string
	}{
		{"Valid", http.StatusOK, expectedComplete, ID},
		{"Valid (with source file)", http.StatusOK, expectedWithSourceFile, IDWithSourceFile},
		{"Not found", http.StatusNotFound, expectedNotFound, "xxxxxxxx-xxxx-Mxxx-Nxxx-xxxxxxxxxxxx"},
		{"Not UUID", http.StatusUnprocessableEntity, expectedUnprocessable, "xxxxxxxx"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a request to pass to our handler.
			req, err := http.NewRequest("GET", fmt.Sprintf("/videos/%s", tt.ID), nil)
			if err != nil {
				t.Fatal(err)
			}

			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rr := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/videos/{id}", controller.GetByID)

			// Change to Gorilla Mux router to pass variables
			router.ServeHTTP(rr, req)

			// Check the content type is what we expect.
			expected := "application/json; charset=UTF-8"
			m := rr.Header()
			if contentType := m.Get("Content-Type"); contentType != expected {
				t.Errorf(
					"handler returned wrong content type: got %v want %v",
					contentType,
					expected,
				)
			}

			// Check the status code is what we expect.
			if status := rr.Code; status != tt.expectedCode {
				t.Errorf(
					"handler returned wrong status code: got %v want %v",
					status,
					tt.expectedCode,
				)
			}

			// Check the response body is what we expect.
			if rr.Body.String() != tt.expectedBody {
				t.Errorf(
					"handler returned unexpected body: got %v want %v",
					rr.Body.String(),
					tt.expectedBody,
				)
			}
		})
	}
}
