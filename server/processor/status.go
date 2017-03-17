package main

import (
	"net"
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

func (this *ProClient) Run() {
	ticker := time.NewTicker(HeartBeat)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			this.conn.Write()
		}
	}
}
