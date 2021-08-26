package main

import (
	"testing"

	"github.com/matryer/is"
)

func TestCreateGameFromManager(t *testing.T) {
	t.Log("TODO")
}

func TestAddPlayerToGame(t *testing.T) {
	is := is.New(t)
	owner := &Player{Name: "Game Owner"}
	player := &Player{Name: "Second Player"}
	gm := newGameManager()
	game := gm.CreateGame("Test Game", owner)
	gm.AddPlayerToGame(player, game.id)
	is.Equal(game.PlayerView().PlayerCount, 2)
}

func TestGetGamePlayerView(t *testing.T) {
	is := is.New(t)
	owner := &Player{Name: "First Player"}
	gm := newGameManager()
	game := gm.CreateGame("Test Game", owner)

	returnedGame, err := gm.GetGamePlayerView(game.id)
	is.NoErr(err)
	is.Equal(returnedGame.Name, game.name)
}
