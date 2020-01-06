package masterserver

import (
	"kiv_ups_server/masterserver/actions"
	"kiv_ups_server/masterserver/interfaces"
	"kiv_ups_server/net/tcp"
	"time"
)

const KeepaliveTimeoutTolerance = 2
const CheatCounterTolerance = 5

// Player structure is a structure that handles that handles
// player's TCP client, name, connected lobby, connected GameServer and last
// keep-alive timestamp that is used to indicate whether is player still connected
type Player struct {
	tcpClient      *tcp.Client
	UID            interfaces.PlayerUID
	Name           string
	Context        interfaces.PlayerContext
	ConnectedLobby *interfaces.Lobby
	GameServer     interfaces.GameServer
	LastKeepAlive  int64
	CheatCounter   int
}

// NewShadowPlayer creates and initializes Player structure with UID 0,
// that is reserved for player that isn't authenticated
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

func (p *Player) LeaveGame() {
	// remove player from game server
	p.GetGameServer().RemovePlayer(p)
	p.SetGameServer(nil)
}

func (p *Player) IsConnected() bool {
	return time.Now().Unix() < p.LastKeepAlive+KeepaliveTimeoutTolerance
}

func (p *Player) RefreshKeepAlive() {
	p.LastKeepAlive = time.Now().Unix()
}

func (p *Player) SetLoggedInMenuContext() {
	p.SetContext(actions.LoggedInMenuContext)
}

func (p *Player) IncrementCheatCounter() {
	if p.CheatCounter >= CheatCounterTolerance {
		p.GetTCPClient().Server.Kick(p.GetTCPClient())
	}

	p.CheatCounter += 1
}

func (p *Player) ResetCheatCounter() {
	p.CheatCounter = 0
}
