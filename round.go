package main

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Round struct {
	number           int
	judge            *Player
	question         *Card
	answers          []*Card
	players          []*Player
	players_unplayed []playerId
}

type JsonRound struct {
	Number   int
	Judge    *Player
	Question *Card
	Answers  []*Card
}

func NewRound(n int, p *Player, q *Card) *Round {
	return &Round{
		number:   n,
		judge:    p,
		question: q,
		answers:  make([]*Card, 0),
	}
}

func (r *Round) display() {
	fmt.Printf("\n========\nRound %d\nCzar: %s\n", r.number, r.judge.name)
	fmt.Println("Black card:")
	r.question.display()
}

func (r *Round) wordsHard() (string, error) {
	b, _ := json.Marshal(r)
	return fmt.Sprintf("%s", b), nil
}

func (r *Round) ToJson() JsonRound {
	return JsonRound{
		Number:   r.number,
		Judge:    r.judge,
		Question: r.question,
		Answers:  r.answers,
	}
}

func (r *Round) PlayCard(p *Player, c *Card) error {
	// Check if player has already played this round
	for _, n := range r.players_unplayed {
		if n == p.id {
			return errors.New("Player has already played")
		}
	}
	r.answers = append(r.answers, c)
	return nil
}
