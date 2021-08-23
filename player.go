package main

import "fmt"

type PlayerId uint64

type Player struct {
	Id   PlayerId `json:"id"`
	Name string   `json:"name"`
}

func (p *Player) String() string {
	return fmt.Sprintf("%v(%d)", p.Name, p.Id)
}
