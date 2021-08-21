package main

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Round struct {
	number           int
	judge            *player
	question         *Card
	answers          []*Card
	players          []*player
	players_unplayed []playerId
}

type JsonRound struct {
	Number   int
	Judge    *player
	Question *Card
	Answers  []*Card
}

func NewRound(n int, p *player, q *Card) *Round {
	return &Round{
		number:   n,
		judge:    p,
		question: q,
		answers:  make([]*Card, 0),
	}
}

func (r *Round) display() {
	fmt.Printf("\n========\nRound %d\nCzar: %s\n", r.number, r.judge.Name)
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

func (r *Round) PlayCard(p *player, c *Card) error {
	// Check if player has already played this round
	for _, n := range r.players_unplayed {
		if n == p.Id {
			return errors.New("Player has already played")
		}
	}
	r.answers = append(r.answers, c)
	return nil
}
