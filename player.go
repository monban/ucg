package main

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type PlayerId uint64

type Player struct {
	Id   PlayerId
	Name string
}

func (p *Player) String() string {
	return fmt.Sprintf("%v(%d)", p.Name, p.Id)
}

func (p *Player) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}{
		Id:   strconv.FormatUint(uint64(p.Id), 10),
		Name: p.Name,
	})
}

func (p *Player) UnmarshalJSON(b []byte) error {
	d := &struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}{}
	err := json.Unmarshal(b, d)
	if err != nil {
		return err
	}
	p.Name = d.Name
	pid, err := strconv.ParseUint(d.Id, 10, 64)
	if err != nil {
		return err
	}
	p.Id = PlayerId(pid)
	return nil
}
