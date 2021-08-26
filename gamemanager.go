package main

import "fmt"

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

func (gm *GameManager) CreateGame(name string, owner *Player) *Game {
	g := NewGame(gm.nextGameId(), name, owner)
	gm.games[0] = g
	return g
}

func (gm *GameManager) ListGames() []*Game {
	n := len(gm.games)
	list := make([]*Game, 0, n)
	i := 0
	for _, g := range gm.games {
		i++
		list[i] = g
	}
	return list
}

func (gm *GameManager) GetGamePlayerView(id gameId) (*PlayerViewGame, error) {
	game, ok := gm.games[id]
	if !ok {
		return nil, fmt.Errorf("Cannot find game with id %d", id)
	}
	view := game.PlayerView()
	return &view, nil
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

func (gm *GameManager) AddPlayerToGame(p *Player, gid gameId) error {
	game, ok := gm.games[gid]
	if !ok {
		return fmt.Errorf("Game with id %v not found", gid)
	}
	game.AddPlayer(p)
	return nil
}
