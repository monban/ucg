package main

const maxGameId gameId = ^gameId(0)

type gameId uint64

type GameManager struct {
	games map[gameId]*Game
}

func newGameManager() *GameManager {
	return &GameManager{
		games: make(map[gameId]*Game, 100),
	}
}

func (gm *GameManager) CreateGame(name string, owner *player) *Game {
	g := NewGame(gm.nextGameId(), name, owner)
	gm.games[0] = g
	return g
}

func (gm *GameManager) ListGames() []ListedGame {
	list := make([]ListedGame, 0, len(gm.games))
	for _, g := range gm.games {
		list = append(list, g.ToListedGame())
	}
	return list
}

func (gm *GameManager) nextGameId() gameId {
	var nextId gameId
	for nextId = 0; nextId < maxGameId; nextId++ {
		if gm.games[nextId] == nil {
			return nextId
		}
	}
	panic("Out of game ids!")
}
