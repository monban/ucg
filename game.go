package main

type Game struct {
	players []*Player
	rounds  []*Round
	name    string
	id      gameId
	owner   *Player
}

type PlayerViewGameState struct {
	Rounds  []Round  `json:"rounds"`
	Players []string `json:"players"`
	Name    string   `json:"name"`
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

func (g *Game) PlayerViewGameState() PlayerViewGameState {
	pvgs := PlayerViewGameState{Name: g.name}
	pvgs.Rounds = make([]Round, len(g.rounds))
	for i := range g.rounds {
		pvgs.Rounds[i] = *g.rounds[i]
	}
	pvgs.Players = make([]string, len(g.players)+1)
	pvgs.Players[0] = g.owner.Name
	for i := range g.players {
		pvgs.Players[i+1] = g.players[i].Name
	}
	return pvgs
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
