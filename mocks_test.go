package main

type MockPlayerManager struct {
	FindPlayerCall struct {
		Receives struct {
			Pid PlayerId
		}
		Returns struct {
			player *Player
			err    error
		}
	}
	NewPlayerCall struct {
		Receives struct {
			Name string
		}
		Returns struct {
			player *Player
		}
	}
}

func (pm *MockPlayerManager) FindPlayer(pid PlayerId) (*Player, error) {
	pm.FindPlayerCall.Receives.Pid = pid
	return pm.FindPlayerCall.Returns.player, pm.FindPlayerCall.Returns.err
}

func (pm *MockPlayerManager) NewPlayer(name string) *Player {
	pm.NewPlayerCall.Receives.Name = name
	return pm.FindPlayerCall.Returns.player
}
