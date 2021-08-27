package main

import (
	"testing"

	"github.com/matryer/is"
)

func TestCreateGame(t *testing.T) {
	is := is.New(t)
	gm := newGameManager()
	p := &Player{}
	g := gm.CreateGame("Test Game", p)
	is.True(g != nil)
}

func TestAddPlayerToGame(t *testing.T) {
	is, gm, g := GameManagerMocks(t)
	p := &Player{Name: "New Player", Id: 1234}
	gm.AddPlayerToGame(p, g.id)
	found := false
	for _, player := range g.players {
		if player.Id == p.Id {
			found = true
		}
	}
	is.True(found)
}

func TestGet(t *testing.T) {
	is, gm, g := GameManagerMocks(t)

	returnedGame, err := gm.GetGamePlayerView(g.id)
	is.NoErr(err)
	is.Equal(returnedGame.Name, g.name)
}

func GameManagerMocks(t *testing.T) (*is.I, *GameManager, *Game) {
	i := is.New(t)
	owner := &Player{Name: "Game Owner"}
	gm := newGameManager()
	game := gm.CreateGame("Test Game", owner)
	return i, gm, game
}
