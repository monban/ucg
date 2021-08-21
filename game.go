package main

import (
	"encoding/json"
)

type Game struct {
	players []*player
	rounds  []*Round
	name    string
	id      gameId
	owner   *player
}

type ListedGame struct {
	Name    string `json:"name"`
	Players int    `json:"players"`
	Owner   string `json:"owner"`
}

func NewGame(id gameId, n string, owner *player) *Game {
	return &Game{
		players: make([]*player, 0),
		rounds:  make([]*Round, 0),
		name:    n,
		id:      id,
		owner:   owner,
	}
}

func (g *Game) AddPlayer(name string) *player {
	newPlayer := player{
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

func (g *Game) JsonGameState() []byte {
	foo := struct {
		Rounds  []interface{} `json:"rounds"`
		Players []string      `json:"players"`
		Name    string        `json:"name"`
	}{}

	foo.Name = g.name

	for _, v := range g.rounds {
		foo.Rounds = append(foo.Rounds, v.ToJson())
	}
	//for _, v := range g.players {
	//foo.Players = append(foo.Players, v)
	//}
	foo.Players = []string{"Fred", "Barney"}
	data, err := json.Marshal(foo)
	if err != nil {
		panic("Unable to marshal game to json")
	}
	return data
}

func (g *Game) ToListedGame() ListedGame {
	return ListedGame{
		Name:    g.name,
		Players: 5,
		Owner:   g.owner.Name,
	}
}

type newGameData struct {
	Name     string   `json:"name"`
	PlayerId playerId `json:"playerId"`
}
