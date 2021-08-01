package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setupServer() *server {
	gm := GameManager{}
	l := log.Default()
	s, err := newServer(l, &gm)
	if err != nil {
		panic(err)
	}
	return s
}

func TestNewServer(t *testing.T) {
	gm := GameManager{}
	l := log.Default()
	_, err := newServer(l, &gm)
	if err != nil {
		t.Error("undefined")
	}
}

func TestCreateGame(t *testing.T) {
	s := setupServer()
	jsonData, _ := json.Marshal(newGameData{Name: "foo"})
	data := bytes.NewReader(jsonData)
	req := httptest.NewRequest("POST", "/games", data)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)
	if w.Result().StatusCode != http.StatusCreated {
		t.Errorf("Expected %v, got %v.", http.StatusCreated, w.Code)
	}
	body, _ := io.ReadAll(w.Result().Body)
	var bd struct{ Name string }
	json.Unmarshal(body, &bd)
	if bd.Name != "foo" {
		t.Errorf("Expected name to be foo, but it was %v", bd.Name)
	}
}

func TestUrlForGame(t *testing.T) {
	s := setupServer()
	g := s.gameManager.CreateGame("SomeName")
	expectedPath := fmt.Sprintf("/games/%d", g.id)
	u := s.urlForGame(g.id)
	if u.Path != expectedPath {
		t.Errorf("Expected %v, got %v", expectedPath, u.Path)
	}
}

func TestNewUser(t *testing.T) {
	s := setupServer()
	data := strings.NewReader("{\"name\":\"Bob\"}")
	req := httptest.NewRequest("POST", "/users", data)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)
	if w.Result().StatusCode != http.StatusCreated {
		t.Errorf("Expected %v, got %v.", http.StatusCreated, w.Code)
	}
	body, _ := io.ReadAll(w.Result().Body)
	var bd struct{ id userId }
	err := json.Unmarshal(body, &bd)
	if err != nil {
		t.Errorf("Unable to unmarshall %v to playerId", string(body))
	}
}
