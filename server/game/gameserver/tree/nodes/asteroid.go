package nodes

import (
	"kiv_ups_server/game/gameserver/tree"
	"kiv_ups_server/game/interfaces"
	"kiv_ups_server/net/tcp/protocol"
	"math/rand"
)

const Width = 1920
const Height = 1080
const Margin = 100
const BorderWidth = 100

type AsteroidWrapper struct {
	Node *tree.Node `json:"-"`
}

func (a *AsteroidWrapper) Init(node *tree.Node) {
	a.Node = node
}

func (a *AsteroidWrapper) Process(playerMessages []interfaces.PlayerMessage, delta float64) {
	if len(a.Node.Children) < 50 {
		// create new asteroid
		// randomize pos x
		var minX, maxX, minY, maxY int
		side := rand.Int() % 4
		if side == 0 {
			// up
			minX = 0 - Margin - BorderWidth
			maxX = Width + Margin + BorderWidth
			minY = 0 - Margin - BorderWidth
			maxY = 0 - Margin
		} else if side == 1 {
			// right
			minX = Width + Margin
			maxX = Width + Margin + BorderWidth
			minY = 0 - Margin - BorderWidth
			maxY = Height + Margin + BorderWidth
		} else if side == 2 {
			// bottom
			minX = 0 - Margin - BorderWidth
			maxX = Width + Margin + BorderWidth
			minY = Height + Margin
			maxY = Height + Margin + BorderWidth
		} else if side == 3 {
			// left
			minX = 0 - Margin - BorderWidth
			maxX = 0 - Margin
			minY = 0 - Margin - BorderWidth
			maxY = Width + Margin + BorderWidth
		}

		posX := rand.Intn(maxX-minX) + minX
		posY := rand.Intn(maxY-minY) + minY
		node := tree.NewNode(a.Node, &Asteroid{
			PosX:      float64(posX),
			PosY:      float64(posY),
			VelocityX: -100 + rand.Float64()*(100+100),
			VelocityY: -100 + rand.Float64()*(100+100),
			Rotation:  0,
			Scale:     rand.Float64(),
			Node:      nil,
		})

		node.Init()
		node.Value.Init(&node)

		a.Node.Children = append(a.Node.Children, &node)
	}
}

func (a *AsteroidWrapper) ListenMessages() []protocol.Message {
	return []protocol.Message{}
}

func (a *AsteroidWrapper) Filter(playerMessages []interfaces.PlayerMessage) []interfaces.PlayerMessage {
	return playerMessages
}

type Asteroid struct {
	PosX      float64    `json:"pos_x"`
	PosY      float64    `json:"pos_y"`
	VelocityX float64    `json:"velocity_x"`
	VelocityY float64    `json:"velocity_y"`
	Rotation  float64    `json:"rotation"`
	Scale     float64    `json:"scale"`
	Value     int        `json:"value"`
	Node      *tree.Node `json:"-"`
}

func (a *Asteroid) Init(node *tree.Node) {
	a.Node = node
	a.Value = int(100 * a.Scale) * 30
}

func (a *Asteroid) Process(playerMessages []interfaces.PlayerMessage, delta float64) {
	a.PosX += a.VelocityX * delta
	a.PosY += a.VelocityY * delta

	if a.PosX < -Margin-BorderWidth || a.PosX > Width+Margin+BorderWidth ||
		a.PosY < -Margin-BorderWidth || a.PosY > Height+Margin+BorderWidth {
		a.Node.Destroy()
		return
	}
}

func (a *Asteroid) ListenMessages() []protocol.Message {
	return []protocol.Message{}
}

func (a *Asteroid) Filter(playerMessages []interfaces.PlayerMessage) []interfaces.PlayerMessage {
	return playerMessages
}
