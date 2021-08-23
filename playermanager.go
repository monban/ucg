package main

import "fmt"

const maxPlayerId PlayerId = ^PlayerId(0)

type PlayerManager struct {
	players map[PlayerId]*Player
}

func newPlayerManager() *PlayerManager {
	return &PlayerManager{
		players: make(map[PlayerId]*Player, 100),
	}
}

func (pm *PlayerManager) setPlayer(p *Player) {
	pm.players[p.Id] = p
}

func (pm *PlayerManager) FindPlayer(id PlayerId) (*Player, error) {
	p, ok := pm.players[id]
	if ok {
		return p, nil
	} else {
		return nil, fmt.Errorf("Cannot find player with id %d", id)
	}
}

func (pm *PlayerManager) NewPlayer(name string) *Player {
	var nextId PlayerId
	for nextId = 0; nextId < maxPlayerId; nextId++ {
		if pm.players[nextId] == nil {
			newPlayer := &Player{Id: nextId, Name: name}
			pm.players[nextId] = newPlayer
			return newPlayer
		}
	}
	panic("Out of game ids!")
}
