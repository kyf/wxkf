package main

import (
	"github.com/kyf/wxkf/protocol"
	"net"
)

type Server struct {
	ln net.Listener
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

func handleConn(conn net.Conn) {
	tcpconn, _ := conn.(*net.TCPConn)
	tcpconn.SetKeepAlive(true)
	tcpconn.SetKeepAlivePeriod(time.Minute * 3)

	pro, err := protocol.initProtocol(PROTOCOL)
	for {
		data, err := pro.Read()
		if err != nil {
			logger.Errorf("protocol read err:%v", err)
			break
		}
		switch data.ToSource {
		case "wx":

		}

	}
}
