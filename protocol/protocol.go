package protocol

import (
	"fmt"
	"net"
)

var (
	ErrNotFoundProtocol = "protocol: [%s] can not be found!"
)

type DataPkg struct {
	From       string
	To         string
	Content    string
	FromSource string
	ToSource   string
}

type Protocol interface {
	SetConn(net.Conn)
	Read() (*DataPkg, error)
	Write(*DataPkg) error
	Close()
	Id() int32
}

type ProtocolMgr map[string]Protocol

var protocolMgr = make(ProtocolMgr, 1)

func Register(name string, p Protocol) {
	protocolMgr[name] = p
}

func initProtocol(proName string, conn net.Conn) (Protocol, error) {
	if p, ok := protocolMgr[proName]; ok {
		err := p.init(conn)
		if err != nil {
			return nil, err
		}
		return p, nil
	} else {
		return nil, fmt.Sprintf(ErrNotFoundProtocol, proName)
	}
}
