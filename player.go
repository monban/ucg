package main

import "encoding/json"

type playerId uint64

type Player struct {
	id   playerId
	name string
}

func newPlayer(id playerId, name string) *Player {
	return &Player{
		id:   id,
		name: name,
	}
}

func (p *Player) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Id   playerId `json:"id"`
		Name string   `json:"name"`
	}{
		Id:   p.id,
		Name: p.name,
	})
}
