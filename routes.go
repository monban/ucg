package main

func (s *server) routes() {
	s.router.HandleFunc("/", s.handleRoot())
	s.router.HandleFunc("/games", s.gameController())
}
