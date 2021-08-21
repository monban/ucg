package main

import "fmt"

const maxPlayerId playerId = ^playerId(0)

type playerManager struct {
	players map[playerId]*player
}

func newPlayerManager() *playerManager {
	return &playerManager{
		players: make(map[playerId]*player, 100),
	}
}

func (pm *playerManager) setPlayer(p *player) {
	pm.players[p.Id] = p
}

func (pm *playerManager) findPlayer(id playerId) (*player, error) {
	p, ok := pm.players[id]
	if ok {
		return p, nil
	} else {
		return nil, fmt.Errorf("Cannot find player with id %d", id)
	}
}

func (pm *playerManager) newPlayer(name string) *player {
	var nextId playerId
	for nextId = 0; nextId < maxPlayerId; nextId++ {
		if pm.players[nextId] == nil {
			newPlayer := &player{Id: nextId, Name: name}
			pm.players[nextId] = newPlayer
			return newPlayer
		}
	}
	panic("Out of game ids!")
}
