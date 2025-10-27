package server

import (
	"net/http"
)

type MyResponseWriter struct {
}

// NewResponseWriter создает новый MyResponseWriter
func NewResponseWriter() *MyResponseWriter {
	return &MyResponseWriter{}
}

// implement MyResponseWriter methods for http.ResponseWriter
var _ http.ResponseWriter = (*MyResponseWriter)(nil)

func (w *MyResponseWriter) Header() http.Header {
	panic("TODO: implement me")
}

func (w *MyResponseWriter) Write(data []byte) (int, error) {
	panic("TODO: implement me")
}

func (w *MyResponseWriter) WriteHeader(statusCode int) {
	panic("TODO: implement me")
}

// implement method for using your ResponseWriter on server

func (w *MyResponseWriter) GetResponse() (*http.Response, error) {
	panic("TODO: implement me")
}
