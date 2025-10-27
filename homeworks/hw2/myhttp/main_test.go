package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/phayes/freeport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tcarzverey/course-go-python/homeworks/hw2/myhttp/client"
	"github.com/tcarzverey/course-go-python/homeworks/hw2/myhttp/server"
)

const (
	defaultPort = 0
)

func TestMyServerHandlerClient(t *testing.T) {
	tests := []struct {
		name       string
		getRequest func(addr string) *http.Request
		check      func(t *testing.T, resp *http.Response)
	}{
		{
			name: "GET without auth header - should return 401",
			getRequest: func(addr string) *http.Request {
				url := fmt.Sprintf("http://%s/myhandler?name=John", addr)
				req, _ := http.NewRequest("GET", url, nil)
				return req
			},
			check: func(t *testing.T, resp *http.Response) {
				assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
			},
		},
		{
			name: "POST without auth header - should return 401",
			getRequest: func(addr string) *http.Request {
				url := fmt.Sprintf("http://%s/myhandler", addr)
				body := `{"name":"John"}`
				req, _ := http.NewRequest("POST", url, bytes.NewReader([]byte(body)))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			check: func(t *testing.T, resp *http.Response) {
				assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
			},
		},
		{
			name: "PUT method not allowed - should return 405",
			getRequest: func(addr string) *http.Request {
				url := fmt.Sprintf("http://%s/myhandler", addr)
				req, _ := http.NewRequest("PUT", url, nil)
				req.Header.Set("Authorization", "Bearer token")
				return req
			},
			check: func(t *testing.T, resp *http.Response) {
				assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
			},
		},
		{
			name: "GET without name param - should return 400",
			getRequest: func(addr string) *http.Request {
				url := fmt.Sprintf("http://%s/myhandler", addr)
				req, _ := http.NewRequest("GET", url, nil)
				req.Header.Set("Authorization", "Bearer token")
				return req
			},
			check: func(t *testing.T, resp *http.Response) {
				assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
			},
		},
		{
			name: "GET with auth and name - should return 200 with plain text",
			getRequest: func(addr string) *http.Request {
				url := fmt.Sprintf("http://%s/myhandler?name=John", addr)
				req, _ := http.NewRequest("GET", url, nil)
				req.Header.Set("Authorization", "Bearer token")
				return req
			},
			check: func(t *testing.T, resp *http.Response) {
				assert.Equal(t, http.StatusOK, resp.StatusCode)
				assert.Equal(t, "text/plain", resp.Header.Get("Content-Type"))
				assert.Equal(t, "success", resp.Header.Get("X-Custom-Result"))

				buf := new(bytes.Buffer)
				buf.ReadFrom(resp.Body)
				assert.Equal(t, "Hello, John!", buf.String())
			},
		},
		{
			name: "POST with valid JSON - should return 200 with JSON",
			getRequest: func(addr string) *http.Request {
				url := fmt.Sprintf("http://%s/myhandler", addr)
				body := `{"name":"John"}`
				req, _ := http.NewRequest("POST", url, bytes.NewReader([]byte(body)))
				req.Header.Set("Authorization", "Bearer token")
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			check: func(t *testing.T, resp *http.Response) {
				assert.Equal(t, http.StatusOK, resp.StatusCode)
				assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
				assert.Equal(t, "success", resp.Header.Get("X-Custom-Result"))

				buf := new(bytes.Buffer)
				buf.ReadFrom(resp.Body)
				assert.Equal(t, `{"greeting":"Hello, John!"}`, buf.String())
			},
		},
		{
			name: "POST with invalid JSON - should return 400",
			getRequest: func(addr string) *http.Request {
				url := fmt.Sprintf("http://%s/myhandler", addr)
				body := `invalid json`
				req, _ := http.NewRequest("POST", url, bytes.NewReader([]byte(body)))
				req.Header.Set("Authorization", "Bearer token")
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			check: func(t *testing.T, resp *http.Response) {
				assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
			},
		},
		{
			name: "POST with empty name - should return 400",
			getRequest: func(addr string) *http.Request {
				url := fmt.Sprintf("http://%s/myhandler", addr)
				body := `{"name":""}`
				req, _ := http.NewRequest("POST", url, bytes.NewReader([]byte(body)))
				req.Header.Set("Authorization", "Bearer token")
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			check: func(t *testing.T, resp *http.Response) {
				assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
			},
		},
	}
	http.HandleFunc("/myhandler", MyHandler)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := server.New()

			var err error
			port := defaultPort
			if port == 0 {
				port, err = freeport.GetFreePort()
				require.NoError(t, err)
			}
			addr := fmt.Sprintf("localhost:%v", port)

			go func() {
				err := srv.ListenAndServe(addr, nil)
				if err != nil {
					log.Println("server error", err)
				}
			}()
			defer srv.Close()

			// ждём пока сервер поднимется
			time.Sleep(time.Millisecond * 100)

			req := tt.getRequest(addr)
			c := client.New()
			resp, err := c.Do(req)
			assert.NoError(t, err)
			if resp != nil && resp.Body != nil {
				defer resp.Body.Close()
			}

			if tt.check != nil {
				tt.check(t, resp)
			}
		})
	}
}
