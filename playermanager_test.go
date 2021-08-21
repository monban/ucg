package main

import (
	"testing"

	"github.com/matryer/is"
)

func TestCreateUser(t *testing.T) {
	is := is.New(t)
	pm := newPlayerManager()
	p := pm.newPlayer("Bob")
	ret, err := pm.findPlayer(p.Id)
	is.NoErr(err)
	is.Equal(ret.Name, p.Name)
}
