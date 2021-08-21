package main

import (
	"fmt"
	"net/http"
)

type methodRouter struct {
	method map[string]*http.ServeMux
}

func (mr *methodRouter) HandleFunc(method string, path string, handler func(http.ResponseWriter, *http.Request)) {
	if mr.method == nil {
		mr.method = make(map[string]*http.ServeMux)
	}
	mux, ok := mr.method[method]
	if !ok {
		mux = http.NewServeMux()
		mr.method[method] = mux
	}
	mux.HandleFunc(path, handler)
}

func (mr *methodRouter) Handle(method string, path string, handler http.Handler) {
	if mr.method == nil {
		mr.method = make(map[string]*http.ServeMux)
	}
	mux, ok := mr.method[method]
	if !ok {
		mux = http.NewServeMux()
		mr.method[method] = mux
	}
	mux.Handle(path, handler)
}

func (mr *methodRouter) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if mr.method == nil {
		err := fmt.Sprintf("No route found for %v %v", req.Method, req.URL)
		http.Error(res, err, http.StatusNotFound)
		return
	}
	mux, ok := mr.method[req.Method]
	if !ok {
		err := fmt.Sprintf("No route found for %v %v", req.Method, req.URL)
		http.Error(res, err, http.StatusNotFound)
		return
	}
	mux.ServeHTTP(res, req)
}
