package main

func (s *server) routes() {
	s.router.Handle("/", s.rootHandler)
	s.router.HandleFunc("/games", s.gameController())
	s.router.HandleFunc("/users", s.newUser())
}
