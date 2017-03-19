package main

import (
	"bytes"
	"encoding/binary"
	"github.com/kyf/wxkf/server/dtype"
	"io"
	"net"
	"time"
)

type TcpServer struct {
	ln net.Listener
}

func NewTcpServer(TcpAddr string) (*TcpServer, error) {
	ln, err := net.Listen("tcp", TcpAddr)
	if err != nil {
		return nil, err
	}

	return &TcpServer{ln: ln}, nil
}

func (this *TcpServer) Run() error {
	for {
		conn, err := this.ln.Accept()
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
}

func handleConn(conn net.Conn) {
	tcpconn, _ := conn.(*net.TCPConn)
	tcpconn.SetKeepAlive(true)
	tcpconn.SetKeepAlivePeriod(time.Minute * 3)

	pro := &Process{Name: "", IP: "", ConnNum: 0}

	defaultProMgr.Register(pro)
	headbuf := make([]byte, 4)
	for {
		tcpconn.SetReadDeadline(time.Now().Add(time.Second * ProcessorHeartBeat))
		var statusData dtype.StatusPkg
		err := statusData.Decode(tcpconn)
		if err != nil {
			logger.Errorf("processor:  decode status pkg err:%v", err)
			continue
		}

	}

	tcpconn.Close()
	defaultProMgr.Unregister(pro)
}

func (this *TcpServer) Shutdown() {
	this.ln.Close()
}
