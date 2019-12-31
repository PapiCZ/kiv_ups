package tcp

import (
	log "github.com/sirupsen/logrus"
	"syscall"
)

// TCP is elementary structure that handles basic data about TCP connection
type TCP struct {
	FD       int
	Sockaddr syscall.Sockaddr
}

// SendBytes writes slice of bytes to file descriptor
func (t TCP) SendBytes(message []byte) (err error) {
	_, err = syscall.Write(t.FD, message)

	return
}

// Close shutdowns and then closes file descriptor
func (t TCP) Close() (err error) {
	err = syscall.Shutdown(t.FD, syscall.SHUT_RDWR)
	if err != nil {
		log.Errorln(err)
	}

	err = syscall.Close(t.FD)
	if err != nil {
		log.Errorln(err)
	}

	return
}
