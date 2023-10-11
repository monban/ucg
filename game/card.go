package game

import "fmt"

type Card struct {
	Id         int
	text       string
	isQuestion bool
}

type JsonCard struct {
	Text       string
	IsQuestion bool
}

func NewCard(i int, t string, b bool) *Card {
	return &Card{
		i,
		t,
		b,
	}
}

func (c *Card) display() {
	var color string
	if c.isQuestion {
		color = "BLACK"
	} else {
		color = "WHITE"
	}
	fmt.Printf("%s CARD\n", color)
	fmt.Println(c.text)
}

func (c *Card) ToJson() JsonCard {
	return JsonCard{
		Text:       c.text,
		IsQuestion: c.isQuestion,
	}
}
