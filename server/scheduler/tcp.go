package main

import (
	"bytes"
	"encoding/binary"
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
		_, err := io.ReadFull(conn, headbuf)
		if err != nil {
			logger.Errorf("[scheduler ]tcp read header err:%v", err)
			break
		}

		var pkglen int64
		binary.Read(bytes.NewReader(headbuf), binary.BigEndian, pkglen)
		pack := make([]byte, pkglen)
		_, err = io.ReadFull(conn, pack)
		if err != nil {
			logger.Errorf("[scheduler ]tcp read data err:%v", err)
			break
		}

		err = decodeProcessPkg(pack, pro)
		if err != nil {
			logger.Errorf("[scheduler ]tcp decode data err:%v", err)
			break
		}

	}

	tcpconn.Close()
	defaultProMgr.Unregister(pro)
}

func (this *TcpServer) Shutdown() {
	this.ln.Close()
}
