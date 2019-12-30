package tcp

import (
	"bufio"
	"io"
	"io/ioutil"
	"kiv_ups_server/net/tcp/protocol"
	"math/rand"
	"syscall"
)

type UID int

type Client struct {
	TCP          *TCP
	UID          UID
	MessageChan  chan *protocol.ProtoMessage
	Protocol     protocol.GameProtocol
	failCounter  int
	decodeReader io.ReadCloser
	decodeWriter io.WriteCloser
}

func newClient(fd int, sockaddr syscall.Sockaddr, protocol protocol.GameProtocol,
	messageChan chan *protocol.ProtoMessage) Client {
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

func (c Client) Send(message protocol.ProtoMessage) (err error) {
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

func (c *Client) SendBytes(message []byte) (err error) {
	_, err = syscall.Write(c.TCP.FD, message)

	return
}

type ClientMessage struct {
	protocol.Message
	protocol.RequestId
	Sender            *Client
	DisconnectRequest bool
}
