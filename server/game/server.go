package game

import (
	log "github.com/sirupsen/logrus"
	"kiv_ups_server/game/actions"
	"kiv_ups_server/game/interfaces"
	"kiv_ups_server/net/tcp"
	"kiv_ups_server/net/tcp/protocol"
	"syscall"
)

type Server struct {
	TCPServer        *tcp.Server
	Players          map[tcp.UID]interfaces.Player
	ActionDefinition actions.ActionDefinition
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
	}

	actions.RegisterAllActions(&ms.ActionDefinition)

	return
}

func (s *Server) Start() (err error) {
	clientMessageChan := make(chan tcp.ClientMessage)

	go s.TCPServer.Start(clientMessageChan)

	for {
		message := <-clientMessageChan
		s.RunAction(message)
	}
}

func (s *Server) RunAction(message tcp.ClientMessage) (err error) {
	p, ok := s.Players[message.Sender.UID]
	var player interfaces.Player
	if ok {
		player = p
	} else {
		pl := NewShadowPlayer(message.Sender, "", actions.DefaultContext)
		player = &pl
	}

	action := s.ActionDefinition.GetAction(message.Message.GetTypeId(), player.GetContext())
	// TODO: should return error
	actionResponse := action.Process(s, &PlayerMessage{
		ClientMessage: &message,
		Player:        player,
	})
	sm := actionResponse.ServerMessage

	for _, target := range actionResponse.Targets {
		log.Tracef("Server answers to client %d: %#v | Data: %#v", target.GetUID(), sm, sm.Data)
		_ = target.GetTCPClient().Send(protocol.ProtoMessage{
			Message:   sm,
			RequestId: message.RequestId,
		})
	}

	return
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
