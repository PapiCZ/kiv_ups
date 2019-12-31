package tcp

import (
	"bufio"
	"io"
	"io/ioutil"
	protocol2 "kiv_ups_server/internal/net/tcp/protocol"
	"math/rand"
	"syscall"
)

// Unique identifier of client. This can be used as map key.
type UID int

// Client structure stores basic data about TCP client
type Client struct {
	TCP          *TCP
	UID          UID
	MessageChan  chan *protocol2.ProtoMessage
	Protocol     protocol2.GameProtocol
	failCounter  int
	decodeReader io.ReadCloser
	decodeWriter io.WriteCloser
}

// newClient initializes new client
func newClient(fd int, sockaddr syscall.Sockaddr, protocol protocol2.GameProtocol,
	messageChan chan *protocol2.ProtoMessage) Client {
	return Client{
		TCP: &TCP{
			FD:       fd,
			Sockaddr: sockaddr,
		},
		UID:         UID(rand.Int()),
		MessageChan: messageChan,
		Protocol:    protocol,
	}
}

// Send encodes given message and sends it to client
func (c Client) Send(message protocol2.ProtoMessage) (err error) {
	reader, writer := io.Pipe()

	go c.Protocol.Encode(message, writer)
	r := bufio.NewReader(reader)

	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return
	}

	_, err = syscall.Write(c.TCP.FD, bytes)

	return
}

// SendBytes allows to send binary data to client.
// You shouldn't need to do it, this function is primarily for TCP server.
func (c *Client) SendBytes(message []byte) (err error) {
	_, err = syscall.Write(c.TCP.FD, message)

	return
}

type ClientMessage struct {
	protocol2.Message
	protocol2.RequestId
	Sender            *Client

	// DisconnectRequest is used to notify master server about
	// client disconnection
	DisconnectRequest bool
}
