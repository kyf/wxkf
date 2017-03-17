package main

import (
	"context"
	"log"
	"net/http"
	"time"

	mlog "github.com/kyf/util/log"
)

var (
	HttpAddr = ":2001"
	TcpAddr  = ":3001"

	LogPath   = ""
	LogPrefix = ""

	logger *mlog.Logger
)

const (
	ProcessorHeartBeat = time.Duration(10)
)

func newHttpServer() *http.Server {
	svr := &http.Server{
		Addr:         HttpAddr,
		Handler:      &Handler{},
		ReadTimeout:  time.Minute * 1,
		WriteTimeout: time.Second * 30,
	}
	svr.SetKeepAlivesEnabled(false)
	return svr
}

func main() {
	var err error
	logger, err = mlog.NewLogger(LogPath, LogPrefix, log.LstdFlags)
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	httpSvr := newHttpServer()
	tcpSvr, err := NewTcpServer(TcpAddr)
	if err != nil {
		log.Fatal(err)
	}

	exit := make(chan int, 1)

	go monitorShutdown(exit)
	go func() {
		err := httpSvr.ListenAndServe()
		if err != nil {
			logger.Fatal(err)
		}
	}()
	go func() {
		err := tcpSvr.Run()
		if err != nil {
			logger.Fatal(err)
		}
	}()

	<-exit

	httpSvr.Shutdown(context.Background())
	tcpSvr.Shutdown()
}
