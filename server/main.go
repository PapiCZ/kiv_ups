package main

import (
	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"kiv_ups_server/game"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	log.SetLevel(log.TraceLevel)
	log.SetFormatter(&prefixed.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.StampMilli,
	})
}

func main() {
	log.Println("Starting pprof at http://localhost:6060/debug/pprof/")
	go http.ListenAndServe("localhost:6060", nil)
	sockaddr := syscall.SockaddrInet4{
		Port: 35000,
		Addr: [4]byte{127, 0, 0, 1},
	}

	gameServer := game.NewServer(&sockaddr)


	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		err := gameServer.Stop()

		if err != nil {
			log.Errorln(err)
			os.Exit(1)
		}
		os.Exit(0)
	}()

	gameServer.Start()
}
