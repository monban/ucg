package main

type MockGameManager struct {
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
}

func (gm *MockGameManager) ListGames() []ListedGame {
	return gm.ListGamesCall.Returns.Games
}

func (gm *MockGameManager) CreateGame(name string, owner *Player) *Game {
	gm.CreateGameCall.Receives.Name = name
	gm.CreateGameCall.Receives.Owner = owner
	return gm.CreateGameCall.Returns.Game
}
