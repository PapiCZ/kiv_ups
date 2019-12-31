package masterserver

import (
	"errors"
	log "github.com/sirupsen/logrus"
	actions2 "kiv_ups_server/internal/masterserver/actions"
	interfaces2 "kiv_ups_server/internal/masterserver/interfaces"
	"kiv_ups_server/internal/net/tcp"
	"kiv_ups_server/internal/net/tcp/protocol"
	"strconv"
	"syscall"
)

const (
	GameServerTypeMin = 400
	GameServerTypeMax = 599
)

// Server structure is a structure for master server that handles connected
// players, game servers and data for routing incoming messages
type Server struct {
	TCPServer        *tcp.Server
	Players          map[tcp.UID]interfaces2.Player
	ActionDefinition actions2.ActionDefinition
	Lobbies          map[string]*interfaces2.Lobby
	GameServers      []interfaces2.GameServer
}

// NewServer creates and initializes master server
func NewServer(sockaddr syscall.Sockaddr) (ms Server) {
	server, err := tcp.NewServer(sockaddr)

	if err != nil {
		log.Panic(err)
	}

	ms = Server{
		TCPServer:        server,
		Players:          make(map[tcp.UID]interfaces2.Player),
		ActionDefinition: actions2.NewDefinition(),
		Lobbies:          make(map[string]*interfaces2.Lobby),
		GameServers:      make([]interfaces2.GameServer, 0),
	}

	actions2.RegisterAllActions(&ms.ActionDefinition)

	return
}

// Start starts master server and makes it ready for incoming messages
func (s *Server) Start() (err error) {
	clientMessageChan := make(chan tcp.ClientMessage)

	go s.TCPServer.Start(clientMessageChan)

	for {
		message := <-clientMessageChan
		err = s.RunAction(message)
		if err != nil {
			log.Errorln(err)
		}
	}
}

// RunAction routes incoming messages to its actions and returns response to clients.
func (s *Server) RunAction(message tcp.ClientMessage) (err error) {
	p, ok := s.Players[message.Sender.UID]
	var player interfaces2.Player
	if ok {
		player = p
	} else {
		// Client isn't authenticated. We need to create ShadowPlayer.
		pl := NewShadowPlayer(message.Sender, "", interfaces2.PlayerContext(0))
		player = &pl
	}

	if message.DisconnectRequest {
		// Disconnect player
		s.OnPlayerDisconnected(player)

		return nil
	}

	if message.GetTypeId() >= GameServerTypeMin && message.GetTypeId() <= GameServerTypeMax {
		gameServer := player.GetGameServer()
		if gameServer != nil {
			// Forward message to game server
			gameServer.GetRequestMessageChan() <- &PlayerMessage{
				ClientMessage: &message,
				Player:        player,
			}
		}

		return nil
	}

	// Get action according to message type and player's context
	action := s.ActionDefinition.GetAction(message.Message.GetTypeId(), player.GetContext())

	if action == nil {
		_ = message.Sender.Send(protocol.ProtoMessage{
			Message: &tcp.ServerMessage{
				Status:  false,
				Message: "Action not found",
				Data:    &protocol.ActionErrorMessage{},
			},
			RequestId: message.RequestId,
		})

		return errors.New("invalid action: " + strconv.Itoa(int(message.Message.GetTypeId())) +
			", player context: " + strconv.Itoa(int(player.GetContext())))
	}

	actionResponse := action.Process(s, &PlayerMessage{
		ClientMessage: &message,
		Player:        player,
	})
	sm := actionResponse.ServerMessage

	// Send response to all targets
	s.SendMessage(sm, message.RequestId, actionResponse.Targets...)

	return
}

// SendMessageWithoutRequest allows to send message from server without request ID
func (s *Server) SendMessageWithoutRequest(sm tcp.ServerMessage, player ...interfaces2.Player) {
	s.SendMessage(sm, "", player...)
}

// SendMessage sends ServerMessage to all given players
func (s *Server) SendMessage(sm tcp.ServerMessage, requestId protocol.RequestId, player ...interfaces2.Player) {
	for _, p := range player {
		log.Tracef("Server answers to client %d: %#v | Data: %#v", p.GetUID(), sm, sm.Data)
		err := p.GetTCPClient().Send(protocol.ProtoMessage{
			Message:   sm,
			RequestId: requestId,
		})

		if err != nil {
			log.Errorln(err)
		}
	}
}

func (s *Server) OnPlayerDisconnected(player interfaces2.Player) {
	lobby := player.GetConnectedLobby()
	if lobby != nil {
		lobby.RemovePlayer(player)

		// Notify all players in lobby about disconnection
		s.SendMessageWithoutRequest(tcp.ServerMessage{
			Data:    &protocol.LobbyPlayerDisconnectedMessage{Name: player.GetName()},
			Status:  true,
			Message: "",
		}, lobby.GetPlayers()...)

		if len(lobby.Players) == 0 {
			s.DeleteLobby(lobby.Name)
		}
	}
}

// Stop stops mater server
func (s *Server) Stop() (err error) {
	return s.TCPServer.Close()
}

// GetTCPServer is getter for tcp server
func (s *Server) GetTCPServer() *tcp.Server {
	return s.TCPServer
}

// GetPlayers is getter for players
func (s *Server) GetPlayers() map[tcp.UID]interfaces2.Player {
	return s.Players
}

// Authenticate authenticates given player
func (s *Server) Authenticate(player interfaces2.Player) {
	s.Players[player.GetTCPClient().UID] = player

	log.Infoln("Authenticated player", player.GetName())
}

// AddLobby adds lobby to master server
func (s *Server) AddLobby(lobby *interfaces2.Lobby) {
	s.Lobbies[lobby.Name] = lobby
}

// DeleteLobby deletes lobby by its name
func (s *Server) DeleteLobby(name string) {
	delete(s.Lobbies, name)
}

// GetLobby returns lobby by its name.
func (s *Server) GetLobby(name string) (*interfaces2.Lobby, error) {
	lobby, ok := s.Lobbies[name]

	if ok {
		return lobby, nil
	}

	return nil, errors.New("unknown lobby")
}

// GetLobbies is getter for all lobbies
func (s *Server) GetLobbies() []*interfaces2.Lobby {
	v := make([]*interfaces2.Lobby, 0, len(s.Lobbies))

	for _, value := range s.Lobbies {
		v = append(v, value)
	}

	return v
}

// AddGameServer adds game server to master server
func (s *Server) AddGameServer(gs interfaces2.GameServer) {
	s.GameServers = append(s.GameServers, gs)
}

type PlayerMessage struct {
	*tcp.ClientMessage
	interfaces2.Player
}

func (p *PlayerMessage) GetMessage() *tcp.ClientMessage {
	return p.ClientMessage
}

func (p PlayerMessage) GetPlayer() interfaces2.Player {
	return p.Player
}
