package interfaces

import (
	"kiv_ups_server/internal/net/tcp"
	"kiv_ups_server/internal/net/tcp/protocol"
)

type Action interface {
	Invoke(s MasterServer, p Player, message protocol.Message) (*tcp.ServerMessage, error)
}
