package main

type playerId uint64

type player struct {
	Id   playerId `json:"id"`
	Name string   `json:"name"`
}
