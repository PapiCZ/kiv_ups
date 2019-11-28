package interfaces

import (
	"kiv_ups_server/net/tcp"
)

type PlayerContext int
type PlayerUID int

type Player interface {
	GetTCPClient() *tcp.Client
	GetUID() PlayerUID
	SetUID(uid PlayerUID)
	GetName() string
	SetName(name string)
	GetContext() PlayerContext
	SetContext(ctx PlayerContext)
	SetConnectedLobby(*Lobby)
	GetConnectedLobby() *Lobby
}

type PlayerMessage interface {
	GetMessage() *tcp.ClientMessage
	GetPlayer()  Player
}
