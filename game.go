package main

import (
	"encoding/json"
)

type Game struct {
	players []*Player
	rounds  []*Round
	name    string
	id      gameId
	owner   *Player
}

type ListedGame struct {
	Name    string `json:"name"`
	Players int    `json:"players"`
	Owner   string `json:"owner"`
}

func NewGame(id gameId, n string, owner *Player) *Game {
	return &Game{
		players: make([]*Player, 0),
		rounds:  make([]*Round, 0),
		name:    n,
		id:      id,
		owner:   owner,
	}
}

func (g *Game) AddPlayer(p *Player) {
	g.players = append(g.players, p)
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
		Players: len(g.players) + 1,
		Owner:   g.owner.Name,
	}
}

type newGameData struct {
	Name     string   `json:"name"`
	PlayerId PlayerId `json:"PlayerId"`
}
