package server

import (
	"net/http"
)

func New() HTTPServer {
	return &myServer{}
}

type myServer struct {
}

func (m *myServer) ListenAndServe(addr string, handler http.Handler) error {
	panic("TODO: implement me")
}

func (m *myServer) Close() error {
	panic("TODO: implement me")
}
