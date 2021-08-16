package main

import "fmt"

const maxPlayerId playerId = ^playerId(0)

type playerManager struct {
	players map[playerId]*Player
}

func newPlayerManager() *playerManager {
	return &playerManager{
		players: make(map[playerId]*Player, 100),
	}
}

func (pm *playerManager) setPlayer(p *Player) {
	pm.players[p.id] = p
}

func (pm *playerManager) findPlayer(id playerId) (*Player, error) {
	p, ok := pm.players[id]
	if ok {
		return p, nil
	} else {
		return nil, fmt.Errorf("Cannot find player with id %d", id)
	}
}

func (pm *playerManager) newPlayer(name string) *Player {
	var nextId playerId
	for nextId = 0; nextId < maxPlayerId; nextId++ {
		if pm.players[nextId] == nil {
			newPlayer := newPlayer(nextId, name)
			pm.players[nextId] = newPlayer
			return newPlayer
		}
	}
	panic("Out of game ids!")
}
