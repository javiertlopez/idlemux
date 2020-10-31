package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
)

type MockedVideos struct {
	mock.Mock
}

func (m *MockedVideos) Insert(ctx context.Context, anyVideo *Video) (*Video, error) {
	uuid := "fcdf5f4e-b086-4b52-8714-bf3623186185"

	response := &Video{
		ID:          &uuid,
		Title:       anyVideo.Title,
		Description: anyVideo.Description,
	}

	if anyVideo.Asset != nil {
		response.Asset = &Asset{
			ID: anyVideo.Asset.ID,
		}
	}

	return response, nil
}

func (m *MockedVideos) GetByID(ctx context.Context, id string) (*Video, error) {
	ID := "fcdf5f4e-b086-4b52-8714-bf3623186185"
	IDWithSourceFile := "a9200233-9b62-489c-9cbc-bb37f2922804"

	switch id {
	case ID:
		return &Video{
			ID:          &ID,
			Title:       "Some Might Say",
			Description: "Oasis song from (What's the Story) Morning Glory? album.",
		}, nil
	case IDWithSourceFile:
		return &Video{
			ID:          &ID,
			Title:       "Some Might Say",
			Description: "Oasis song from (What's the Story) Morning Glory? album.",
			Asset: &Asset{
				ID: "5iNFJg9dIww2AgUryhgghbP00Dc4ogoxn00gzitOdjICg",
			},
		}, nil
	}

	return nil, ErrVideoNotFound
}

func TestCreateVideoHandler(t *testing.T) {
	logger := logrus.New()
	logger.Out = ioutil.Discard

	completeVideo := &Video{
		Title:       "Some Might Say",
		Description: "Oasis song from (What's the Story) Morning Glory? album.",
		SourceURL:   "https://storage.googleapis.com/muxdemofiles/mux-video-intro.mp4",
	}

	incompleteVideo := &Video{
		Title:       "Some Might Say",
		Description: "Oasis song from (What's the Story) Morning Glory? album.",
	}

	invalidVideo := &Video{
		Title:       "Some Might Say",
		Description: "Oasis song from (What's the Story) Morning Glory? album.",
		SourceURL:   "invalidURL",
	}

	titleVideo := &Video{
		Title: "Some Might Say",
	}

	descriptionVideo := &Video{
		Description: "Oasis song from (What's the Story) Morning Glory? album.",
	}

	emptyVideo := &Video{}

	mockedVideos := new(MockedVideos)
	mockedAssets := new(MockedAssets)

	// Create an app
	awesome := App{
		logger: logger,
		videos: mockedVideos,
		assets: mockedAssets,
	}

	expectedComplete := `{"id":"fcdf5f4e-b086-4b52-8714-bf3623186185","title":"Some Might Say","description":"Oasis song from (What's the Story) Morning Glory? album.","asset":{"id":"5iNFJg9dIww2AgUryhgghbP00Dc4ogoxn00gzitOdjICg"}}`
	expectedIncomplete := `{"id":"fcdf5f4e-b086-4b52-8714-bf3623186185","title":"Some Might Say","description":"Oasis song from (What's the Story) Morning Glory? album."}`
	expectedUnprocessable := `{"message":"Unprocessable Entity","status":422}`
	expectedBadRequest := `{"message":"Invalid request","status":400}`

	tests := []struct {
		name         string
		expectedCode int
		expectedBody string
		body         *Video
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
			handler := http.HandlerFunc(awesome.CreateVideoHandler)

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
		handler := http.HandlerFunc(awesome.CreateVideoHandler)

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

func TestReadVideoHandler(t *testing.T) {
	logger := logrus.New()
	logger.Out = ioutil.Discard

	ID := "fcdf5f4e-b086-4b52-8714-bf3623186185"
	IDWithSourceFile := "a9200233-9b62-489c-9cbc-bb37f2922804"

	mockedVideos := new(MockedVideos)
	mockedAssets := new(MockedAssets)

	// Create an app
	awesome := App{
		logger: logger,
		videos: mockedVideos,
		assets: mockedAssets,
	}

	expectedComplete := `{"id":"fcdf5f4e-b086-4b52-8714-bf3623186185","title":"Some Might Say","description":"Oasis song from (What's the Story) Morning Glory? album."}`
	expectedWithSourceFile := `{"id":"fcdf5f4e-b086-4b52-8714-bf3623186185","title":"Some Might Say","description":"Oasis song from (What's the Story) Morning Glory? album.","asset":{"id":"5iNFJg9dIww2AgUryhgghbP00Dc4ogoxn00gzitOdjICg"}}`
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
			router.HandleFunc("/videos/{id}", awesome.ReadVideoHandler)

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
