package game

import (
	"kiv_ups_server/game/interfaces"
	"kiv_ups_server/net/tcp"
	"math/rand"
)

type Player struct {
	tcpClient *tcp.Client
	UID       interfaces.PlayerUID
	Name      string
	Context   interfaces.PlayerContext
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
