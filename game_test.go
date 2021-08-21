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
	g := NewGame(0, gameName, &player{})
	lg := g.ToListedGame()
	is.Equal(lg.Name, gameName)
}
