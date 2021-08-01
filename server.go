package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type server struct {
	router      *http.ServeMux
	log         *log.Logger
	gameManager *GameManager
	rootHandler http.Handler
}

func newServer(l *log.Logger, gm *GameManager) (*server, error) {
	l.Println("Setting up new server")
	s := &server{
		router:      http.NewServeMux(),
		log:         l,
		gameManager: gm,
		rootHandler: http.FileServer(http.Dir("www")),
	}
	s.routes()
	return s, nil
}

func (s *server) handleRoot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		i, err := ioutil.ReadFile("index.html")
		if err != nil {
			http.NotFound(w, r)
			return
		}
		fmt.Fprint(w, string(i))
	}
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.log.Printf("%v %v %v%v", r.Proto, r.Method, r.Host, r.URL)
	s.router.ServeHTTP(w, r)
}

func (s *server) gameController() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			data, err := json.Marshal(s.gameManager.ListGames())
			if err != nil {
				http.Error(w, "Error marshaling games", 500)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, string(data))
		case http.MethodPost:
			postData, _ := io.ReadAll(r.Body)
			var jd newGameData
			err := json.Unmarshal(postData, &jd)
			if err != nil {
				log.Printf("ERROR: failed to unmarshal postData: %v", err.Error())
			}
			log.Printf("Creating new game: %v\n", jd)
			g, _ := s.gameManager.CreateGame(jd.Name).JsonGameState()
			if g == nil {
				http.Error(w, "Unable to create game", 500)
			}
			w.WriteHeader(http.StatusCreated)
			s.log.Printf("%v\n", string(g))
			w.Write(g)
		}
	}
}

func (s *server) urlForGame(id gameId) url.URL {
	// TODO: Check game with id exists
	p := fmt.Sprintf("/games/%d", id)
	u := url.URL{}
	u.Path = p
	return u
}

type cardsRouter struct {
}

type gameId uint64
