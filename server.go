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
	router      router
	log         printfer
	gm          gameManager
	pm          playerManager
	rootHandler http.Handler
}

type router interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
	Handle(string, string, http.Handler)
	HandleFunc(string, string, func(http.ResponseWriter, *http.Request))
}

type printfer interface {
	Printf(string, ...interface{})
}

type gameManager interface {
	ListGames() []ListedGame
	CreateGame(string, *Player) *Game
}

type playerManager interface {
	FindPlayer(PlayerId) (*Player, error)
	NewPlayer(string) *Player
}

func newServer(l printfer, gm gameManager, pm playerManager) (*server, error) {
	l.Printf("Setting up new server")
	s := &server{
		router:      &methodRouter{},
		log:         l,
		gm:          gm,
		rootHandler: http.FileServer(http.Dir("www")),
		pm:          pm,
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

func (s *server) getGames() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := json.Marshal(s.gm.ListGames())
		if err != nil {
			http.Error(w, "Error marshaling games", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(data))
	}
}

func (s *server) createGame() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		postData, _ := io.ReadAll(r.Body)
		var jd newGameData
		err := json.Unmarshal(postData, &jd)
		if err != nil {
			log.Printf("ERROR: failed to unmarshal postData: %v", err.Error())
			http.Error(w, "Unable to create game with provided data", http.StatusBadRequest)
			return
		}
		p, err := s.pm.FindPlayer(jd.PlayerId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("%v creating new game: %v\n", p, jd)
		g := s.gm.CreateGame(jd.Name, p)
		if g == nil {
			http.Error(w, "Unable to create game", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		body, _ := json.Marshal(g.PlayerViewGameState())
		w.Write(body)
	}
}

func (s *server) newUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "You can only POST to this endpoint", http.StatusMethodNotAllowed)
			return
		}
		body, _ := io.ReadAll(r.Body)
		s.log.Printf("%v", string(body))
		pd := Player{}
		err := json.Unmarshal(body, &pd)
		if err != nil {
			http.Error(w, "Invalid format", http.StatusBadRequest)
			return
		}
		p := s.pm.NewPlayer(pd.Name)
		s.log.Printf("Creating new player: %v(%d)", p.Name, p.Id)
		rbody, _ := json.Marshal(p)
		w.WriteHeader(http.StatusCreated)
		w.Write(rbody)
	}
}

func (s *server) urlForGame(id gameId) url.URL {
	// TODO: Check game with id exists
	p := fmt.Sprintf("/games/%d", id)
	u := url.URL{}
	u.Path = p
	return u
}
