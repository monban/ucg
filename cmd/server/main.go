package main

import (
	"net/http"
	"os"

	"github.com/charmbracelet/log"
	"github.com/monban/ucg/game"
	"github.com/monban/ucg/server"
)

func main() {
	if err := (runServer()); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func runServer() error {
	addr := "0.0.0.0:8080"
	gm := game.NewGameManager()
	pm := game.NewPlayerManager()
	l := log.Default()
	s, _ := server.New(l, gm, pm)
	l.Info("Setting up new server", "addr", addr)
	return http.ListenAndServe(addr, s)
}
