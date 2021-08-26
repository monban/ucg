package main

import "encoding/json"

type Game struct {
	players []*Player
	rounds  []*Round
	name    string
	id      gameId
	owner   *Player
}

type PlayerViewGame struct {
	Name        string   `json:"name"`
	PlayerCount int      `json:"playerCount"`
	PlayerNames []string `json:"playerNames"`
	Owner       string   `json:"owner"`
	Id          gameId   `json:"id"`
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

func (g *Game) PlayerNames() []string {
	c := len(g.players) + 1
	names := make([]string, c, c)
	names[0] = g.owner.Name
	if c == 1 {
		return names
	}
	for i := 1; i < c; i++ {
		names[i] = g.players[i-1].Name
	}
	return names
}

func (g *Game) PlayerView() PlayerViewGame {
	return PlayerViewGame{
		Name:        g.name,
		PlayerCount: len(g.players) + 1,
		Owner:       g.owner.Name,
		PlayerNames: g.PlayerNames(),
		Id:          g.id,
	}
}

func (g *Game) MarshalJSON() ([]byte, error) {
	return json.Marshal(g.PlayerView())
}

type newGameData struct {
	Name     string   `json:"name"`
	PlayerId PlayerId `json:"PlayerId"`
}
