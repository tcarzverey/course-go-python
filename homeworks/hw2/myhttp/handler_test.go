package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMyHandler(t *testing.T) {
	tests := []struct {
		name            string
		createRequest   func() *http.Request
		expectedStatus  int
		expectedBody    string
		expectedHeaders map[string]string
	}{
		{
			name: "GET without auth header",
			createRequest: func() *http.Request {
				req, _ := http.NewRequest(http.MethodGet, "/?name=John", nil)
				return req
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "POST without auth header",
			createRequest: func() *http.Request {
				body := `{"name":"John"}`
				req, _ := http.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(body)))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "PUT method not allowed",
			createRequest: func() *http.Request {
				req, _ := http.NewRequest(http.MethodPut, "/", nil)
				req.Header.Set("Authorization", "Bearer token")
				return req
			},
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name: "DELETE method not allowed",
			createRequest: func() *http.Request {
				req, _ := http.NewRequest(http.MethodDelete, "/", nil)
				req.Header.Set("Authorization", "Bearer token")
				return req
			},
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name: "GET without name parameter",
			createRequest: func() *http.Request {
				req, _ := http.NewRequest(http.MethodGet, "/", nil)
				req.Header.Set("Authorization", "Bearer token")
				return req
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "GET with empty name parameter",
			createRequest: func() *http.Request {
				req, _ := http.NewRequest(http.MethodGet, "/?name=", nil)
				req.Header.Set("Authorization", "Bearer token")
				return req
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "GET with name parameter - success",
			createRequest: func() *http.Request {
				req, _ := http.NewRequest(http.MethodGet, "/?name=John", nil)
				req.Header.Set("Authorization", "Bearer token")
				return req
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "Hello, John!",
			expectedHeaders: map[string]string{
				"Content-Type":    "text/plain",
				"X-Custom-Result": "success",
			},
		},
		{
			name: "POST without name field",
			createRequest: func() *http.Request {
				body := `{}`
				req, _ := http.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(body)))
				req.Header.Set("Authorization", "Bearer token")
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "POST with empty name field",
			createRequest: func() *http.Request {
				body := `{"name":""}`
				req, _ := http.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(body)))
				req.Header.Set("Authorization", "Bearer token")
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "POST with invalid JSON",
			createRequest: func() *http.Request {
				body := `invalid json`
				req, _ := http.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(body)))
				req.Header.Set("Authorization", "Bearer token")
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "POST with name field - success",
			createRequest: func() *http.Request {
				body := `{"name":"John"}`
				req, _ := http.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(body)))
				req.Header.Set("Authorization", "Bearer token")
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"greeting":"Hello, John!"}`,
			expectedHeaders: map[string]string{
				"Content-Type":    "application/json",
				"X-Custom-Result": "success",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := tt.createRequest()
			rr := httptest.NewRecorder()

			MyHandler(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code, "Status code mismatch")

			if tt.expectedBody != "" {
				assert.Equal(t, tt.expectedBody, rr.Body.String(), "Response body mismatch")
			}

			for header, expectedValue := range tt.expectedHeaders {
				assert.Equal(t, expectedValue, rr.Header().Get(header), "Header %s mismatch", header)
			}
		})
	}
}
