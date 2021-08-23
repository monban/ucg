package main

import (
	"testing"

	"github.com/matryer/is"
)

func TestNewGame(t *testing.T) {
	t.Skip("TODO")
}

func TestListedGame(t *testing.T) {
	is := is.New(t)
	const gameName = "GameName"
	g := NewGame(0, gameName, &Player{})
	lg := g.ToListedGame()
	is.Equal(lg.Name, gameName)
}

func TestPlayerCount(t *testing.T) {
	is := is.New(t)
	g := NewGame(0, "Foo", &Player{Name: "Game Owner"})
	newPlayer := &Player{Name: "Second Player"}
	g.AddPlayer(newPlayer)
	is.Equal(len(g.players)+1, g.ToListedGame().Players)
}

func TestPvgs(t *testing.T) {
	is := is.New(t)
	gameName := "Foo"
	p1 := &Player{Name: "Game Owner"}
	p2 := &Player{Name: "Second Player"}
	g := NewGame(0, gameName, p1)
	g.AddPlayer(p2)

	pvgs := g.PlayerViewGameState()
	is.Equal([]string{p1.Name, p2.Name}, pvgs.Players)
	is.Equal(gameName, pvgs.Name)
}
