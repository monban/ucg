package main

import (
	"fmt"
	"github.com/charmbracelet/log"
	"net/http"
	"os"
)

func main() {
	if err := (runServer()); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func runServer() error {
	gm := NewGameManager()
	pm := newPlayerManager()
	l := log.Default()
	s, _ := newServer(l, gm, pm)
	return http.ListenAndServe(":8080", s)
}
