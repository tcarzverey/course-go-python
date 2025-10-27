package server

import (
	"io"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_myResponseWriter(t *testing.T) {
	tests := []struct {
		name    string
		prepare func(t *testing.T, w http.ResponseWriter)
		check   func(t *testing.T, resp *http.Response)
	}{
		{
			name:    "success: default values",
			prepare: func(t *testing.T, w http.ResponseWriter) {},
			check: func(t *testing.T, resp *http.Response) {
				assert.Equal(t, http.StatusOK, resp.StatusCode)
				assert.Equal(t, "200 OK", resp.Status)
				assert.Equal(t, "HTTP/1.1", resp.Proto)
				assert.Equal(t, int64(0), resp.ContentLength)
				assert.Equal(t, "0", resp.Header.Get("Content-Length"))
				data, _ := io.ReadAll(resp.Body)
				assert.Empty(t, string(data))
			},
		},
		{
			name: "success: write body only",
			prepare: func(t *testing.T, w http.ResponseWriter) {
				n, err := w.Write([]byte("hello"))
				require.NoError(t, err)
				assert.Equal(t, 5, n)
			},
			check: func(t *testing.T, resp *http.Response) {
				assert.Equal(t, int64(5), resp.ContentLength)
				assert.Equal(t, "5", resp.Header.Get("Content-Length"))

				body, err := io.ReadAll(resp.Body)
				require.NoError(t, err)
				assert.Equal(t, "hello", string(body))
				assert.Equal(t, http.StatusOK, resp.StatusCode)
			},
		},
		{
			name: "success: write header and body",
			prepare: func(t *testing.T, w http.ResponseWriter) {
				w.Header().Set("Content-Type", "text/plain")
				w.WriteHeader(http.StatusAccepted)
				_, _ = w.Write([]byte("response body"))
			},
			check: func(t *testing.T, resp *http.Response) {
				assert.Equal(t, http.StatusAccepted, resp.StatusCode)
				assert.Equal(t, "202 Accepted", resp.Status)
				assert.Equal(t, "text/plain", resp.Header.Get("Content-Type"))
				assert.Equal(t, "HTTP/1.1", resp.Proto)
				assert.Equal(t, int64(len("response body")), resp.ContentLength)
				assert.Equal(t, strconv.Itoa(len("response body")), resp.Header.Get("Content-Length"))
				data, _ := io.ReadAll(resp.Body)
				assert.Equal(t, "response body", string(data))
			},
		},
		{
			name: "success: multiple headers",
			prepare: func(t *testing.T, w http.ResponseWriter) {
				w.Header().Add("X-Test", "a")
				w.Header().Add("X-Test", "b")
				w.WriteHeader(http.StatusCreated)
			},
			check: func(t *testing.T, resp *http.Response) {
				assert.EqualValues(t, []string{"a", "b"}, resp.Header.Values("X-Test"))
				assert.Equal(t, http.StatusCreated, resp.StatusCode)
				assert.Equal(t, "201 Created", resp.Status)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := NewResponseWriter()
			tt.prepare(t, w)

			resp, err := w.GetResponse()
			require.NoError(t, err)
			require.NotNil(t, resp)

			if tt.check != nil {
				tt.check(t, resp)
			}
		})
	}
}
