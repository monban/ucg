package main

import (
	"testing"
	"log"
	"net/http"
	"net/http/httptest"
	"encoding/json"
	"bytes"
	"io"
)

func setupServer() *server {
	gm := GameManager{}
	l := log.Default()
	s,err := newServer(l, &gm)
	if err != nil {
		panic(err)
	}
	return s
}

func TestNewServer(t *testing.T) {
	gm := GameManager{}
	l := log.Default()
	_,err := newServer(l, &gm)
	if err != nil {
		t.Error("undefined")
	}
}

type newGameData struct {
	Name string `json:'name'`
}

func TestCreateGame(t *testing.T) {
	s := setupServer()
	jsonData,_ := json.Marshal(newGameData{Name: "foo"})
	data := bytes.NewReader(jsonData)
	req := httptest.NewRequest("POST", "/games", data)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)
	if w.Result().StatusCode != http.StatusCreated {
		t.Errorf("Expected %v, got %v.", http.StatusCreated, w.Code)
	}
	body,_ := io.ReadAll(w.Result().Body)
	var bd struct{Name string}
	json.Unmarshal(body, &bd)
	if bd.Name != "foo" {
		t.Errorf("Expected name to be foo, but it was %v", bd.Name)
	}
}

