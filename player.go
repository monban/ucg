package main

type Player struct {
	Id   int
	Name string
}

func NewPlayer(id int, name string) Player {
	return Player{
		Id:   id,
		Name: name,
	}
}
