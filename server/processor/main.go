package main

import (
	"github.com/kyf/wxkf/protocol"
	_ "github.com/kyf/wxkf/protocol/websocket"
	"net"

	mlog "github.com/kyf/util/log"
	"log"
)

var (
	SchedulerHost = "127.0.0.1:3001"
	Addr          = ":4001"

	LogPath   = "/var/log/wxkf/processor/"
	LogPrefix = "[wxkf-processor]"

	logger *mlog.Logger
)

const (
	HeartBeat = time.Second * 5
)

var (
	PROTOCOL = "websocket"
)

func main() {
	var err error
	logger, err = mlog.NewLogger(LogPath, LogPrefix, log.LstdFlags)
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	svr, err := NewServer()
	if err != nil {
		logger.Fatal(err)
	}

	exit := make(chan int, 1)
	go monitorShutdown(exit)

	go func() {
		err := svr.Run()
		if err != nil {
			logger.Fatal(err)
		}
	}()

	<-exit
	svr.Shutdown()
}
