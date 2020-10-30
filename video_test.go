package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateVideoHandler(t *testing.T) {
	mockedVideo := &Video{
		Title:       "Some Might Say",
		Description: "Oasis song from (What's the Story) Morning Glory? album.",
	}

	output, _ := json.Marshal(mockedVideo)

	// Create a request to pass to our handler.
	req, err := http.NewRequest("POST", "/videos", bytes.NewBuffer(output))
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateVideoHandler)

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
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf(
			"handler returned wrong status code: got %v want %v",
			status,
			http.StatusCreated,
		)
	}

	// Check the response body is what we expect.
	expected = `{"id":"fcdf5f4e-b086-4b52-8714-bf3623186185","title":"Some Might Say","description":"Oasis song from (What's the Story) Morning Glory? album."}`
	if rr.Body.String() != expected {
		t.Errorf(
			"handler returned unexpected body: got %v want %v",
			rr.Body.String(),
			expected,
		)
	}
}
