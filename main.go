package main

import (
	"fmt"
	"log"
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
	l := log.Default()
	s, _ := newServer(l, gm)
	return http.ListenAndServe(":8080", s)
}
