package convert

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseRequest(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr assert.ErrorAssertionFunc
		check   func(t *testing.T, req *http.Request)
	}{
		{
			name: "success: simple GET request",
			input: "GET / HTTP/1.1\r\n" +
				"Host: example.com\r\n" +
				"User-Agent: test\r\n" +
				"\r\n",
			wantErr: assert.NoError,
			check: func(t *testing.T, req *http.Request) {
				assert.Equal(t, "GET", req.Method)
				assert.Equal(t, "/", req.URL.Path)
				assert.Equal(t, "HTTP/1.1", req.Proto)
				assert.Equal(t, "example.com", req.Host)
				assert.Equal(t, "test", req.Header.Get("User-Agent"))
			},
		},
		{
			name: "success: POST request with body",
			input: "POST /submit HTTP/1.1\r\n" +
				"Host: api.example.com\r\n" +
				"Content-Type: application/json\r\n" +
				"Content-Length: 16\r\n" +
				"\r\n" +
				`{"key": "value"}`,
			wantErr: assert.NoError,
			check: func(t *testing.T, req *http.Request) {
				assert.Equal(t, "POST", req.Method)
				assert.Equal(t, "/submit", req.URL.Path)
				assert.Equal(t, "application/json", req.Header.Get("Content-Type"))

				if req.Body != nil {
					buf := new(bytes.Buffer)
					buf.ReadFrom(req.Body)
					body := buf.String()
					assert.Equal(t, `{"key": "value"}`, body)
				}
			},
		},
		{
			name: "success: request with query parameters",
			input: "GET /search?q=golang&page=1 HTTP/1.1\r\n" +
				"Host: google.com\r\n" +
				"\r\n",
			wantErr: assert.NoError,
			check: func(t *testing.T, req *http.Request) {
				assert.Equal(t, "q=golang&page=1", req.URL.RawQuery)
			},
		},
		{
			name: "error: invalid request line: too much",
			input: "GET / HTTP/1.1 extra\r\n" +
				"Host: example.com\r\n" +
				"\r\n",
			wantErr: assert.Error,
			check:   nil,
		},
		{
			name: "error: invalid request line: not enough",
			input: "GET /\r\n" +
				"Host: example.com\r\n" +
				"\r\n",
			wantErr: assert.Error,
			check:   nil,
		},
		{
			name:    "error: empty request",
			input:   "",
			wantErr: assert.Error,
			check:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			req, err := ParseRequest(reader)

			tt.wantErr(t, err)

			if tt.check != nil {
				tt.check(t, req)
			}
		})
	}
}

func splitParts(t *testing.T, output string) (headerPart string, headers map[string]string, bodyPart string) {
	parts := strings.SplitN(output, "\r\n\r\n", 2)
	assert.NotZero(t, len(parts))
	headerPart = parts[0]
	if len(parts) > 1 {
		bodyPart = parts[1]
	}

	headers = make(map[string]string)
	headerLines := strings.Split(headerPart, "\r\n")
	headerPart = headerLines[0]
	headerLines = headerLines[1:]
	for _, line := range headerLines {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, ": ", 2)
		if len(parts) == 2 {
			headers[parts[0]] = parts[1]
		}
	}
	return
}

func TestWriteRequest(t *testing.T) {
	tests := []struct {
		name    string
		request *http.Request
		check   func(t *testing.T, output string)
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success: simple GET",
			request: &http.Request{
				Method: "GET",
				URL:    &url.URL{Path: "/"},
				Proto:  "HTTP/1.1",
				Host:   "example.com",
				Header: http.Header{
					"User-Agent": []string{"test-client"},
				},
			},
			check: func(t *testing.T, output string) {
				header, headers, body := splitParts(t, output)

				assert.Equal(t, "GET / HTTP/1.1", header)
				assert.Equal(t, "example.com", headers["Host"])
				assert.Equal(t, "test-client", headers["User-Agent"])
				assert.Zero(t, body)
			},
			wantErr: assert.NoError,
		},
		{
			name: "success: minimal query",
			request: &http.Request{
				URL: &url.URL{Host: "example.com", Path: "/"},
			},
			check: func(t *testing.T, output string) {
				header, headers, body := splitParts(t, output)

				assert.Equal(t, "GET / HTTP/1.1", header)
				assert.Equal(t, "example.com", headers["Host"])
				assert.Zero(t, body)
			},
			wantErr: assert.NoError,
		},
		{
			name: "success: GET with query params",
			request: &http.Request{
				Method: "GET",
				URL: &url.URL{
					Path:     "/query",
					RawQuery: "abc=1&param=xyz",
				},
				Proto: "HTTP/1.1",
				Host:  "example.com",
				Header: http.Header{
					"User-Agent": []string{"test-client"},
				},
			},
			check: func(t *testing.T, output string) {
				header, headers, body := splitParts(t, output)

				assert.Equal(t, "GET /query?abc=1&param=xyz HTTP/1.1", header)
				assert.Equal(t, "example.com", headers["Host"])
				assert.Equal(t, "test-client", headers["User-Agent"])
				assert.Zero(t, body)
			},
			wantErr: assert.NoError,
		},
		{
			name: "success: POST request with body",
			request: &http.Request{
				Method: "POST",
				URL:    &url.URL{Scheme: "http", Host: "api.example.com", Path: "/submit"},
				Proto:  "HTTP/1.1",
				Host:   "api.example.com",
				Header: http.Header{
					"Content-Type":   []string{"application/json"},
					"Content-Length": []string{"16"},
				},
				Body: io.NopCloser(strings.NewReader(`{"key": "value"}`)),
			},
			check: func(t *testing.T, output string) {
				header, headers, body := splitParts(t, output)

				assert.True(t, strings.HasPrefix(header, "POST /submit HTTP/1.1"))
				assert.Equal(t, "api.example.com", headers["Host"])
				assert.Equal(t, "application/json", headers["Content-Type"])
				assert.Equal(t, "16", headers["Content-Length"])
				assert.Equal(t, `{"key": "value"}`, body)
			},
			wantErr: assert.NoError,
		},
		{
			name: "success: POST request with content length in request",
			request: &http.Request{
				Method:        "POST",
				URL:           &url.URL{Scheme: "http", Host: "api.example.com", Path: "/submit"},
				Proto:         "HTTP/1.1",
				Host:          "api.example.com",
				ContentLength: 16,
				Header: http.Header{
					"Content-Type": []string{"application/json"},
				},
				Body: io.NopCloser(strings.NewReader(`{"key": "value"}`)),
			},
			check: func(t *testing.T, output string) {
				header, headers, body := splitParts(t, output)

				assert.True(t, strings.HasPrefix(header, "POST /submit HTTP/1.1"))
				assert.Equal(t, "api.example.com", headers["Host"])
				assert.Equal(t, "application/json", headers["Content-Type"])
				assert.Equal(t, "16", headers["Content-Length"])
				assert.Equal(t, `{"key": "value"}`, body)
			},
			wantErr: assert.NoError,
		},
		{ // if no content-length provided - we should use transfer-encoding with one chunk of data
			name: "success: POST with no content-length, should use transfer-encoding",
			request: &http.Request{
				Method: "POST",
				URL:    &url.URL{Scheme: "http", Host: "api.example.com", Path: "/submit"},
				Proto:  "HTTP/1.1",
				Host:   "api.example.com",
				Body:   io.NopCloser(strings.NewReader(`{"key": "value"}`)),
			},
			check: func(t *testing.T, output string) {
				header, headers, body := splitParts(t, output)

				assert.True(t, strings.HasPrefix(header, "POST /submit HTTP/1.1"))
				assert.Equal(t, "api.example.com", headers["Host"])
				assert.Equal(t, "chunked", headers["Transfer-Encoding"])
				assert.Equal(t, "16\r\n{\"key\": \"value\"}\r\n0", body)
			},
			wantErr: assert.NoError,
		},
		{
			name: "success: POST request with too long body, cropped",
			request: &http.Request{
				Method: "POST",
				URL:    &url.URL{Scheme: "http", Host: "api.example.com", Path: "/submit"},
				Proto:  "HTTP/1.1",
				Host:   "api.example.com",
				Header: http.Header{
					"Content-Length": []string{"16"},
				},
				Body: io.NopCloser(strings.NewReader(`{"key": "value"}_too_long_body`)),
			},
			check: func(t *testing.T, output string) {
				header, headers, body := splitParts(t, output)

				assert.True(t, strings.HasPrefix(header, "POST /submit HTTP/1.1"))
				assert.Equal(t, "api.example.com", headers["Host"])
				assert.Equal(t, "16", headers["Content-Length"])
				assert.Equal(t, `{"key": "value"}`, body)
			},
			wantErr: assert.NoError,
		},
		{
			name: "error: too short body",
			request: &http.Request{
				Method:        "POST",
				URL:           &url.URL{Scheme: "http", Host: "api.example.com", Path: "/submit"},
				Proto:         "HTTP/1.1",
				Host:          "api.example.com",
				ContentLength: 26,
				Header: http.Header{
					"Content-Type": []string{"application/json"},
				},
				Body: io.NopCloser(strings.NewReader(`{"key": "value"}`)),
			},
			wantErr: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := WriteRequest(&buf, tt.request)

			tt.wantErr(t, err)

			if tt.check != nil {
				tt.check(t, buf.String())
			}
		})
	}
}

func TestParseResponse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr assert.ErrorAssertionFunc
		check   func(t *testing.T, resp *http.Response)
	}{
		{
			name: "200 OK response",
			input: "HTTP/1.1 200 OK\r\n" +
				"Content-Type: text/plain\r\n" +
				"Content-Length: 13\r\n" +
				"\r\n" +
				"Hello, World!",
			wantErr: assert.NoError,
			check: func(t *testing.T, resp *http.Response) {
				assert.Equal(t, 200, resp.StatusCode)
				assert.Equal(t, "200 OK", resp.Status)
				assert.Equal(t, "HTTP/1.1", resp.Proto)
				assert.Equal(t, "text/plain", resp.Header.Get("Content-Type"))

				if resp.Body != nil {
					buf := new(bytes.Buffer)
					buf.ReadFrom(resp.Body)
					body := buf.String()
					assert.Equal(t, "Hello, World!", body)
				}
			},
		},
		{
			name: "404 Not Found response",
			input: "HTTP/1.1 404 Not Found\r\n" +
				"\r\n",
			wantErr: assert.NoError,
			check: func(t *testing.T, resp *http.Response) {
				assert.Equal(t, 404, resp.StatusCode)
				assert.Equal(t, "404 Not Found", resp.Status)
			},
		},
		{
			name: "invalid status line",
			input: "HTTP/1.1 200\r\n" +
				"Content-Type: text/plain\r\n" +
				"\r\n",
			wantErr: assert.Error,
			check:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			resp, err := ParseResponse(reader)

			tt.wantErr(t, err)

			if !tt.wantErr(t, err) && tt.check != nil {
				tt.check(t, resp)
			}
		})
	}
}

func TestWriteResponse(t *testing.T) {
	tests := []struct {
		name     string
		response *http.Response
		check    func(t *testing.T, output string)
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			name: "200 OK response",
			response: &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Proto:      "HTTP/1.1",
				Header: http.Header{
					"Content-Type":   []string{"text/plain"},
					"Content-Length": []string{"13"},
				},
				Body: io.NopCloser(strings.NewReader("Hello, World!")),
			},
			check: func(t *testing.T, output string) {
				header, headers, body := splitParts(t, output)

				assert.True(t, strings.HasPrefix(header, "HTTP/1.1 200 200 OK"))
				assert.Equal(t, "text/plain", headers["Content-Type"])
				assert.Equal(t, "13", headers["Content-Length"])
				assert.Equal(t, "Hello, World!", body)
			},
			wantErr: assert.NoError,
		},
		{
			name: "response without body",
			response: &http.Response{
				StatusCode: 204,
				Status:     "No Content",
				Proto:      "HTTP/1.1",
				Header:     http.Header{},
			},
			check: func(t *testing.T, output string) {
				header, headers, body := splitParts(t, output)

				assert.Equal(t, "HTTP/1.1 204 No Content", header)
				assert.Empty(t, headers)
				assert.Zero(t, body)
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := WriteResponse(&buf, tt.response)

			tt.wantErr(t, err)

			if tt.check != nil {
				tt.check(t, buf.String())
			}
		})
	}
}

func TestRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		req  *http.Request
	}{
		{
			name: "GET request roundtrip",
			req: &http.Request{
				Method: "GET",
				URL:    &url.URL{Path: "/test", RawQuery: "a=1&b=2"},
				Proto:  "HTTP/1.1",
				Host:   "example.com",
				Header: http.Header{
					"User-Agent":    []string{"test-agent"},
					"Authorization": []string{"Bearer token"},
				},
			},
		},
		{
			name: "POST request with body roundtrip",
			req: &http.Request{
				Method: "POST",
				URL:    &url.URL{Path: "/api"},
				Proto:  "HTTP/1.1",
				Host:   "api.example.com",
				Header: http.Header{
					"Content-Type": []string{"application/json"},
				},
				Body: io.NopCloser(strings.NewReader(`{"data": "value"}`)),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := WriteRequest(&buf, tt.req)
			assert.NoError(t, err)

			parsedReq, err := ParseRequest(&buf)
			assert.NoError(t, err)

			assert.Equal(t, tt.req.Method, parsedReq.Method)
			assert.Equal(t, tt.req.URL.Path, parsedReq.URL.Path)
			assert.Equal(t, tt.req.Host, parsedReq.Host)

			for key, wantValues := range tt.req.Header {
				gotValues := parsedReq.Header[key]
				assert.Equal(t, len(wantValues), len(gotValues))
				for i, want := range wantValues {
					if i < len(gotValues) {
						assert.Equal(t, want, gotValues[i])
					}
				}
			}
		})
	}
}

func TestResponseRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		resp *http.Response
	}{
		{
			name: "200 OK roundtrip",
			resp: &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Proto:      "HTTP/1.1",
				Header: http.Header{
					"Content-Type":   []string{"application/json"},
					"Content-Length": []string{"25"},
				},
				Body: io.NopCloser(strings.NewReader(`{"message": "success"}`)),
			},
		},
		{
			name: "error response roundtrip",
			resp: &http.Response{
				StatusCode: 500,
				Status:     "500 Internal Server Error",
				Proto:      "HTTP/1.1",
				Header: http.Header{
					"Content-Type": []string{"text/plain"},
				},
				Body: io.NopCloser(strings.NewReader("Server error")),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := WriteResponse(&buf, tt.resp)
			assert.NoError(t, err)

			parsedResp, err := ParseResponse(&buf)
			assert.NoError(t, err)

			assert.Equal(t, tt.resp.StatusCode, parsedResp.StatusCode)
			assert.Equal(t, tt.resp.Proto, parsedResp.Proto)

			for key, wantValues := range tt.resp.Header {
				gotValues := parsedResp.Header[key]
				assert.Equal(t, len(wantValues), len(gotValues))
			}
		})
	}
}
