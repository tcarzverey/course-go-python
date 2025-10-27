package convert

import (
	"io"
	"net/http"
)

// ParseRequest парсит HTTP запрос из потока байт
func ParseRequest(r io.Reader) (*http.Request, error) {
	// TODO: implement HTTP request parsing
	panic("TODO: implement me")
}

// WriteRequest записывает HTTP запрос в поток байт
func WriteRequest(w io.Writer, req *http.Request) error {
	// TODO: implement HTTP request writing
	panic("TODO: implement me")
}

// ParseResponse парсит HTTP ответ из потока байт
func ParseResponse(r io.Reader) (*http.Response, error) {
	// TODO: implement HTTP response parsing
	panic("TODO: implement me")
}

// WriteResponse записывает HTTP ответ в поток байт
func WriteResponse(w io.Writer, resp *http.Response) error {
	// TODO: implement HTTP response writing
	panic("TODO: implement me")
}
