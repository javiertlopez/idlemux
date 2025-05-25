package controller

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

const (
	commit  = "2a4ca47"
	version = "1.2.3"
)

func TestHealthz(t *testing.T) {
	logger := logrus.New()
	logger.Out = io.Discard

	// Create an app
	controller := controller{
		commit:  commit,
		version: version,
	}
	// Create a request to pass to our handler.
	req, err := http.NewRequest("GET", "/app/health", nil)
	assert.NoError(t, err)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controller.Healthz)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the content type, status code and body
	assert.Equal(t, "application/json; charset=UTF-8", rr.Header().Get("Content-Type"), "Should return JSON content type")
	assert.Equal(t, http.StatusOK, rr.Code, "Should return OK status code")
	assert.Equal(t, `{"message":"Hello World!","status":200}`, rr.Body.String(), "Response body should match expected")
}

func TestStatusz(t *testing.T) {
	logger := logrus.New()
	logger.Out = io.Discard

	// Create an app
	controller := controller{
		commit:  commit,
		version: version,
	}

	// Create a request to pass to our handler.
	req, err := http.NewRequest("GET", "/app/statusz", nil)
	assert.NoError(t, err)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controller.Statusz)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the content type, status code and body
	assert.Equal(t, "application/json; charset=UTF-8", rr.Header().Get("Content-Type"), "Should return JSON content type")
	assert.Equal(t, http.StatusOK, rr.Code, "Should return OK status code")
	assert.Equal(t, `{"commit":"2a4ca47","version":"1.2.3"}`, rr.Body.String(), "Response body should match expected")
}
