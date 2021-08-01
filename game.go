package main

import (
	"encoding/json"
)

type Game struct {
	players []*Player
	rounds  []*Round
	name    string
	id      gameId
}

type ListedGame struct {
	Name    string `json:"name"`
	Players int    `json:"players"`
}

func NewGame(n string) *Game {
	return &Game{
		players: make([]*Player, 0),
		rounds:  make([]*Round, 0),
		name:    n,
	}
}

func (g *Game) AddPlayer(name string) *Player {
	newPlayer := Player{
		Id:   0,
		Name: name,
	}
	g.players = append(g.players, &newPlayer)
	return &newPlayer
}

func (g *Game) StartNewRound() *Round {
	foo := NewCard(0, "foo", true)
	r := NewRound(len(g.rounds), g.players[0], foo)
	g.rounds = append(g.rounds, r)
	return r
}

func (g *Game) JsonGameState() ([]byte, error) {
	foo := struct {
		Rounds  []interface{}
		Players []interface{}
		Name    string
	}{}

	foo.Name = g.name

	for _, v := range g.rounds {
		foo.Rounds = append(foo.Rounds, v.ToJson())
	}
	for _, v := range g.players {
		foo.Players = append(foo.Players, v)
	}
	data, err := json.Marshal(foo)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (g *Game) ToListedGame() ListedGame {
	return ListedGame{
		Name:    g.name,
		Players: 5,
	}
}

type newGameData struct {
	Name string `json:'name'`
}
