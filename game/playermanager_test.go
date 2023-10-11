package game

import (
	"testing"

	"github.com/matryer/is"
)

func TestCreateUser(t *testing.T) {
	is := is.New(t)
	pm := NewPlayerManager()
	p := pm.NewPlayer("Bob")
	ret, err := pm.FindPlayer(p.Id)
	is.NoErr(err)
	is.Equal(ret.Name, p.Name)
}

func TestPlayerIdNotPredictable(t *testing.T) {
	is := is.New(t)
	pm := NewPlayerManager()
	p := pm.NewPlayer("Test Player")
	t.Logf("New player id is %d", p.Id)
	is.True(p.Id != 0)
}

func TestFindPlayer_missing(t *testing.T) {
	is := is.New(t)
	pm := NewPlayerManager()
	p, err := pm.FindPlayer(1234)
	is.True(err != nil)
	is.True(p == nil)
}

func TestFindPlayer_extant(t *testing.T) {
	is := is.New(t)
	pm := NewPlayerManager()
	p := pm.NewPlayer("Test Player")
	p, err := pm.FindPlayer(p.Id)
	is.NoErr(err)
	is.True(p != nil)
}
