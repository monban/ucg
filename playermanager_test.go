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
