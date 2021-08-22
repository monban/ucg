package main

type playerId uint64

type Player struct {
	Id   playerId `json:"id"`
	Name string   `json:"name"`
}
