package tcp

import (
	log "github.com/sirupsen/logrus"
	"syscall"
)

type TCP struct {
	FD       int
	Sockaddr syscall.Sockaddr
}

func (t TCP) SendBytes(message []byte) (err error) {
	_, err = syscall.Write(t.FD, message)

	return
}

func (t TCP) Close() (err error) {
	t.SendBytes([]byte("Bye!"))
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
