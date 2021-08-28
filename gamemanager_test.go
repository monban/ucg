package main

import (
	"sort"
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

	returnedGame, err := gm.Get(g.id)
	is.NoErr(err)
	is.Equal(returnedGame.name, g.name)
}

func TestList(t *testing.T) {
	is, gm, g := GameManagerMocks(t)
	p2 := &Player{}
	g2 := gm.CreateGame("Second Game", p2)
	expectedList := []gameId{g.id, g2.id}
	gameList := gm.List()
	t.Logf("Game list: %+v", gameList)
	l := len(gameList)
	idList := make([]gameId, 0, l)
	is.Equal(len(gameList), len(expectedList))

	for _, game := range gameList {
		idList = append(idList, game.id)
	}
	sort.Slice(idList, func(i, j int) bool { return idList[i] < idList[j] })
	is.Equal(expectedList, idList)
}

func GameManagerMocks(t *testing.T) (*is.I, *GameManager, *Game) {
	i := is.New(t)
	owner := &Player{Name: "Game Owner"}
	gm := newGameManager()
	game := gm.CreateGame("Test Game", owner)
	return i, gm, game
}
