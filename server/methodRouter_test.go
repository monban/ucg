package server

import (
	"net/http"
	"testing"
)

func TestAddRoute(t *testing.T) {
	mr := &methodRouter{}
	mr.HandleFunc("GET", "/", func(http.ResponseWriter, *http.Request) {})
}
