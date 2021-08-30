package main

import (
	"testing"

	"github.com/matryer/is"
)

func TestCreateUser(t *testing.T) {
	is := is.New(t)
	pm := newPlayerManager()
	p := pm.NewPlayer("Bob")
	ret, err := pm.FindPlayer(p.Id)
	is.NoErr(err)
	is.Equal(ret.Name, p.Name)
}

func TestFindPlayer_missing(t *testing.T) {
	is := is.New(t)
	pm := newPlayerManager()
	p, err := pm.FindPlayer(1234)
	is.True(err != nil)
	is.True(p == nil)
}

func TestFindPlayer_extant(t *testing.T) {
	is := is.New(t)
	pm := newPlayerManager()
	p := pm.NewPlayer("Test Player")
	p, err := pm.FindPlayer(p.Id)
	is.NoErr(err)
	is.True(p != nil)
}
