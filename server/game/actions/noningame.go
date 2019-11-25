package actions

import (
	"kiv_ups_server/game/interfaces"
	"kiv_ups_server/net/tcp"
	"kiv_ups_server/net/tcp/protocol"
)

func (a AuthenticateAction) Process(s interfaces.MasterServer, m interfaces.PlayerMessage) ActionResponse {
	authenticateData := m.GetMessage().Message.(*protocol.AuthenticateMessage)
	//server.Authenticate(ConvertShadowPlayerToPlayer(player, authenticateData.Name))

	return ActionResponse{
		ServerMessage: tcp.ServerMessage{
			Data:        authenticateData,
			Status:      true,
			Message:     "It works!",
		},
		Targets: []interfaces.Player{m.GetPlayer()},
	}
}
