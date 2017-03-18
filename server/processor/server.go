package main

import (
	"github.com/kyf/wxkf/protocol"
	"net"
	"sync"
)

type Server struct {
	ln      net.Listener
	ConnMgr map[int32]protocol.Protocol
	sync.Mutex
}

func NewServer() (*Server, error) {
	ln, err := net.Listen("tcp", Addr)
	if err != nil {
		return nil, err
	}

	return &Server{ln: ln}, nil
}

func (this *Server) Run() error {
	for {
		conn, err := ln.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				time.Sleep(time.Second * 1)
				continue
			} else {
				return err
			}
		}

		go handleConn(conn)
	}

	return nil
}

func (this *Server) RegisterConn(pro Protocol) {
	this.Lock()
	defer this.Unlock()

	this.ConnMgr[pro.Id()] = pro
}

func (this *Server) UnregisterConn(pro Protocol) {
	this.Lock()
	defer this.Unlock()

	id := pro.Id()
	for k, v := range this.ConnMgr {
		if k == id {
			this.ConnMgr[id] = nil
			delete(this.ConnMgr, id)
		}
	}
}

func handleConn(conn net.Conn, svr *Server) {
	tcpconn, _ := conn.(*net.TCPConn)
	tcpconn.SetKeepAlive(true)
	tcpconn.SetKeepAlivePeriod(time.Minute * 3)
	defer tcpconn.Close()

	pro, err := protocol.initProtocol(PROTOCOL)
	if err != nil {
		logger.Errorf("process: initProtocol err:%v", err)
		return
	}
	svr.RegisterConn(pro)
	defer svr.UnregisterConn(pro)
	defer pro.Close()
	for {
		data, err := pro.Read()
		if err != nil {
			logger.Errorf("protocol read err:%v", err)
			break
		}
		sendMessage(data)
	}
}

func sendMessage(data *DataPkg) {
	switch data.ToSource {
	case SourceWX:

	default:
		isonline, islocal, target := getOnlineConn(data.To)
		if !isonline {
			sendOffline(data)
			continue
		}
		if islocal {

		}
	}

}
