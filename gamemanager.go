package main

import "fmt"

const maxGameId gameId = ^gameId(0)

type gameId uint64

type GameManager struct {
	games map[gameId]*Game
}

func NewGameManager() *GameManager {
	return &GameManager{
		games: make(map[gameId]*Game, 100),
	}
}

func (gm *GameManager) CreateGame(name string, owner *Player) *Game {
	id := gm.nextGameId()
	g := NewGame(id, name, owner)

	gm.games[id] = g
	return g
}

func (gm *GameManager) List() []*Game {
	n := len(gm.games)
	list := make([]*Game, n, n)
	i := 0
	for _, g := range gm.games {
		list[i] = g
		i++
	}
	return list
}

func (gm *GameManager) Get(id gameId) (*Game, error) {
	game, ok := gm.games[id]
	if !ok {
		return nil, fmt.Errorf("Cannot find game with id %d", id)
	}
	return game, nil
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
