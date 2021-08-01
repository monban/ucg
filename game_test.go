package main

import (
	"testing"
)

func TestNewGame(t *testing.T) {
	t.Skip("TODO")
}

func TestListedGame(t *testing.T) {
	const gameName = "GameName"
	g := NewGame(gameName)
	lg := g.ToListedGame()
	if lg.Name != gameName {
		t.Errorf("Expected %v but got %v", gameName, lg.Name)
	}
}
