package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"text/template"

	"github.com/monban/ucg/game"
)

type server struct {
	router      router
	log         logger
	gm          gameManager
	pm          playerManager
	rootHandler http.Handler
	template    *template.Template
}

type router interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
	Handle(string, string, http.Handler)
	HandleFunc(string, string, func(http.ResponseWriter, *http.Request))
}

type logger interface {
	Printf(string, ...interface{})
	Info(any, ...any)
	Error(any, ...any)
}

type gameManager interface {
	List() []*game.Game
	CreateGame(string, *game.Player) *game.Game
	AddPlayerToGame(*game.Player, game.GameId) error
	Get(game.GameId) (*game.Game, error)
}

type playerManager interface {
	FindPlayer(game.PlayerId) (*game.Player, error)
	NewPlayer(string) *game.Player
}

func New(l logger, gm gameManager, pm playerManager) (*server, error) {
	s := &server{
		router:   &methodRouter{},
		log:      l,
		gm:       gm,
		pm:       pm,
		template: template.Must(template.ParseFiles("templates/main.html.tmpl")),
	}
	s.rootHandler = s.handleRoot()
	s.routes()
	return s, nil
}

func (s *server) handleRoot() http.HandlerFunc {
	t := template.Must(s.template.Clone())
	template.Must(t.ParseFiles("templates/index.html.tmpl"))

	return func(w http.ResponseWriter, r *http.Request) {
		t.Execute(w, "n/a")
	}
}

func (s *server) newGameHandler() http.HandlerFunc {
	t := template.Must(s.template.Clone())
	template.Must(t.ParseFiles("templates/games/new.html.tmpl"))

	return func(w http.ResponseWriter, r *http.Request) {
		t.Execute(w, "n/a")
	}
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.log.Info("request", "protocol", r.Proto, "method", r.Method, "host", r.Host, "uri", r.URL)
	s.router.ServeHTTP(w, r)
}

func (s *server) getGames() http.HandlerFunc {
	t := template.Must(s.template.Clone())
	template.Must(t.ParseFiles("templates/games/index.html.tmpl"))

	return func(w http.ResponseWriter, r *http.Request) {
		s.template.Execute(w, s.gm.List())
	}
}

func (s *server) postGamesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		player, err := PlayerFromReq(r, s.pm)
		if err != nil {
			s.log.Error("Error reading body", "err", err.Error())
			http.Error(w, "Error occured reading body", http.StatusBadRequest)
			return
		}
		pathElements := strings.Split(r.URL.Path, "/")
		s.log.Info("Last element of path is %v", pathElements[len(pathElements)-1])
		if pathElements[len(pathElements)-1] == "join" {
			gidint, err := strconv.ParseUint(pathElements[len(pathElements)-2], 10, 64)
			if err != nil {
				s.log.Error("parsing game id", "err", err)
				http.Error(w, "Error occured parsing game id", http.StatusBadRequest)
			}
			s.gm.AddPlayerToGame(player, game.GameId(gidint))
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
		g, err := s.gm.Get(game.GameId(idInt))
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
		// 	var err error
		// 	p, err := PlayerFromReq(r, s.pm)
		// 	if err != nil {
		// 		s.log.Error("PlayerFromReq", "err", err)
		// 		http.Error(w, "Unable to create game with this user", http.StatusBadRequest)
		// 		return
		// 	}
		// 	postData, _ := io.ReadAll(r.Body)
		// 	var jd newGameData
		// 	err = json.Unmarshal(postData, &jd)
		// 	if err != nil {
		// 		s.log.Error("ERROR: failed to unmarshal postData", "err", err)
		// 		http.Error(w, "Unable to create game with provided data", http.StatusBadRequest)
		// 		return
		// 	}
		// s.log.Info("creating new game: %v\n", "player", p, "gameData", jd)
		// g := s.gm.CreateGame(jd.Name, p)
		// if g == nil {
		// 	http.Error(w, "Unable to create game", http.StatusInternalServerError)
		// 	return
		// }
		// w.Header().Set("Content-Type", "application/json")
		// w.WriteHeader(http.StatusCreated)
		// body, _ := json.Marshal(g)
		// w.Write(body)
	}
}

func (s *server) newUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.PostFormValue("name")
		if name == "" {
			s.log.Error("newUser with empty name")
			http.Error(w, "name is required", http.StatusBadRequest)
			return
		}
		p := s.pm.NewPlayer(r.PostFormValue("name"))
		s.log.Info("Creating new player", "name", p.Name, "id", p.Id)
		rbody, _ := json.Marshal(p)
		w.WriteHeader(http.StatusCreated)
		w.Write(rbody)
	}
}

func valuesToString(u url.Values) string {
	b := &strings.Builder{}
	for k, v := range u {
		fmt.Fprintf(b, "%s=>%s ", k, v)
	}
	return b.String()
}

func (s *server) urlForGame(id game.GameId) url.URL {
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

func PlayerFromReq(r *http.Request, pm playerManager) (*game.Player, error) {
	pidstring := r.Header.Get("X-Player-Id")
	pidint, err := strconv.ParseUint(pidstring, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("unable to parse uint from string ''%s': %w", pidstring, err)
	}
	player, err := pm.FindPlayer(game.PlayerId(pidint))
	if err != nil {
		return nil, fmt.Errorf("unable to find player with pid %d: %w", pidint, err)
	}
	return player, nil
}

type serverMessage struct {
	action string
}
