package nodes

import (
	"kiv_ups_server/game/gameserver/tree"
	"kiv_ups_server/game/interfaces"
	"kiv_ups_server/net/tcp/protocol"
)

type RootNode struct {
	Node *tree.Node `json:"-"`
}

func (rn RootNode) Init(node *tree.Node) {
	rn.Node = node
}

func (rn RootNode) Process(messages []interfaces.PlayerMessage, delta float64) {

}

func (rn RootNode) ListenMessages() []protocol.Message {
	return []protocol.Message{}
}

func (rn RootNode) Filter(messages []interfaces.PlayerMessage) []interfaces.PlayerMessage {
	return messages
}
