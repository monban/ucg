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

	"github.com/matryer/is"
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
	is := is.New(t)
	gm := GameManager{}
	l := log.Default()
	_, err := newServer(l, &gm)
	is.NoErr(err)
}

func TestCreateGame(t *testing.T) {
	is := is.New(t)
	s := setupServer()
	var jsonData []byte
	var data io.Reader
	var req *http.Request

	// First we try without including a player id
	jsonData, _ = json.Marshal(newGameData{Name: "foo"})
	data = bytes.NewReader(jsonData)
	req = httptest.NewRequest("POST", "/games", data)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)
	is.Equal(w.Result().StatusCode, http.StatusBadRequest)

	// Try again, but include the id
	jsonData, _ = json.Marshal(newGameData{Name: "foo", PlayerId: 0})
	req = httptest.NewRequest("POST", "/games", data)
	s.ServeHTTP(w, req)
	body, _ := io.ReadAll(w.Result().Body)
	t.Logf("word: %v", string(body))
	is.Equal(w.Result().StatusCode, http.StatusCreated)
	var bd struct{ Name string }
	json.Unmarshal(body, &bd)

	is.Equal(bd.Name, "foo")
}

func TestUrlForGame(t *testing.T) {
	s := setupServer()
	g := s.gm.CreateGame("foo", &Player{})
	expectedPath := fmt.Sprintf("/games/%d", g.id)
	u := s.urlForGame(g.id)
	if u.Path != expectedPath {
		t.Errorf("Expected %v, got %v", expectedPath, u.Path)
	}
}

func TestNewUser(t *testing.T) {
	is := is.New(t)
	s := setupServer()
	data := strings.NewReader("{\"name\":\"Bob\"}")
	req := httptest.NewRequest("POST", "/users", data)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)
	is.Equal(w.Result().StatusCode, http.StatusCreated)
	body, _ := io.ReadAll(w.Result().Body)
	var bd struct{ id playerId }
	err := json.Unmarshal(body, &bd)
	is.NoErr(err)
}
