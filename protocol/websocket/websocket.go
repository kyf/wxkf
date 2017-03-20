package websocket

import (
	"bufio"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/kyf/wxkf/protocol"
	"io"
	"net"
	"net/http"
	"sync"
)

const (
	HeartBeat     = time.Second * 10
	HeartBeatPing = HeartBeat * 9 / 10
)

type Proto struct {
	conn websocket.Conn
	id   int32
	user []byte
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

	ReadUserTimeout := time.Second * 10
	this.conn.SetReadDeadline(time.Now().Add(ReadUserTimeout))
	_, err = io.ReadFull(this.conn, this.user)

	//user auth

	this.conn.SetPingHandler(this.PingHandler)
	go this.HeartBeat()
}

func (this *Proto) HeartBeat() {
	ticker := time.NewTicker(HeartBeatPing)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			err := this.conn.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				break
			}
		}
	}
}

func (this *Proto) Read() (*protocol.DataPkg, error) {
	var result DataPkg
begin:
	mtype, content, err := this.conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	switch mtype {
	case websocket.TextMessage:
		err = json.Unmarshal(content, &result)
		if err != nil {
			return nil, err
		}
	case websocket.BinaryMessage:
	case websocket.CloseMessage:
		return websocket.ErrCloseSent
	case websocket.PingMessage:
		goto begin
	case websocket.PongMessage:
		this.conn.SetPongHandler(this.PongHandler)
		goto begin
	}

	return &result, nil
}

func (this *Proto) PingHandler(appdata string) error {
	return this.conn.WriteMessage(websocket.PongMessage, nil)
}

func (this *Proto) PongHandler(appdata string) error {
	this.conn.SetReadDeadline(time.Now().Add(HeartBeat))
}

func (this *Proto) Write(dp *DataPkg) error {}

func (this *Proto) Id() int32 {
	return this.id
}

func (this *Proto) User() []byte {
	return this.user
}

func (this *Proto) Close() {
	this.conn.Close()
}

var id int32

func init() {
	p := &Proto{id: atomic.AddInt32(id, 1), user: make([]byte, 24)}

	protocol.Register("websocket", p)
}
