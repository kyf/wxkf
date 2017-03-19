package main

import (
	"net"

	"github.com/kyf/wxkf/server/dtype"
)

type ProClient struct {
	conn net.Conn
	IP   string
	Name string
}

func NewProClient(host string) (*ProClient, error) {
	conn, err := net.Dial("tcp", host)
	if err != nil {
		return nil, err
	}

	return &ProClient{conn: conn}
}

func (this *ProClient) Run(svr *Server) {
	ticker := time.NewTicker(HeartBeat)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			data := dtype.StatusPkg{IP: svr.IP, Name: svr.Name, ConnNum: svr.Count()}
			err := data.Encode(this.conn)
			if err != nil {
				logger.Errorf("proclient:  send status data  err :%v", err)
			}

		}
	}
}
