package main

import (
	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"kiv_ups_server/masterserver"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strconv"
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
	go http.ListenAndServe("0.0.0.0:6060", nil)

	if len(os.Args) < 3 {
		log.Panicln("You need to pass host and port!")
	}

	host := net.ParseIP(os.Args[1])
	port, err := strconv.Atoi(os.Args[2])

	if err != nil {
		log.Panicln(err)
	}

	sockaddr := syscall.SockaddrInet4{
		Port: port,
		Addr: [...]byte{host[12], host[13], host[14], host[15]},
	}

	gameServer := masterserver.NewServer(&sockaddr)

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
