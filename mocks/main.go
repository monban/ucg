package mocks

import ucg "github.com/monban/ucg"

type PlayerManager struct {
	FindPlayerCall struct {
		Receives struct {
		}
		Returns struct {
			Pid ucg.PlayerId
		}
	}
}

func (pm *PlayerManager) FindPlayer(pid PlayerId) (uint64, error) {
	return pm.FindPlayerCall.Returns.Pid, nil
}
