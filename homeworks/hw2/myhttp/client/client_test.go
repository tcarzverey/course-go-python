package client

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"slices"
	"testing"
	"time"

	"github.com/phayes/freeport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	defaultPort = 0
)

func Test_myClient_Do_RealHTTP(t *testing.T) {
	type args = struct {
		req            *http.Request
		handlerPattern string
		handler        func(http.ResponseWriter, *http.Request)
	}

	tests := []struct {
		name    string
		args    args
		check   func(t *testing.T, response *http.Response)
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success: simple GET",
			args: args{
				req: &http.Request{
					URL: &url.URL{
						Scheme: "http",
						Host:   "localhost",
						Path:   "/query",
					},
				},
				handlerPattern: "/query",
				handler:        func(rw http.ResponseWriter, r *http.Request) {},
			},
			check: func(t *testing.T, response *http.Response) {
				assert.Equal(t, http.StatusOK, response.StatusCode)
				assert.Equal(t, "200 OK", response.Status)
				assert.Equal(t, "HTTP/1.1", response.Proto)
				assert.Equal(t, 1, response.ProtoMajor)
				assert.Equal(t, 1, response.ProtoMinor)
				assert.Zero(t, response.ContentLength)
			},
			wantErr: assert.NoError,
		},
		{
			name: "success: simple POST with body",
			args: args{
				req: &http.Request{
					URL: &url.URL{
						Scheme: "http",
						Host:   "localhost",
						Path:   "/post",
					},
					Body: io.NopCloser(bytes.NewReader([]byte(`{"key": "value"}`))),
					Header: map[string][]string{
						"Content-Length": {"16"},
						"Content-Type":   {"application/json"},
					},
				},
				handlerPattern: "/post",
				handler: func(rw http.ResponseWriter, r *http.Request) {
					buf, err := io.ReadAll(r.Body)
					if err != nil {
						rw.WriteHeader(http.StatusInternalServerError)
					}
					if string(buf) != `{"key": "value"}` {
						rw.WriteHeader(http.StatusBadRequest)
					}

					rw.Header().Add("Response-Header", "abc")
					rw.WriteHeader(http.StatusAccepted)
					_, _ = rw.Write([]byte("response"))
				},
			},
			check: func(t *testing.T, response *http.Response) {
				assert.Equal(t, http.StatusAccepted, response.StatusCode)
				assert.Equal(t, "202 Accepted", response.Status)
				assert.Equal(t, "abc", response.Header.Get("Response-Header"))
				assert.Equal(t, int64(8), response.ContentLength)
				bytes, err := io.ReadAll(response.Body)
				assert.NoError(t, err)
				assert.Equal(t, "response", string(bytes))
			},
			wantErr: assert.NoError,
		},
		{
			name: "success: simple POST with body and query",
			args: args{
				req: &http.Request{
					URL: &url.URL{
						Scheme:   "http",
						Host:     "localhost",
						Path:     "/post-with-query",
						RawQuery: "param1=abc&param2=2",
					},
					Body: io.NopCloser(bytes.NewReader([]byte(`{"key": "value"}`))),
					Header: map[string][]string{
						"Content-Length": {"16"},
						"Content-Type":   {"application/json"},
					},
				},
				handlerPattern: "/post-with-query",
				handler: func(rw http.ResponseWriter, r *http.Request) {
					buf, err := io.ReadAll(r.Body)
					if err != nil {
						rw.WriteHeader(http.StatusInternalServerError)
					}
					if string(buf) != `{"key": "value"}` {
						rw.WriteHeader(http.StatusBadRequest)
					}
					if r.URL.Query().Get("param1") != "abc" {
						rw.WriteHeader(http.StatusBadRequest)
					}
					if r.URL.Query().Get("param2") != "2" {
						rw.WriteHeader(http.StatusBadRequest)
					}

					rw.Header().Add("Response-Header", "abc")
					rw.WriteHeader(http.StatusAccepted)
					_, _ = rw.Write([]byte("response"))
				},
			},
			check: func(t *testing.T, response *http.Response) {
				assert.Equal(t, http.StatusAccepted, response.StatusCode)
				assert.Equal(t, "202 Accepted", response.Status)
				assert.Equal(t, "abc", response.Header.Get("Response-Header"))
				assert.Equal(t, int64(8), response.ContentLength)
				bytes, err := io.ReadAll(response.Body)
				assert.NoError(t, err)
				assert.Equal(t, "response", string(bytes))
			},
			wantErr: assert.NoError,
		},
		{
			name: "success: multiple headers",
			args: args{
				req: &http.Request{
					URL: &url.URL{
						Scheme: "http",
						Host:   "localhost",
						Path:   "/multi-header",
					},
					Header: func() http.Header {
						h := http.Header{}
						h.Add("X-Test", "1")
						h.Add("X-Test", "2")
						return h
					}(),
				},
				handlerPattern: "/multi-header",
				handler: func(rw http.ResponseWriter, r *http.Request) {
					values := r.Header.Values("X-Test")
					if !slices.Equal(values, []string{"1", "2"}) && !slices.Equal(values, []string{"2", "1"}) {
						rw.WriteHeader(http.StatusBadRequest)
					}
					rw.Header().Add("X-Resp-Test", "3")
					rw.Header().Add("X-Resp-Test", "4")
					rw.WriteHeader(http.StatusOK)
				},
			},
			check: func(t *testing.T, resp *http.Response) {
				assert.Equal(t, http.StatusOK, resp.StatusCode)
				assert.ElementsMatch(t, []string{"3", "4"}, resp.Header.Values("X-Resp-Test"))
			},
			wantErr: assert.NoError,
		},
		{
			name: "error: nil request",
			args: args{
				req: nil,
			},
			wantErr: assert.Error,
		},
		{
			name: "error: request without url",
			args: args{
				req: &http.Request{
					URL: nil,
				},
			},
			wantErr: assert.Error,
		},
		{
			name: "error: not found handler",
			args: args{
				req: &http.Request{
					URL: &url.URL{
						Scheme: "http",
						Host:   "localhost",
						Path:   "/nonexsistent",
					},
				},
				handlerPattern: "/query",
				handler: func(rw http.ResponseWriter, r *http.Request) {

				},
			},
			check: func(t *testing.T, response *http.Response) {
				assert.Equal(t, http.StatusNotFound, response.StatusCode)
				assert.Equal(t, "404 Not Found", response.Status)
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			port := defaultPort
			if port == 0 {
				port, err = freeport.GetFreePort()
				require.NoError(t, err)
			}

			var serverCloseFunc func() error
			go func() {
				if tt.args.handlerPattern == "" {
					return
				}
				mux := http.NewServeMux()
				mux.HandleFunc(tt.args.handlerPattern, tt.args.handler)
				server := http.Server{
					Addr:    fmt.Sprintf(":%v", port),
					Handler: mux,
				}
				serverCloseFunc = server.Close
				if err := server.ListenAndServe(); err != nil {
					log.Println("server error", err)
				}
			}()

			if serverCloseFunc != nil {
				defer serverCloseFunc()
			}

			time.Sleep(time.Millisecond * 100) // дожидаемся чтобы сервер запустился

			if tt.args.req != nil && tt.args.req.Host != "" {
				tt.args.req.Host = fmt.Sprintf("%v:%v", tt.args.req.Host, port)
			}
			if tt.args.req != nil && tt.args.req.URL != nil && tt.args.req.URL.Host != "" {
				tt.args.req.URL.Host = fmt.Sprintf("%v:%v", tt.args.req.URL.Host, port)
			}

			client := New()
			resp, err := client.Do(tt.args.req)

			tt.wantErr(t, err)
			if resp != nil && resp.Body != nil {
				defer resp.Body.Close()
			}
			if tt.check != nil {
				tt.check(t, resp)
				if t.Failed() {
					fmt.Printf("testing %s | additionalLog: response=%+v\n", tt.name, resp)
				}
			}

		})
	}
}
