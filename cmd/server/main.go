package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/charmbracelet/log"
	"github.com/monban/ucg/game"
	"github.com/monban/ucg/server"
)

func main() {
	if err := (runServer()); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func runServer() error {
	gm := game.NewGameManager()
	pm := game.NewPlayerManager()
	l := log.Default()
	s, _ := server.New(l, gm, pm)
	return http.ListenAndServe(":8080", s)
}
