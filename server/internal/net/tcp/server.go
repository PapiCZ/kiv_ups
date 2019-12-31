package tcp

import (
	log "github.com/sirupsen/logrus"
	"io"
	protocol2 "kiv_ups_server/internal/net/tcp/protocol"
	"syscall"
	"unsafe"
)

const (
	// DecodeErrThreshold means that client will be kicked out after consecutive
	// protocol errors
	DecodeErrThreshold = 10

	// FdBits is for FD_SET and FD_ISSET functions. Don't touch it!
	FdBits = int(unsafe.Sizeof(0) * 8)

	// BuffLen is length of TCP read buffer
	BuffLen = 1024
)

// Server is struct for TCP server that handles TCP connection,
// connected clients, and communication protocol
type Server struct {
	TCP      *TCP
	Clients  map[int]*Client
	Protocol protocol2.GameProtocol
}

// NewServer initializes TCP server
func NewServer(sockaddr syscall.Sockaddr) (server *Server, err error) {
	server = &Server{}

	// create socket
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		return
	}

	// set socket options
	err = syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
	if err != nil {
		return
	}

	// bind socket to address
	err = syscall.Bind(fd, sockaddr)
	if err != nil {
		return
	}

	// listen to socket
	err = syscall.Listen(fd, 5)
	if err != nil {
		return
	}

	err = nil

	// Initialize protocol and its definition
	def := protocol2.NewDefinition()
	protocol2.RegisterAllMessages(&def)
	proto := protocol2.GameProtocol{Def: def}

	// Initialize Server structure
	server = &Server{
		TCP: &TCP{
			FD:       fd,
			Sockaddr: sockaddr,
		},
		Clients:  make(map[int]*Client),
		Protocol: proto,
	}

	return
}

// Start starts TCP server that now accepts incoming connections
func (s *Server) Start(clientMessageChan chan ClientMessage) {
	log.Info("Accepting connections")

	// Perpare writer map, buffer map, and file descriptor set
	writers := make(map[UID]io.Writer)
	buff := make([]byte, BuffLen)
	rfds := &syscall.FdSet{}

	for {
		// zero file descriptor set
		FD_ZERO(rfds)
		// add server FD to FD set
		FD_SET(rfds, s.TCP.FD)
		maxFd := s.TCP.FD

		// Add all clients to FD set
		for _, client := range s.Clients {
			FD_SET(rfds, client.TCP.FD)

			if client.TCP.FD > maxFd {
				maxFd = client.TCP.FD
			}
		}

		// Wait for incoming message from one of file descriptors
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

			// Initialize new client and
			client := newClient(clientFd, sockaddr, s.Protocol, make(chan *protocol2.ProtoMessage))
			s.Clients[client.TCP.FD] = &client
			reader, writer := io.Pipe()
			if err != nil {
				log.Errorln(err)
				continue
			}
			writers[client.UID] = writer
			client.decodeReader = reader
			client.decodeWriter = writer

			// statusChan is used to receive statuses from decoder
			statusChan := make(chan bool)
			go s.Protocol.InfiniteDecode(reader, client.MessageChan, statusChan)

			// Starts goroutine that reads messages from client, and forwards it to game server
			go func(messageChan chan ClientMessage, c *Client) {
				for {
					msg, ok := <-c.MessageChan

					if !ok {
						log.Traceln("Closed ClientMessage channel...")
						return
					}

					// Forwards message to server
					messageChan <- ClientMessage{
						Message:   msg.Message,
						RequestId: msg.RequestId,
						Sender:    c,
					}
				}
			}(clientMessageChan, &client)

			// Starts goroutine that checks status of protocol decoder. After
			// DecodeErrThreshold of consecutive failures it kicks out client
			// from the server.
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

					// Kicks out client if client fail counter is greater than
					// DecodeErrThreshold
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
			// Checks if client want to talk or disconnect
			for _, client := range s.Clients {
				if FD_ISSET(rfds, client.TCP.FD) {
					// Read message from client
					n, err := syscall.Read(client.TCP.FD, buff)

					if err != nil {
						log.Errorln("Read error:", err)
						break
					}

					if n == 0 {
						s.Kick(clientMessageChan, client)
					} else {
						// Write message to client buffer
						_, _ = writers[client.UID].Write(buff[:n])
					}
				}
			}
		}
	}
}

// Close kicks out all clients and then shutdowns and closes TCP server
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

// Kick kicks given client from TCP server and set disconnect notification to master server
func (s *Server) Kick(clientMessageChan chan ClientMessage, client *Client) {
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
	Status  bool              `json:"status"`
	Message string            `json:"message"`
	Data    protocol2.Message `json:"data"`
}

func (sm ServerMessage) GetTypeId() protocol2.MessageType {
	return sm.Data.GetTypeId()
}
