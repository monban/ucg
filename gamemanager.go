package main

type GameManager struct {
	games []*Game
}

func (gm *GameManager) CreateGame() *Game {
	g := NewGame()
	gm.games = append(gm.games, g)
	return g
}

func (gm *GameManager) ListGames() []ListedGame {
	list := make([]ListedGame, 0, len(gm.games))
	for _, g := range gm.games {
		list = append(list, g.ToListedGame())
	}
	return list
}
