package client

import (
	"net/http"
)

func New() HTTPClient {
	return &myClient{}
}

type myClient struct{}

func (m *myClient) Do(req *http.Request) (*http.Response, error) {
	panic("TODO: implement me")
}
