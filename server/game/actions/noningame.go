package actions

import (
	log "github.com/sirupsen/logrus"
	"kiv_ups_server/game/interfaces"
	"kiv_ups_server/net/tcp"
	"kiv_ups_server/net/tcp/protocol"
)

func (a KeepAliveAction) Process(s interfaces.MasterServer, m interfaces.PlayerMessage) ActionResponse {
	keepAliveData := m.GetMessage().Message.(*protocol.KeepAliveMessage)

	if keepAliveData.Ping == "pong" {
		return ActionResponse{
			ServerMessage: tcp.ServerMessage{
				Data:    &protocol.KeepAliveMessage{Ping: "ping-pong"},
				Status:  true,
				Message: "",
			},
			Targets: []interfaces.Player{m.GetPlayer()},
		}
	} else {
		return ActionResponse{
			ServerMessage: tcp.ServerMessage{
				Data:    nil,
				Status:  false,
				Message: "",
			},
			Targets: []interfaces.Player{m.GetPlayer()},
		}
	}
}

func (a AuthenticateAction) Process(s interfaces.MasterServer, m interfaces.PlayerMessage) ActionResponse {
	authenticateData := m.GetMessage().Message.(*protocol.AuthenticateMessage)
	s.Authenticate(ConvertShadowPlayerToPlayer(m.GetPlayer(), authenticateData.Name))
	m.GetPlayer().SetContext(LoggedInMenuContext)

	return ActionResponse{
		ServerMessage: tcp.ServerMessage{
			Data:    authenticateData,
			Status:  true,
			Message: "",
		},
		Targets: []interfaces.Player{m.GetPlayer()},
	}
}

func (a CreateLobbyAction) Process(s interfaces.MasterServer, m interfaces.PlayerMessage) ActionResponse {
	createLobbyData := m.GetMessage().Message.(*protocol.CreateLobbyMessage)
	s.AddLobby(interfaces.Lobby{
		Name:    createLobbyData.Name,
		Owner:   m.GetPlayer(),
		Players: make(map[interfaces.PlayerUID]interfaces.Player),
	})

	log.Infof("Added lobby %s", createLobbyData.Name)

	return ActionResponse{
		ServerMessage: tcp.ServerMessage{
			Data:    &protocol.CreatedLobbyResponseMessage{},
			Status:  true,
			Message: "",
		},
		Targets: []interfaces.Player{m.GetPlayer()},
	}
}

func (a DeleteLobbyAction) Process(s interfaces.MasterServer, m interfaces.PlayerMessage) ActionResponse {
	deleteLobbyData := m.GetMessage().Message.(*protocol.DeleteLobbyMessage)
	lobby, err := s.GetLobby(deleteLobbyData.Name)

	if err != nil || lobby.Owner != m.GetPlayer() {
		return ActionResponse{
			ServerMessage: tcp.ServerMessage{
				Data:    &protocol.DeleteLobbyResponseMessage{},
				Status:  false,
				Message: "You can't delete this lobby!",
			},
			Targets: []interfaces.Player{m.GetPlayer()},
		}
	}

	log.Infof("Deleted lobby %s", lobby.Name)

	return ActionResponse{
		ServerMessage: tcp.ServerMessage{
			Data:    &protocol.DeleteLobbyResponseMessage{},
			Status:  true,
			Message: "",
		},
		Targets: []interfaces.Player{m.GetPlayer()},
	}
}
