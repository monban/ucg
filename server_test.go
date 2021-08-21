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
	var rec *httptest.ResponseRecorder
	var p player

	// First we try without including a player id
	jsonData, _ = json.Marshal(newGameData{Name: "foo"})
	data = bytes.NewReader(jsonData)
	req = httptest.NewRequest("POST", "/games", data)
	rec = httptest.NewRecorder()
	s.ServeHTTP(rec, req)
	is.Equal(rec.Result().StatusCode, http.StatusBadRequest)

	// Create player
	pData, _ := json.Marshal(struct {
		Name string `json:"name"`
	}{
		Name: "Bob"})
	req = httptest.NewRequest("POST", "/users", bytes.NewBuffer(pData))
	s.ServeHTTP(rec, req)
	returnData, _ := io.ReadAll(rec.Result().Body)
	t.Logf("returnData: %v\n", string(returnData))
	json.Unmarshal(returnData, &p)
	pid := p.Id

	// Create game
	rec = httptest.NewRecorder()
	jsonData, _ = json.Marshal(newGameData{Name: "foo", PlayerId: pid})
	req = httptest.NewRequest("POST", "/games", bytes.NewBuffer(jsonData))
	s.ServeHTTP(rec, req)
	body, _ := io.ReadAll(rec.Result().Body)
	t.Logf("word: %v", string(body))
	is.Equal(rec.Result().StatusCode, http.StatusCreated)
	var bd struct{ Name string }
	json.Unmarshal(body, &bd)

	is.Equal(bd.Name, "foo")
}

func TestUrlForGame(t *testing.T) {
	s := setupServer()
	g := s.gm.CreateGame("foo", &player{})
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
