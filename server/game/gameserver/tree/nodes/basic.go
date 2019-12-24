package nodes

import (
	"kiv_ups_server/game/interfaces"
	"kiv_ups_server/net/tcp/protocol"
	"time"
)

type RootNode struct {

}

func (s RootNode) Process(messages []interfaces.PlayerMessage, delta time.Duration) {

}

func (s RootNode) ListenMessages() []protocol.Message {
	return []protocol.Message{}
}

func (s RootNode) Filter(messages []interfaces.PlayerMessage) []interfaces.PlayerMessage {
	return messages
}
