package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestJSONResponse tests the JSONResponse helper function
func TestJSONResponse(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Arrange
		w := httptest.NewRecorder()
		code := http.StatusOK
		response := map[string]string{"key": "value"}

		// Act
		JSONResponse(w, code, response)

		// Assert
		assert.Equal(t, code, w.Code)
		assert.Equal(t, "application/json; charset=UTF-8", w.Header().Get("Content-Type"))
		assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))

		// Verify JSON content
		var result map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &result)
		assert.NoError(t, err)
		assert.Equal(t, response, result)
	})

	t.Run("Error with marshal", func(t *testing.T) {
		// Arrange
		w := httptest.NewRecorder()
		code := http.StatusOK
		// Create a value that can't be marshaled to JSON
		response := make(chan int) // Channels can't be marshaled to JSON

		// Act
		JSONResponse(w, code, response)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Internal server error")
	})
}
