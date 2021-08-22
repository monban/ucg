package main

import "fmt"

const maxPlayerId playerId = ^playerId(0)

type PlayerManager struct {
	players map[playerId]*Player
}

func newPlayerManager() *PlayerManager {
	return &PlayerManager{
		players: make(map[playerId]*Player, 100),
	}
}

func (pm *PlayerManager) setPlayer(p *Player) {
	pm.players[p.Id] = p
}

func (pm *PlayerManager) FindPlayer(id playerId) (*Player, error) {
	p, ok := pm.players[id]
	if ok {
		return p, nil
	} else {
		return nil, fmt.Errorf("Cannot find player with id %d", id)
	}
}

func (pm *PlayerManager) NewPlayer(name string) *Player {
	var nextId playerId
	for nextId = 0; nextId < maxPlayerId; nextId++ {
		if pm.players[nextId] == nil {
			newPlayer := &Player{Id: nextId, Name: name}
			pm.players[nextId] = newPlayer
			return newPlayer
		}
	}
	panic("Out of game ids!")
}
