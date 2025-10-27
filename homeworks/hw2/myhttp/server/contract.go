package server

import "net/http"

type HTTPServer interface {
	ListenAndServe(addr string, handler http.Handler) error
	Close() error
}
