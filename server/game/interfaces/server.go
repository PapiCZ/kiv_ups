package interfaces

import (
	"kiv_ups_server/net/tcp"
)

type Lobby struct {
	Name    string
	Players map[PlayerUID]Player
}

type MasterServer interface {
	Start() (err error)
	Stop() (err error)
	RunAction(message tcp.ClientMessage) (err error)
	GetTCPServer() *tcp.Server
	GetPlayers() map[tcp.UID]Player
	Authenticate(player Player)
	AddLobby(lobby Lobby)
}
