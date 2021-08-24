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
	is.Equal(len(game.PlayerViewGameState().Players), 2)
}
