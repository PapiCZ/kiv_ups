package nodes

import (
	tree2 "kiv_ups_server/internal/masterserver/gameserver/tree"
	interfaces2 "kiv_ups_server/internal/masterserver/interfaces"
	"kiv_ups_server/internal/net/tcp/protocol"
)

type Score struct {
	Score      int                `json:"score"`
	PlayerName string             `json:"player_name"`
	Player     interfaces2.Player `json:"-"`
	Node       *tree2.Node        `json:"-"`
}

func (s *Score) Init(node *tree2.Node) {
	s.Node = node
	s.PlayerName = s.Player.GetName()
}

func (s *Score) Process(playerMessages []interfaces2.PlayerMessage, delta float64) {

}

func (s *Score) ListenMessages() []protocol.Message {
	return []protocol.Message{}
}

func (s *Score) Filter(playerMessages []interfaces2.PlayerMessage) []interfaces2.PlayerMessage {
	return playerMessages
}
