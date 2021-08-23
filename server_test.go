package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/matryer/is"
)

type logTesting struct {
	t *testing.T
}

func (lt logTesting) Printf(a string, b ...interface{}) {
	lt.t.Logf(a, b...)
}

func TestNewServer(t *testing.T) {
	is := is.New(t)
	lt := logTesting{t: t}
	_, err := newServer(lt, &GameManager{}, &PlayerManager{})
	is.NoErr(err)
}

func TestCreateGame(t *testing.T) {
	var rec *httptest.ResponseRecorder
	var req *http.Request
	var jsonData []byte
	var postData io.Reader
	var resultBody []byte
	var pm *MockPlayerManager = &MockPlayerManager{}
	var is *is.I = is.New(t)

	// Set up mocks
	pm.FindPlayerCall.Returns.player = nil
	pm.FindPlayerCall.Returns.err = errors.New("Player not found")
	gm := newGameManager()
	srv, _ := newServer(&logTesting{t: t}, gm, pm)

	// First we try without including a player id

	// Set up the request
	jsonData, _ = json.Marshal(newGameData{Name: "foo"})
	postData = bytes.NewReader(jsonData)
	req = httptest.NewRequest("POST", "/games", postData)
	rec = httptest.NewRecorder()

	// Make the request
	srv.ServeHTTP(rec, req)

	// Check the results
	resultBody, _ = io.ReadAll(rec.Result().Body)
	t.Logf("resultBody: %v", string(resultBody))
	is.Equal(rec.Result().StatusCode, http.StatusBadRequest)

	// This time we'll include a player id
	// Set up mocks
	pm.FindPlayerCall.Returns.player = &Player{Name: "TestPlayer"}
	pm.FindPlayerCall.Returns.err = nil

	// Set up the request
	jsonData, _ = json.Marshal(newGameData{Name: "foo", PlayerId: 0})
	postData = bytes.NewReader(jsonData)
	req = httptest.NewRequest("POST", "/games", postData)
	rec = httptest.NewRecorder()

	// Make the request
	srv.ServeHTTP(rec, req)

	// Check the results
	resultBody, _ = io.ReadAll(rec.Result().Body)
	t.Logf("resultBody: %v", string(resultBody))

}

func TestUrlForGame(t *testing.T) {
	s, _ := newServer(&logTesting{t}, &GameManager{}, &PlayerManager{})
	g := s.gm.CreateGame("foo", &Player{})
	expectedPath := fmt.Sprintf("/games/%d", g.id)
	u := s.urlForGame(g.id)
	if u.Path != expectedPath {
		t.Errorf("Expected %v, got %v", expectedPath, u.Path)
	}
}

func TestNewUser(t *testing.T) {
	is := is.New(t)
	pm := &MockPlayerManager{}
	pm.FindPlayerCall.Returns.player = &Player{Name: "TestPlayer"}
	s, _ := newServer(&logTesting{t}, &GameManager{}, pm)
	data := strings.NewReader("{\"name\":\"Bob\"}")
	req := httptest.NewRequest("POST", "/users", data)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)
	is.Equal(w.Result().StatusCode, http.StatusCreated)
	is.Equal(pm.NewPlayerCall.Receives.Name, "TestPlayer")
	body, _ := io.ReadAll(w.Result().Body)
	var bd struct{ id PlayerId }
	err := json.Unmarshal(body, &bd)
	is.NoErr(err)
}
