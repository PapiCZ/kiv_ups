package interfaces

import (
	"kiv_ups_server/net/tcp"
	"kiv_ups_server/net/tcp/protocol"
)

type Lobby struct {
	Name         string
	Owner        Player
	Players      map[PlayerUID]Player
	PlayersLimit int
}

func (l *Lobby) AddPlayer(player Player) {
	l.Players[player.GetUID()] = player
}

func (l *Lobby) GetPlayers() []Player {
	players := make([]Player, 0)

	for _, player := range l.Players {
		players = append(players, player)
	}

	return players
}

func (l *Lobby) KickPlayer(player Player) {
	panic("implement me")
}

func (l *Lobby) KickPlayers() {
	panic("implement me")
}

type MasterServer interface {
	Start() (err error)
	Stop() (err error)
	RunAction(message tcp.ClientMessage) (err error)
	GetTCPServer() *tcp.Server
	GetPlayers() map[tcp.UID]Player
	Authenticate(player Player)
	AddLobby(lobby *Lobby)
	DeleteLobby(name string)
	GetLobby(name string) (*Lobby, error)
	GetLobbies() []*Lobby
	AddGameServer(server GameServer)
	SendMessageWithoutRequest(sm tcp.ServerMessage, player ...Player)
	SendMessage(sm tcp.ServerMessage, requestId protocol.RequestId, player ...Player)
}

type GameServer interface {
	Start()
	AddPlayer(player Player)
	GetPlayers() []Player
	GetRequestMessageChan() chan PlayerMessage
}
