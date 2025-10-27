package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/phayes/freeport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	defaultPort = 0
)

func Test_myServer_ListenAndServe(t *testing.T) {
	tests := []struct {
		name      string
		handler   http.Handler
		doRequest func(addr string) (*http.Response, error)
		check     func(t *testing.T, resp *http.Response)
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name: "success: basic GET request",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/hello", r.URL.Path)
				w.Header().Set("Content-Type", "text/plain")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("world"))
			}),
			doRequest: func(addr string) (*http.Response, error) {
				return http.Get(fmt.Sprintf("http://%s/hello", addr))
			},
			check: func(t *testing.T, resp *http.Response) {
				require.Equal(t, http.StatusOK, resp.StatusCode)
				assert.Equal(t, "text/plain", resp.Header.Get("Content-Type"))
				data, err := io.ReadAll(resp.Body)
				require.NoError(t, err)
				assert.Equal(t, "world", string(data))
			},
			wantErr: assert.NoError,
		},
		{
			name: "success: POST request with body",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				body, _ := io.ReadAll(r.Body)
				assert.Equal(t, "ping", string(body))
				w.WriteHeader(http.StatusCreated)
				_, _ = w.Write([]byte("pong"))
			}),
			doRequest: func(addr string) (*http.Response, error) {
				return http.Post(fmt.Sprintf("http://%s/data", addr), "text/plain", io.NopCloser(io.Reader(&stringReader{"ping"})))
			},
			check: func(t *testing.T, resp *http.Response) {
				require.Equal(t, http.StatusCreated, resp.StatusCode)
				data, err := io.ReadAll(resp.Body)
				require.NoError(t, err)
				assert.Equal(t, "pong", string(data))
			},
			wantErr: assert.NoError,
		},
		{
			name:    "error: not found",
			handler: nil, // должен использовать http.DefaultServeMux
			doRequest: func(addr string) (*http.Response, error) {
				return http.Get(fmt.Sprintf("http://%s/unknown", addr))
			},
			check: func(t *testing.T, resp *http.Response) {
				require.Equal(t, http.StatusNotFound, resp.StatusCode)
				data, _ := io.ReadAll(resp.Body)
				assert.Contains(t, string(data), "404 page not found")
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := New()

			var err error
			port := defaultPort
			if port == 0 {
				port, err = freeport.GetFreePort()
				require.NoError(t, err)
			}
			addr := fmt.Sprintf("localhost:%v", port)

			go func() {
				err := srv.ListenAndServe(addr, tt.handler)
				if err != nil {
					log.Println("server error", err)
				}
			}()
			defer srv.Close()

			// ждём пока сервер поднимется
			time.Sleep(time.Millisecond * 100)

			resp, err := tt.doRequest(addr)
			tt.wantErr(t, err)
			if resp != nil && resp.Body != nil {
				defer resp.Body.Close()

			}

			if tt.check != nil {
				tt.check(t, resp)
			}
		})
	}
}

// stringReader — маленький io.Reader для простых строк
type stringReader struct {
	s string
}

func (r *stringReader) Read(p []byte) (int, error) {
	if len(r.s) == 0 {
		return 0, io.EOF
	}
	n := copy(p, r.s)
	r.s = r.s[n:]
	return n, nil
}
