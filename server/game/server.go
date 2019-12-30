package game

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"kiv_ups_server/game/actions"
	"kiv_ups_server/game/interfaces"
	"kiv_ups_server/net/tcp"
	"kiv_ups_server/net/tcp/protocol"
	"strconv"
	"syscall"
)

type Server struct {
	TCPServer        *tcp.Server
	Players          map[tcp.UID]interfaces.Player
	ActionDefinition actions.ActionDefinition
	Lobbies          map[string]*interfaces.Lobby
	GameServers      []interfaces.GameServer
}

func NewServer(sockaddr syscall.Sockaddr) (ms Server) {
	server, err := tcp.NewServer(sockaddr)

	if err != nil {
		log.Panic(err)
	}

	ms = Server{
		TCPServer:        server,
		Players:          make(map[tcp.UID]interfaces.Player),
		ActionDefinition: actions.NewDefinition(),
		Lobbies:          make(map[string]*interfaces.Lobby),
		GameServers:      make([]interfaces.GameServer, 0),
	}

	actions.RegisterAllActions(&ms.ActionDefinition)

	return
}

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

func (s *Server) RunAction(message tcp.ClientMessage) (err error) {
	p, ok := s.Players[message.Sender.UID]
	var player interfaces.Player
	if ok {
		player = p
	} else {
		pl := NewShadowPlayer(message.Sender, "", interfaces.PlayerContext(0))
		player = &pl
	}

	if message.DisconnectRequest {
		s.OnPlayerDisconnected(player)

		return nil
	}

	if message.GetTypeId() >= 300 && message.GetTypeId() <= 500 {
		gameServer := player.GetGameServer()
		if gameServer != nil {
			gameServer.GetRequestMessageChan() <- &PlayerMessage{
				ClientMessage: &message,
				Player:        player,
			}
		} else {
			// TODO: error
		}

		return nil
	}
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

	s.SendMessage(sm, message.RequestId, actionResponse.Targets...)

	return
}

func (s *Server) SendMessageWithoutRequest(sm tcp.ServerMessage, player ...interfaces.Player) {
	s.SendMessage(sm, "", player...)
}

func (s *Server) SendMessage(sm tcp.ServerMessage, requestId protocol.RequestId, player ...interfaces.Player) {
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

func (s *Server) OnPlayerDisconnected(player interfaces.Player) {
	if player.GetConnectedLobby() != nil && player.GetConnectedLobby().Owner == player {
		s.DeleteLobby(player.GetConnectedLobby().Name)
	}

	// TODO: Kick players from lobby
}

func (s *Server) Stop() (err error) {
	return s.TCPServer.Close()
}

func (s *Server) GetTCPServer() *tcp.Server {
	return s.TCPServer
}

func (s *Server) GetPlayers() map[tcp.UID]interfaces.Player {
	return s.Players
}

func (s *Server) Authenticate(player interfaces.Player) {
	s.Players[player.GetTCPClient().UID] = player

	log.Infoln("Authenticated player", player.GetName())
}

func (s *Server) AddLobby(lobby *interfaces.Lobby) {
	s.Lobbies[lobby.Name] = lobby
}

func (s *Server) DeleteLobby(name string) {
	delete(s.Lobbies, name)
}

func (s *Server) GetLobby(name string) (*interfaces.Lobby, error) {
	lobby, ok := s.Lobbies[name]

	if ok {
		return lobby, nil
	}

	return &interfaces.Lobby{}, errors.New("unknown lobby")
}

func (s *Server) GetLobbies() []*interfaces.Lobby {
	v := make([]*interfaces.Lobby, 0, len(s.Lobbies))

	for _, value := range s.Lobbies {
		v = append(v, value)
	}

	return v
}

func (s *Server) AddGameServer(gs interfaces.GameServer) {
	s.GameServers = append(s.GameServers, gs)
}

type PlayerMessage struct {
	*tcp.ClientMessage
	interfaces.Player
}

func (p *PlayerMessage) GetMessage() *tcp.ClientMessage {
	return p.ClientMessage
}

func (p PlayerMessage) GetPlayer() interfaces.Player {
	return p.Player
}
