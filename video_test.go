package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
)

type MockedVideos struct {
	mock.Mock
}

func (m *MockedVideos) Insert(ctx context.Context, anyVideo *Video) (*Video, error) {
	uuid := "fcdf5f4e-b086-4b52-8714-bf3623186185"

	anyVideo.ID = &uuid
	return anyVideo, nil
}

func TestCreateVideoHandler(t *testing.T) {
	logger := logrus.New()
	logger.Out = ioutil.Discard

	completeVideo := &Video{
		Title:       "Some Might Say",
		Description: "Oasis song from (What's the Story) Morning Glory? album.",
	}

	titleVideo := &Video{
		Title: "Some Might Say",
	}

	descriptionVideo := &Video{
		Description: "Oasis song from (What's the Story) Morning Glory? album.",
	}

	emptyVideo := &Video{}

	mocked := new(MockedVideos)

	// Create an app
	awesome := App{
		logger: logger,
		videos: mocked,
	}

	expectedComplete := `{"id":"fcdf5f4e-b086-4b52-8714-bf3623186185","title":"Some Might Say","description":"Oasis song from (What's the Story) Morning Glory? album."}`
	expectedUnprocessable := `{"message":"Unprocessable Entity","status":422}`

	tests := []struct {
		name         string
		expectedCode int
		expectedBody string
		body         *Video
	}{
		{"Valid", http.StatusCreated, expectedComplete, completeVideo},
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
}
