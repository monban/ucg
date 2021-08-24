package main

type MockGameManager struct {
	log           printfer
	ListGamesCall struct {
		Receives struct {
		}
		Returns struct {
			Games []ListedGame
		}
	}
	CreateGameCall struct {
		Receives struct {
			Name  string
			Owner *Player
		}
		Returns struct {
			Game *Game
		}
	}
	AddPlayerToGameCall struct {
		Receives struct {
			p   *Player
			gid gameId
		}
		Returns struct {
			Error error
		}
	}
}

func (gm *MockGameManager) ListGames() []ListedGame {
	return gm.ListGamesCall.Returns.Games
}

func (gm *MockGameManager) CreateGame(name string, owner *Player) *Game {
	gm.CreateGameCall.Receives.Name = name
	gm.CreateGameCall.Receives.Owner = owner
	return gm.CreateGameCall.Returns.Game
}

func (gm *MockGameManager) AddPlayerToGame(p *Player, gid gameId) error {
	gm.log.Printf("Adding player %+v to game with id %d", p, gid)
	gm.AddPlayerToGameCall.Receives.p = p
	gm.AddPlayerToGameCall.Receives.gid = gid
	return gm.AddPlayerToGameCall.Returns.Error
}
