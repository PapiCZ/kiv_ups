package nodes

import (
	"kiv_ups_server/game/gameserver/tree"
	"kiv_ups_server/game/interfaces"
	"kiv_ups_server/net/tcp/protocol"
)

type Spaceship struct {
	PosX       float64           `json:"pos_x"`
	PosY       float64           `json:"pos_y"`
	VelocityX  float64           `json:"velocity_x"`
	VelocityY  float64           `json:"velocity_y"`
	Rotation   float64           `json:"rotation"`
	PlayerName string            `json:"player_name"`
	Player     interfaces.Player `json:"-"`
	Score      *Score            `json:"-"`
	Node       *tree.Node        `json:"-"`
}

func (s *Spaceship) Init(node *tree.Node) {
	s.Node = node

	if s.Player != nil {
		s.PlayerName = s.Player.GetName()
	}
}

func (s *Spaceship) Process(playerMessages []interfaces.PlayerMessage, delta float64) {
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
