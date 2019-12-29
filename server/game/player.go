package game

import (
	"kiv_ups_server/game/interfaces"
	"kiv_ups_server/net/tcp"
	"math/rand"
	"time"
)

type Player struct {
	tcpClient      *tcp.Client
	UID            interfaces.PlayerUID
	Name           string
	Context        interfaces.PlayerContext
	ConnectedLobby *interfaces.Lobby
	GameServer     interfaces.GameServer
	LastKeepAlive  int64
}

func NewPlayer(client *tcp.Client, name string, context interfaces.PlayerContext) Player {
	return Player{
		tcpClient: client,
		UID:       interfaces.PlayerUID(rand.Int()),
		Name:      name,
		Context:   context,
	}
}

func NewShadowPlayer(client *tcp.Client, name string, context interfaces.PlayerContext) Player {
	return Player{
		tcpClient: client,
		UID:       0,
		Name:      name,
		Context:   context,
	}
}

func (p *Player) SetTCPClient(tcpClient *tcp.Client) {
	p.tcpClient = tcpClient
}

func (p *Player) GetTCPClient() *tcp.Client {
	return p.tcpClient
}

func (p *Player) SetUID(uid interfaces.PlayerUID) {
	p.UID = uid
}

func (p *Player) GetUID() interfaces.PlayerUID {
	return p.UID
}

func (p *Player) SetName(name string) {
	p.Name = name
}

func (p *Player) GetName() string {
	return p.Name
}

func (p *Player) GetContext() interfaces.PlayerContext {
	return p.Context
}

func (p *Player) SetContext(ctx interfaces.PlayerContext) {
	p.Context = ctx
}

func (p *Player) SetConnectedLobby(lobby *interfaces.Lobby) {
	p.ConnectedLobby = lobby
}

func (p *Player) GetConnectedLobby() *interfaces.Lobby {
	return p.ConnectedLobby
}

func (p *Player) SetGameServer(gs interfaces.GameServer) {
	p.GameServer = gs
}

func (p *Player) GetGameServer() interfaces.GameServer {
	return p.GameServer
}

func (p *Player) IsConnected() bool {
	return time.Now().Unix() < p.LastKeepAlive + 2
}

func (p *Player) RefreshKeepAlive() {
	p.LastKeepAlive = time.Now().Unix()
}
