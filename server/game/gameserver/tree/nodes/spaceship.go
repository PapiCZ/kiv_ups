package nodes

import (
	"kiv_ups_server/game/interfaces"
	"kiv_ups_server/net/tcp/protocol"
	"time"
)

type Spaceship struct {
	PosX       int               `json:"pos_x"`
	PosY       int               `json:"pos_y"`
	VelocityX  float64           `json:"velocity_x"`
	VelocityY  float64           `json:"velocity_y"`
	Rotation   float64           `json:"rotation"`
	PlayerName string            `json:"player_name"`
	Player     interfaces.Player `json:"-"`
}

func (s *Spaceship) Process(playerMessages []interfaces.PlayerMessage, delta time.Duration) {
	for _, playerMessage := range playerMessages {
		message := playerMessage.GetMessage().Message
		if v, ok := message.(*protocol.PlayerMoveMessage); ok {
			s.PosX = v.PosX
			s.PosY = v.PosY
			s.VelocityX = v.VelocityX
			s.VelocityY = v.VelocityY
			s.Rotation = v.Rotation
			s.PlayerName = s.Player.GetName()
		}
	}
}

func (s *Spaceship) ListenMessages() []protocol.Message {
	return []protocol.Message{protocol.PlayerMoveMessage{}}
}

func (s *Spaceship) Filter(playerMessages []interfaces.PlayerMessage) []interfaces.PlayerMessage {
	filteredPlayerMessages := make([]interfaces.PlayerMessage, 0)

	for _, playerMessage := range playerMessages {
		if _, ok := playerMessage.GetMessage().Message.(*protocol.PlayerMoveMessage); ok &&
			playerMessage.GetPlayer() != s.Player {
			continue
		}

		filteredPlayerMessages = append(filteredPlayerMessages, playerMessage)
	}

	return filteredPlayerMessages
}
