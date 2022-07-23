package controller

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
)

const (
	commit  = "2a4ca47"
	version = "1.2.3"
)

func TestHealthz(t *testing.T) {
	logger := logrus.New()
	logger.Out = ioutil.Discard

	// Create an app
	controller := controller{
		commit:  commit,
		version: version,
	}
	// Create a request to pass to our handler.
	req, err := http.NewRequest("GET", "/app/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controller.Healthz)

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
	if status := rr.Code; status != http.StatusOK {
		t.Errorf(
			"handler returned wrong status code: got %v want %v",
			status,
			http.StatusOK,
		)
	}

	// Check the response body is what we expect.
	expected = `{"message":"Hello World!","status":200}`
	if rr.Body.String() != expected {
		t.Errorf(
			"handler returned unexpected body: got %v want %v",
			rr.Body.String(),
			expected,
		)
	}
}

func TestStatusz(t *testing.T) {
	logger := logrus.New()
	logger.Out = ioutil.Discard

	// Create an app
	controller := controller{
		commit:  commit,
		version: version,
	}

	// Create a request to pass to our handler.
	req, err := http.NewRequest("GET", "/app/statusz", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controller.Statusz)

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
	if status := rr.Code; status != http.StatusOK {
		t.Errorf(
			"handler returned wrong status code: got %v want %v",
			status,
			http.StatusOK,
		)
	}

	// Check the response body is what we expect.
	expected = `{"commit":"2a4ca47","version":"1.2.3"}`
	if rr.Body.String() != expected {
		t.Errorf(
			"handler returned unexpected body: got %v want %v",
			rr.Body.String(),
			expected,
		)
	}
}
