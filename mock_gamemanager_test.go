package main

type MockGameManager struct {
	log           printfer
	ListGamesCall struct {
		Receives struct {
		}
		Returns struct {
			Games []*Game
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
	GetGamePlayerViewCall struct {
		Receives struct {
			id gameId
		}
		Returns struct {
			game *PlayerViewGame
			err  error
		}
	}
}

func (gm *MockGameManager) ListGames() []*Game {
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

func (gm *MockGameManager) GetGamePlayerView(id gameId) (*PlayerViewGame, error) {
	gm.GetGamePlayerViewCall.Receives.id = id
	return gm.GetGamePlayerViewCall.Returns.game, gm.GetGamePlayerViewCall.Returns.err
}
