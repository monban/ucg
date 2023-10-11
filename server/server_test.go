package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
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
	_, err := New(lt, &GameManager{}, &PlayerManager{})
	is.NoErr(err)
}

func TestCreateGameWithoutPlayer(t *testing.T) {
	is, pm, _, srv := ServerMocks(t)

	// Set up mocks
	pm.FindPlayerCall.Returns.player = nil
	pm.FindPlayerCall.Returns.err = errors.New("Player not found")

	// First we try without including a player id

	// Set up the request
	jsonData, _ := json.Marshal(newGameData{Name: "foo"})
	postData := bytes.NewReader(jsonData)
	req := httptest.NewRequest("POST", "/games", postData)
	rec := httptest.NewRecorder()

	// Make the request
	srv.ServeHTTP(rec, req)

	// Check the results
	resultBody, _ := io.ReadAll(rec.Result().Body)
	t.Logf("resultBody: %v", string(resultBody))
	is.Equal(rec.Result().StatusCode, http.StatusBadRequest)

	// This time we'll include a player id
}

func TestCreateGameWithPlayer(t *testing.T) {
	is, pm, gm, srv := ServerMocks(t)

	// Set up mocks
	testPlayer := &Player{Name: "TestPlayer"}
	pm.FindPlayerCall.Returns.player = testPlayer
	pm.FindPlayerCall.Returns.err = nil
	gm.CreateGameCall.Returns.Game = &Game{owner: testPlayer}

	// Set up the request
	gameData := newGameData{Name: "foo"}
	jsonData, _ := json.Marshal(gameData)
	postData := bytes.NewReader(jsonData)
	req := httptest.NewRequest("POST", "/games", postData)
	req.Header.Add("X-Player-Id", fmt.Sprint(testPlayer.Id))
	rec := httptest.NewRecorder()

	// Make the request
	srv.ServeHTTP(rec, req)

	// Check the results
	is.Equal(gm.CreateGameCall.Receives.Name, gameData.Name)
	is.Equal(gm.CreateGameCall.Receives.Owner, testPlayer)
	resultBody, _ := io.ReadAll(rec.Result().Body)
	t.Logf("resultBody: %v", string(resultBody))
	is.Equal(rec.Result().StatusCode, http.StatusCreated)
}

func TestUrlForGame(t *testing.T) {
	_, _, _, srv := ServerMocks(t)
	expectedPath := fmt.Sprintf("/games/%d", 1234)
	u := srv.urlForGame(1234)
	if u.Path != expectedPath {
		t.Errorf("Expected %v, got %v", expectedPath, u.Path)
	}
}

func TestNewUser(t *testing.T) {
	is, pm, _, srv := ServerMocks(t)
	testPlayer := &Player{Name: "TestPlayer"}
	pm.FindPlayerCall.Returns.player = testPlayer
	postData, _ := json.Marshal(testPlayer)
	postBytes := bytes.NewBuffer(postData)
	req := httptest.NewRequest("POST", "/users", postBytes)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	is.Equal(w.Result().StatusCode, http.StatusCreated)
	is.Equal(pm.NewPlayerCall.Receives.Name, "TestPlayer")
	body, _ := io.ReadAll(w.Result().Body)
	t.Logf("Returned body: %v", string(body))
	var bd struct{ id PlayerId }
	err := json.Unmarshal(body, &bd)
	is.NoErr(err)
}

func TestUserCanJoinGame(t *testing.T) {
	is, pm, gm, srv := ServerMocks(t)
	owner := &Player{Name: "Game Owner", Id: 0}
	game := &Game{id: 0, name: "Test Game", owner: owner}
	newPlayer := &Player{Name: "Second Player"}
	pm.FindPlayerCall.Returns.player = newPlayer
	postData, _ := json.Marshal(struct{ PlayerId PlayerId }{newPlayer.Id})

	url := srv.urlForGame(game.id)
	url.Path = fmt.Sprintf("%v/join", url.Path)

	req := httptest.NewRequest("POST", url.Path, bytes.NewBuffer(postData))
	req.Header.Add("X-Player-Id", fmt.Sprint(newPlayer.Id))
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	body, _ := io.ReadAll(rec.Result().Body)
	t.Logf("Resultbody: %+v", string(body))
	t.Logf("StatusCode: %d", rec.Code)
	is.Equal(gm.AddPlayerToGameCall.Receives.p, newPlayer)
}

func TestGetGame(t *testing.T) {
	is, _, gm, srv := ServerMocks(t)
	owner := &Player{}
	game := NewGame(137, "Super Happy Fun Time!", owner)
	gm.GetCall.Returns.game = game

	uri := fmt.Sprintf("/games/%d", game.id)
	req := httptest.NewRequest("GET", uri, nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	is.Equal(rec.Result().StatusCode, http.StatusOK)
	returnedGame := &PlayerViewGame{}
	UnmarshalRecorder(rec, returnedGame)
	is.Equal(rec.Result().Header.Get("Content-Type"), "application/json")
	is.Equal(gm.GetCall.Receives.id, game.id)
	is.Equal(game.PlayerView(), *returnedGame)
}

func ServerMocks(t *testing.T) (*is.I, *MockPlayerManager, *MockGameManager, *server) {
	i := is.NewRelaxed(t)
	lt := &logTesting{t}
	pm := &MockPlayerManager{}
	gm := &MockGameManager{log: lt}
	srv, _ := New(lt, gm, pm)
	return i, pm, gm, srv
}

func UnmarshalRecorder(rec *httptest.ResponseRecorder, d interface{}) {
	body, _ := io.ReadAll(rec.Result().Body)
	json.Unmarshal(body, d)
}
