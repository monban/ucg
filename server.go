package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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
	ListGames() []*Game
	CreateGame(string, *Player) *Game
	AddPlayerToGame(*Player, gameId) error
	Get(gameId) (*Game, error)
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

func (s *server) postGamesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		player, err := PlayerFromReq(r, s.pm)
		if err != nil {
			s.log.Printf(err.Error())
			http.Error(w, "Error occured reading body", http.StatusBadRequest)
			return
		}
		pathElements := strings.Split(r.URL.Path, "/")
		s.log.Printf("Last element of path is %v", pathElements[len(pathElements)-1])
		if pathElements[len(pathElements)-1] == "join" {
			gidint, err := strconv.ParseUint(pathElements[len(pathElements)-2], 10, 64)
			if err != nil {
				s.log.Printf(err.Error())
				http.Error(w, "Error occured parsing game id", http.StatusBadRequest)
			}
			s.gm.AddPlayerToGame(player, gameId(gidint))
			w.WriteHeader(http.StatusNoContent)

			return
		}
	}
}

func (s *server) getGamesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pathElements := strings.Split(r.URL.Path, "/")
		if len(pathElements) != 3 {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		idInt, err := strconv.ParseUint(pathElements[len(pathElements)-1], 10, 64)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		g, err := s.gm.Get(gameId(idInt))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		data, err := json.Marshal(g)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}

func (s *server) createGameHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		p, err := PlayerFromReq(r, s.pm)
		if err != nil {
			s.log.Printf(err.Error())
			http.Error(w, "Unable to create game with this user", http.StatusBadRequest)
			return
		}
		postData, _ := io.ReadAll(r.Body)
		var jd newGameData
		err = json.Unmarshal(postData, &jd)
		if err != nil {
			log.Printf("ERROR: failed to unmarshal postData: %v", err.Error())
			http.Error(w, "Unable to create game with provided data", http.StatusBadRequest)
			return
		}
		s.log.Printf("%v creating new game: %v\n", p, jd)
		g := s.gm.CreateGame(jd.Name, p)
		if g == nil {
			http.Error(w, "Unable to create game", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		body, _ := json.Marshal(g)
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

func JsonBody(r *http.Request, d interface{}) error {
	var err error
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, d)
	if err != nil {
		return err
	}
	return nil
}

func PlayerFromReq(r *http.Request, pm playerManager) (*Player, error) {
	pidstring := r.Header.Get("X-Player-Id")
	pidint, err := strconv.ParseUint(pidstring, 10, 64)
	if err != nil {
		return nil, err
	}
	player, err := pm.FindPlayer(PlayerId(pidint))
	if err != nil {
		return nil, err
	}
	return player, nil
}

type serverMessage struct {
	action string
}
