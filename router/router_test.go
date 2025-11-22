package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNew(t *testing.T) {
	mockController := NewMockController(t)

	router := New(mockController)

	assert.NotNil(t, router)
	assert.IsType(t, &mux.Router{}, router)
}

func TestRouter_Routes(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		path         string
		expectedCode int
	}{
		{
			name:         "Healthz endpoint",
			method:       "GET",
			path:         "/app/healthz",
			expectedCode: http.StatusOK,
		},
		{
			name:         "Statusz endpoint",
			method:       "GET",
			path:         "/app/statusz",
			expectedCode: http.StatusOK,
		},
		{
			name:         "Create endpoint",
			method:       "POST",
			path:         "/videos",
			expectedCode: http.StatusCreated,
		},
		{
			name:         "Get by ID endpoint",
			method:       "GET",
			path:         "/videos/123",
			expectedCode: http.StatusOK,
		},
		{
			name:         "List endpoint",
			method:       "GET",
			path:         "/videos",
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockController := NewMockController(t)
			mockController.On("Healthz", mock.Anything, mock.Anything).Maybe().Run(func(args mock.Arguments) {
				w := args.Get(0).(http.ResponseWriter)
				w.WriteHeader(http.StatusOK)
			}).Return()
			mockController.On("Statusz", mock.Anything, mock.Anything).Maybe().Run(func(args mock.Arguments) {
				w := args.Get(0).(http.ResponseWriter)
				w.WriteHeader(http.StatusOK)
			}).Return()
			mockController.On("Create", mock.Anything, mock.Anything).Maybe().Run(func(args mock.Arguments) {
				w := args.Get(0).(http.ResponseWriter)
				w.WriteHeader(http.StatusCreated)
			}).Return()
			mockController.On("GetByID", mock.Anything, mock.Anything).Maybe().Run(func(args mock.Arguments) {
				w := args.Get(0).(http.ResponseWriter)
				w.WriteHeader(http.StatusOK)
			}).Return()
			mockController.On("List", mock.Anything, mock.Anything).Maybe().Run(func(args mock.Arguments) {
				w := args.Get(0).(http.ResponseWriter)
				w.WriteHeader(http.StatusOK)
			}).Return()

			router := New(mockController)

			req, err := http.NewRequest(tt.method, tt.path, nil)
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)
		})
	}
}

func TestRouter_InvalidRoutes(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		path         string
		expectedCode int
	}{
		{
			name:         "Nonexistent route",
			method:       "GET",
			path:         "/nonexistent",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "Method not allowed on healthz",
			method:       "POST",
			path:         "/app/healthz",
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name:         "Method not allowed on videos",
			method:       "DELETE",
			path:         "/videos",
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name:         "Method not allowed on videos by ID",
			method:       "PUT",
			path:         "/videos/123",
			expectedCode: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockController := NewMockController(t)
			router := New(mockController)

			req, err := http.NewRequest(tt.method, tt.path, nil)
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)
		})
	}
}

func TestRouter_PathParameters(t *testing.T) {
	t.Run("Extracts path parameter", func(t *testing.T) {
		mockController := NewMockController(t)
		mockController.On("GetByID", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			req := args.Get(1).(*http.Request)
			vars := mux.Vars(req)
			assert.Equal(t, "test-id-123", vars["id"])
		}).Return()

		router := New(mockController)

		req, err := http.NewRequest("GET", "/videos/test-id-123", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})
}
