package nodes

import (
	"kiv_ups_server/masterserver/gameserver/tree"
	"kiv_ups_server/masterserver/interfaces"
	"kiv_ups_server/net/tcp/protocol"
	"math/rand"
)

type Spaceship struct {
	PosX           float64           `json:"pos_x"`
	PosY           float64           `json:"pos_y"`
	VelocityX      float64           `json:"velocity_x"`
	VelocityY      float64           `json:"velocity_y"`
	Rotation       float64           `json:"rotation"`
	PlayerName     string            `json:"player_name"`
	Immune         bool              `json:"immune"`
	ReloadPosition bool              `json:"reload_position"`
	Player         interfaces.Player `json:"-"`
	Score          *Score            `json:"-"`
	ImmuneTime     float64           `json:"-"`
	Radius         float64           `json:"-"`
	Node           *tree.Node        `json:"-"`
}

func (s *Spaceship) Init(node *tree.Node) {
	s.Node = node
	s.Radius = 50

	if s.Player != nil {
		s.PlayerName = s.Player.GetName()
	}
}

func (s *Spaceship) Process(playerMessages []interfaces.PlayerMessage, delta float64) {
	if s.ReloadPosition {
		s.ReloadPosition = false
	}

	for _, playerMessage := range playerMessages {
		message := playerMessage.GetMessage().Message
		if v, ok := message.(*protocol.PlayerMoveMessage); ok {
			s.PosX = v.PosX
			s.PosY = v.PosY
			s.VelocityX = v.VelocityX
			s.VelocityY = v.VelocityY
			s.Rotation = v.Rotation
			s.PlayerName = s.Player.GetName()
		} else if _, ok := message.(*protocol.ShootProjectileMessage); ok {
			velocity := Vector{0, -1}
			velocity = velocity.Rotated(s.Rotation).Normalized()

			node := tree.NewNode(s.Node.Parent, &Projectile{
				PosX:      s.PosX,
				PosY:      s.PosY,
				VelocityX: velocity.X * 600,
				VelocityY: velocity.Y * 600,
				Rotation:  s.Rotation,
				Player:    s.Player,
				Node:      nil,
			})
			node.Init()
			node.Value.Init(&node)

			s.Node.Parent.Children = append(s.Node.Parent.Children, &node)
		}
	}

	// check collision with asteroid
	if !s.Immune {
		for _, node := range s.Node.GetRoot().FindAllChildrenByType("asteroid") {
			asteroid := node.Value.(*Asteroid)
			asteroidCollider := Circle{
				asteroid.PosX,
				asteroid.PosY,
				asteroid.Radius,
			}

			if asteroidCollider.IsPointInside(s.PosX, s.PosY) {
				s.Die()
			}
		}
	}

	// remove immunity
	if s.Immune {
		s.ImmuneTime += delta
		if s.ImmuneTime > 3 {
			s.Immune = false
		}
	}
}

func (s *Spaceship) ListenMessages() []protocol.Message {
	return []protocol.Message{protocol.PlayerMoveMessage{}, protocol.ShootProjectileMessage{}}
}

func (s *Spaceship) Filter(playerMessages []interfaces.PlayerMessage) []interfaces.PlayerMessage {
	filteredPlayerMessages := make([]interfaces.PlayerMessage, 0)

	for _, playerMessage := range playerMessages {
		if _, ok := playerMessage.GetMessage().Message.(*protocol.PlayerMoveMessage); ok &&
			playerMessage.GetPlayer().GetName() != s.Player.GetName() {
			continue
		}

		if _, ok := playerMessage.GetMessage().Message.(*protocol.ShootProjectileMessage); ok &&
			playerMessage.GetPlayer().GetName() != s.Player.GetName() {
			continue
		}

		filteredPlayerMessages = append(filteredPlayerMessages, playerMessage)
	}

	return filteredPlayerMessages
}

func (s *Spaceship) RandomizePosition() {
	s.PosX = float64(rand.Intn(1920))
	s.PosY = float64(rand.Intn(1080))
	s.ReloadPosition = true
}

func (s *Spaceship) Immunity() {
	s.Immune = true
	s.ImmuneTime = 0
}

func (s *Spaceship) Die() {
	s.Immunity()
	s.RandomizePosition()
}
