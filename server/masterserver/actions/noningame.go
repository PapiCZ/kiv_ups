package actions

import (
	log "github.com/sirupsen/logrus"
	"kiv_ups_server/masterserver/gameserver"
	"kiv_ups_server/masterserver/interfaces"
	"kiv_ups_server/net/tcp"
	"kiv_ups_server/net/tcp/protocol"
)

func (a KeepAliveAction) Process(s interfaces.MasterServer, m interfaces.PlayerMessage) ActionResponse {
	keepAliveData := m.GetMessage().Message.(*protocol.KeepAliveMessage)

	// Check if player sent valid data
	if keepAliveData.Ping == "pong" {
		m.GetPlayer().RefreshKeepAlive()
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
	connectedPlayer := m.GetPlayer()

	// Check if player with the same name is already connected
	for _, player := range s.GetPlayers() {
		if player.IsConnected() && player.GetName() == authenticateData.Name {
			return ActionResponse{
				ServerMessage: tcp.ServerMessage{
					Data:    authenticateData,
					Status:  false,
					Message: "the player with the specified name is already connected",
				},
				Targets: []interfaces.Player{m.GetPlayer()},
			}
		}
	}

	// Reconnect?
	for _, player := range s.GetPlayers() {
		if player.GetName() == authenticateData.Name && player.GetGameServer() != nil &&
			player.GetGameServer().IsRunning() {
			// Set new TCP client to old, disconnected player
			player.SetTCPClient(m.GetPlayer().GetTCPClient())
			connectedPlayer = player
			break
		}
	}

	// Authenticate player
	s.Authenticate(ConvertShadowPlayerToPlayer(connectedPlayer, authenticateData.Name))
	connectedPlayer.SetContext(LoggedInMenuContext)

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
	if m.GetPlayer().GetConnectedLobby() == nil || len(m.GetPlayer().GetConnectedLobby().GetPlayers()) == 0 {
		// Create new lobby and add it to server
		createLobbyData := m.GetMessage().Message.(*protocol.CreateLobbyMessage)
		lobby := interfaces.Lobby{
			Name:         createLobbyData.Name,
			Owner:        m.GetPlayer(),
			Players:      make(map[interfaces.PlayerUID]interfaces.Player),
			PlayersLimit: createLobbyData.PlayersLimit,
		}
		s.AddLobby(&lobby)

		// Add lobby owner to lobby
		lobby.AddPlayer(m.GetPlayer())

		log.Infof("Added lobby %s", createLobbyData.Name)

		// Move player to lobby context
		m.GetPlayer().SetContext(LobbyContext)
		m.GetPlayer().SetConnectedLobby(&lobby)
		return ActionResponse{
			ServerMessage: tcp.ServerMessage{
				Data:    &protocol.CreatedLobbyResponseMessage{},
				Status:  true,
				Message: "",
			},
			Targets: []interfaces.Player{m.GetPlayer()},
		}
	} else {
		return ActionResponse{
			ServerMessage: tcp.ServerMessage{
				Data:    &protocol.CreatedLobbyResponseMessage{},
				Status:  false,
				Message: "You can't create more than one lobby",
			},
			Targets: []interfaces.Player{m.GetPlayer()},
		}
	}
}

func (a DeleteLobbyAction) Process(s interfaces.MasterServer, m interfaces.PlayerMessage) ActionResponse {
	deleteLobbyData := m.GetMessage().Message.(*protocol.DeleteLobbyMessage)
	lobby, err := s.GetLobby(deleteLobbyData.Name)

	// Delete lobby by given name
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
	m.GetPlayer().SetContext(LoggedInMenuContext)
	m.GetPlayer().SetConnectedLobby(nil)

	return ActionResponse{
		ServerMessage: tcp.ServerMessage{
			Data:    &protocol.DeleteLobbyResponseMessage{},
			Status:  true,
			Message: "",
		},
		Targets: []interfaces.Player{m.GetPlayer()},
	}
}

func (a ListLobbiesAction) Process(s interfaces.MasterServer, m interfaces.PlayerMessage) ActionResponse {
	lobbies := s.GetLobbies()
	outputLobbies := make([]map[string]interface{}, 0, len(lobbies))

	for _, v := range lobbies {
		lobby := make(map[string]interface{})
		lobby["name"] = v.Name
		lobby["connected_players"] = len(v.Players)
		lobby["players_limit"] = v.PlayersLimit

		outputLobbies = append(outputLobbies, lobby)
	}

	return ActionResponse{
		ServerMessage: tcp.ServerMessage{
			Data:    &protocol.ListLobbiesResponseMessage{Lobbies: outputLobbies},
			Status:  true,
			Message: "",
		},
		Targets: []interfaces.Player{m.GetPlayer()},
	}
}

func (a JoinLobbyAction) Process(s interfaces.MasterServer, m interfaces.PlayerMessage) ActionResponse {
	joinLobbyData := m.GetMessage().Message.(*protocol.JoinLobbyMessage)
	lobby, err := s.GetLobby(joinLobbyData.Name)
	if err != nil || len(lobby.Players) >= lobby.PlayersLimit {
		// Send error if lobby doesn't exist or lobby contains maximum
		// number of players
		return ActionResponse{
			ServerMessage: tcp.ServerMessage{
				Data:    &protocol.JoinLobbyResponseMessage{},
				Status:  false,
				Message: "",
			},
			Targets: []interfaces.Player{m.GetPlayer()},
		}
	}

	// Notify all players about new player
	s.SendMessageWithoutRequest(tcp.ServerMessage{
		Data:    &protocol.LobbyPlayerConnectedMessage{Name: m.GetPlayer().GetName()},
		Status:  true,
		Message: "",
	}, lobby.GetPlayers()...)

	lobby.AddPlayer(m.GetPlayer())
	m.GetPlayer().SetContext(LobbyContext)
	m.GetPlayer().SetConnectedLobby(lobby)

	return ActionResponse{
		ServerMessage: tcp.ServerMessage{
			Data:    &protocol.JoinLobbyResponseMessage{},
			Status:  true,
			Message: "",
		},
		Targets: []interfaces.Player{m.GetPlayer()},
	}
}

func (a LeaveLobbyAction) Process(s interfaces.MasterServer, m interfaces.PlayerMessage) ActionResponse {
	playerLobby := m.GetPlayer().GetConnectedLobby()
	if playerLobby == nil {
		return ActionResponse{
			ServerMessage: tcp.ServerMessage{
				Data:    &protocol.LeaveLobbyResponseMessage{},
				Status:  false,
				Message: "",
			},
			Targets: []interfaces.Player{m.GetPlayer()},
		}
	}

	// Remove player from lobby and change its context
	playerLobby.RemovePlayer(m.GetPlayer())
	m.GetPlayer().SetConnectedLobby(nil)
	m.GetPlayer().SetContext(LoggedInMenuContext)

	if len(playerLobby.Players) == 0 {
		// Delete lobby if all players has disconnected
		s.DeleteLobby(playerLobby.Name)
	} else {
		// Notify all remaining player about disconnection
		s.SendMessageWithoutRequest(tcp.ServerMessage{
			Data:    &protocol.LobbyPlayerDisconnectedMessage{Name: m.GetPlayer().GetName()},
			Status:  true,
			Message: "",
		}, playerLobby.GetPlayers()...)
	}

	return ActionResponse{
		ServerMessage: tcp.ServerMessage{
			Data:    &protocol.LeaveLobbyResponseMessage{},
			Status:  true,
			Message: "",
		},
		Targets: []interfaces.Player{m.GetPlayer()},
	}
}

func (a ListLobbyPlayersAction) Process(s interfaces.MasterServer, m interfaces.PlayerMessage) ActionResponse {
	lobby := m.GetPlayer().GetConnectedLobby()

	if lobby == nil {
		return ActionResponse{
			ServerMessage: tcp.ServerMessage{
				Data:    &protocol.ListLobbyPlayersResponseMessage{},
				Status:  false,
				Message: "",
			},
			Targets: []interfaces.Player{m.GetPlayer()},
		}
	}

	// Convert all connected players to slice
	playerNames := make([]string, 0)
	for _, player := range lobby.GetPlayers() {
		playerNames = append(playerNames, player.GetName())
	}

	return ActionResponse{
		ServerMessage: tcp.ServerMessage{
			Data:    &protocol.ListLobbyPlayersResponseMessage{Players: playerNames},
			Status:  true,
			Message: "",
		},
		Targets: []interfaces.Player{m.GetPlayer()},
	}
}

func (a GameReconnectAvailableAction) Process(s interfaces.MasterServer, m interfaces.PlayerMessage) ActionResponse {
	// Check if player can reconnect to game
	return ActionResponse{
		ServerMessage: tcp.ServerMessage{
			Data: &protocol.GameReconnectAvailableResponseMessage{
				Available: m.GetPlayer().GetGameServer() != nil,
			},
			Status:  true,
			Message: "",
		},
		Targets: []interfaces.Player{m.GetPlayer()},
	}
}

func (a ReconnectAction) Process(s interfaces.MasterServer, m interfaces.PlayerMessage) ActionResponse {
	// Reconnect player to game
	gs := m.GetPlayer().GetGameServer()
	if gs != nil {
		gs.RemoveDisconnectedPlayer(m.GetPlayer())

		// notify players
		s.SendMessageWithoutRequest(tcp.ServerMessage{
			Status:  true,
			Message: "",
			Data:    &protocol.PlayerConnectedMessage{PlayerName: m.GetPlayer().GetName()},
		}, gs.GetPlayers()...)

		return ActionResponse{
			ServerMessage: tcp.ServerMessage{
				Data:    &protocol.ReconnectMessage{},
				Status:  true,
				Message: "",
			},
			Targets: []interfaces.Player{m.GetPlayer()},
		}
	} else {
		return ActionResponse{
			ServerMessage: tcp.ServerMessage{
				Data:    &protocol.ReconnectMessage{},
				Status:  false,
				Message: "you can't reconnect to the game",
			},
			Targets: []interfaces.Player{m.GetPlayer()},
		}
	}
}

func (a LeaveGameAction) Process(s interfaces.MasterServer, m interfaces.PlayerMessage) ActionResponse {
	m.GetPlayer().LeaveGame()
	m.GetPlayer().SetConnectedLobby(nil)
	// Remove player from game
	return ActionResponse{
		ServerMessage: tcp.ServerMessage{
			Data: &protocol.GameReconnectAvailableResponseMessage{
				Available: m.GetPlayer().GetGameServer() != nil,
			},
			Status:  true,
			Message: "",
		},
		Targets: []interfaces.Player{m.GetPlayer()},
	}
}

func (a StartGameAction) Process(s interfaces.MasterServer, m interfaces.PlayerMessage) ActionResponse {
	gs := gameserver.NewGameServer()

	// Add all players to game server and start the game
	for _, player := range m.GetPlayer().GetConnectedLobby().GetPlayers() {
		gs.AddPlayer(player)
		player.SetGameServer(&gs)
	}
	go gs.Start()
	s.DeleteLobby(m.GetPlayer().GetConnectedLobby().Name)

	return ActionResponse{
		ServerMessage: tcp.ServerMessage{
			Data:    &protocol.StartGameResponseMessage{},
			Status:  true,
			Message: "",
		},
		Targets: gs.GetPlayers(),
	}
}
