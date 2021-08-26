package main

import (
	"testing"

	"github.com/matryer/is"
)

func TestNewGame(t *testing.T) {
	t.Skip("TODO")
}

func TestPlayerCount(t *testing.T) {
	is := is.New(t)
	g := NewGame(0, "Foo", &Player{Name: "Game Owner"})
	newPlayer := &Player{Name: "Second Player"}
	g.AddPlayer(newPlayer)
	is.Equal(len(g.players)+1, g.PlayerView().PlayerCount)
}

func TestPlayerNames(t *testing.T) {
	is := is.New(t)
	playernames := []string{"Alice", "Bob", "Charles"}
	g := NewGame(0, "Foo", &Player{Name: playernames[0]})
	g.AddPlayer(&Player{Name: playernames[1]})
	g.AddPlayer(&Player{Name: playernames[2]})
	is.Equal(playernames, g.PlayerNames())
}
