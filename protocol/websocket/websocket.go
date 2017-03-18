package websocket

import (
	"bufio"
	"github.com/gorilla/websocket"
	"github.com/kyf/wxkf/protocol"
	"net"
	"net/http"
)

type Proto struct {
	conn websocket.Conn
}

func (this *Proto) init(conn net.Conn) {
	reader := bufio.NewReader(reader)
	req, err := http.ReadRequest(reader)
	if err != nil {
		logger.Errorf(err)
		conn.Close()
		return
	}
	w, err := http.ReadResponse(reader, req)
	if err != nil {
		logger.Errorf(err)
		conn.Close()
		return
	}

	c, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if err != nil {
		logger.Errorf(err)
		conn.Close()
		return
	}

	this.conn = c
}

func (this *Proto) Read() (*protocol.DataPkg, error) {
	mtype, content, err := this.conn.ReadMessage()
	if err != nil {
		return nil, err
	}

}

func (this *Proto) Write(dp *DataPkg) error {}

func (this *Proto) Close() {
	this.conn.Close()
}

func init() {
	p := &Proto{}

	protocol.Register("websocket", p)
}
