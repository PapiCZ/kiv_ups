package tcp

import (
	log "github.com/sirupsen/logrus"
	"io"
	"kiv_ups_server/net/tcp/protocol"
	"syscall"
	"unsafe"
)

const DecodeErrThreshold = 10
const FdBits = int(unsafe.Sizeof(0) * 8)

const (
	BuffLen = 1024
)

type Server struct {
	TCP      *TCP
	Clients  map[int]*Client
	Protocol protocol.GameProtocol
}

func NewServer(sockaddr syscall.Sockaddr) (server *Server, err error) {
	server = &Server{}
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		return
	}

	err = syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
	if err != nil {
		return
	}

	err = syscall.Bind(fd, sockaddr)
	if err != nil {
		return
	}

	err = syscall.Listen(fd, 5)
	if err != nil {
		return
	}

	err = nil

	def := protocol.NewDefinition()
	protocol.RegisterAllMessages(&def)
	proto := protocol.GameProtocol{Def: def}

	s := Server{
		TCP: &TCP{
			FD:       fd,
			Sockaddr: sockaddr,
		},
		Clients:  make(map[int]*Client),
		Protocol: proto,
	}

	server = &s

	return
}

func (s *Server) Start(clientMessageChan chan ClientMessage) {
	log.Info("Accepting connections")

	writers := make(map[UID]io.Writer)
	buff := make([]byte, BuffLen)
	rfds := &syscall.FdSet{}

	for {
		FD_ZERO(rfds)
		FD_SET(rfds, s.TCP.FD)
		maxFd := s.TCP.FD

		for _, client := range s.Clients {
			FD_SET(rfds, client.TCP.FD)

			if client.TCP.FD > maxFd {
				maxFd = client.TCP.FD
			}
		}

		activeFd, err := syscall.Select(maxFd+1, rfds, nil, nil, nil)
		if err != nil {
			log.Errorln("Select error:", err)
			continue
		}
		if activeFd < 0 {
			log.Errorln("Negative activeFd")
			continue
		}

		if FD_ISSET(rfds, s.TCP.FD) {
			clientFd, sockaddr, err := syscall.Accept(s.TCP.FD)

			if err != nil {
				log.Errorln("Accept error:", err)
				continue
			}

			client := newClient(clientFd, sockaddr, s.Protocol, make(chan *protocol.ProtoMessage))
			s.Clients[client.TCP.FD] = &client
			reader, writer := io.Pipe()
			if err != nil {
				log.Errorln(err)
				continue
			}
			writers[client.UID] = writer
			client.decodeReader = reader
			client.decodeWriter = writer

			statusChan := make(chan bool)
			go s.Protocol.InfiniteDecode(reader, client.MessageChan, statusChan)

			go func(messageChan chan ClientMessage, c *Client) {
				for {
					msg, ok := <-c.MessageChan

					if !ok {
						log.Traceln("Closed ClientMessage channel...")
						return
					}

					messageChan <- ClientMessage{
						Message:   msg.Message,
						RequestId: msg.RequestId,
						Sender:    c,
					}
				}
			}(clientMessageChan, &client)

			go func(sc chan bool, clientMessageChan chan ClientMessage, c *Client) {
				for {
					status, ok := <-sc

					if !ok {
						return
					}

					if !status {
						c.failCounter++
					} else {
						c.failCounter = 0
					}

					if c.failCounter > DecodeErrThreshold {
						s.Kick(clientMessageChan, c)
						return
					}
				}
			}(statusChan, clientMessageChan, &client)

			log.Infof("New Client[FD %v]: %v:%v",
				client.TCP.FD,
				client.TCP.Sockaddr.(*syscall.SockaddrInet4).Addr,
				client.TCP.Sockaddr.(*syscall.SockaddrInet4).Port,
			)
		} else {
			for _, client := range s.Clients {
				if FD_ISSET(rfds, client.TCP.FD) {
					n, err := syscall.Read(client.TCP.FD, buff)

					if err != nil {
						log.Errorln("Read error:", err)
						break
					}

					if n == 0 {
						s.Kick(clientMessageChan, client)
					} else {
						// Client wanna talk
						_, _ = writers[client.UID].Write(buff[:n])
					}
				}
			}
		}
	}
}

func (s Server) Close() (err error) {
	log.Infoln("Shutting down...")

	for _, c := range s.Clients {
		_ = c.TCP.Close()
	}

	err = syscall.Shutdown(s.TCP.FD, syscall.SHUT_RDWR)
	if err != nil {
		log.Errorln(err)
	}

	err = syscall.Close(s.TCP.FD)

	if err != nil {
		log.Errorln(err)
	}

	return
}

func (s *Server) Kick(clientMessageChan chan ClientMessage, client *Client) {
	// Client wanna disconnect
	log.Infof("Client disconnected [FD %v]: %v:%v",
		client.TCP.FD,
		client.TCP.Sockaddr.(*syscall.SockaddrInet4).Addr,
		client.TCP.Sockaddr.(*syscall.SockaddrInet4).Port,
	)

	clientMessageChan <- ClientMessage{
		Message:           nil,
		RequestId:         "",
		DisconnectRequest: true,
		Sender:            client,
	}
	close(client.MessageChan)
	_ = client.decodeReader.Close()
	_ = client.decodeWriter.Close()
	_ = syscall.Close(client.TCP.FD)

	delete(s.Clients, client.TCP.FD)
}

func FD_SET(p *syscall.FdSet, fd int) {
	p.Bits[fd/FdBits] |= int64(uint(1) << (uint(fd) % uint(FdBits)))
}

func FD_ISSET(p *syscall.FdSet, fd int) bool {
	return (p.Bits[fd/FdBits] & int64(uint(1)<<(uint(fd)%uint(FdBits)))) != 0
}

func FD_ZERO(p *syscall.FdSet) {
	for i := range p.Bits {
		p.Bits[i] = 0
	}
}

type ServerMessage struct {
	Status  bool             `json:"status"`
	Message string           `json:"message"`
	Data    protocol.Message `json:"data"`
}

func (sm ServerMessage) GetTypeId() protocol.MessageType {
	return sm.Data.GetTypeId()
}
