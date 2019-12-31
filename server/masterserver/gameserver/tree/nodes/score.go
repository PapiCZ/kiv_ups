package nodes

import (
	"kiv_ups_server/masterserver/gameserver/tree"
	"kiv_ups_server/masterserver/interfaces"
	"kiv_ups_server/net/tcp/protocol"
)

type Score struct {
	Score      int               `json:"score"`
	PlayerName string            `json:"player_name"`
	Player     interfaces.Player `json:"-"`
	Node       *tree.Node        `json:"-"`
}

func (s *Score) Init(node *tree.Node) {
	s.Node = node
	s.PlayerName = s.Player.GetName()
}

func (s *Score) Process(playerMessages []interfaces.PlayerMessage, delta float64) {

}

func (s *Score) ListenMessages() []protocol.Message {
	return []protocol.Message{}
}

func (s *Score) Filter(playerMessages []interfaces.PlayerMessage) []interfaces.PlayerMessage {
	return playerMessages
}
