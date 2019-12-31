package nodes

import (
	tree2 "kiv_ups_server/internal/masterserver/gameserver/tree"
	interfaces2 "kiv_ups_server/internal/masterserver/interfaces"
	"kiv_ups_server/internal/net/tcp/protocol"
)

type RootNode struct {
	Node *tree2.Node `json:"-"`
}

func (rn RootNode) Init(node *tree2.Node) {
	rn.Node = node
}

func (rn RootNode) Process(messages []interfaces2.PlayerMessage, delta float64) {

}

func (rn RootNode) ListenMessages() []protocol.Message {
	return []protocol.Message{}
}

func (rn RootNode) Filter(messages []interfaces2.PlayerMessage) []interfaces2.PlayerMessage {
	return messages
}
