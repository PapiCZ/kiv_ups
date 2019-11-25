package interfaces

import (
	"kiv_ups_server/net/tcp"
)

type MasterServer interface {
	Start() (err error)
	Stop() (err error)
	RunAction(message tcp.ClientMessage) (err error)
	GetTCPServer() *tcp.Server
	GetPlayers() map[tcp.UID]Player
	Authenticate(player Player)
}
