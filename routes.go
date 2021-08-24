package main

import "net/http"

func (s *server) routes() {
	s.router.Handle(http.MethodGet, "/", s.rootHandler)
	s.router.HandleFunc(http.MethodGet, "/games", s.getGames())
	s.router.HandleFunc(http.MethodPost, "/games", s.createGameHandler())
	s.router.HandleFunc(http.MethodPost, "/games/", s.postGamesHandler())
	s.router.HandleFunc(http.MethodPost, "/users", s.newUser())
}
