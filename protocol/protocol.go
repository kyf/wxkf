package protocol

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
)

var (
	ErrNotFoundProtocol = "protocol: [%s] can not be found!"

	ErrNil = errors.New("DataPkg body is nil")
)

type DataPkg struct {
	From       string
	To         string
	Content    string
	FromSource string
	ToSource   string
}

func (this *DataPkg) Encode(w io.Writer) error {
	body, err := json.Marshal(this)
	if err != nil {
		return err
	}

	size := int64(len(body))
	err = binary.Write(w, binary.BigEndian, size)
	if err != nil {
		return err
	}

	_, err = w.Write(body)
	if err != nil {
		return err
	}

	return nil
}

func (this *DataPkg) Decode(r *io.Reader) error {
	var size int64
	err := binary.Read(r, binary.BigEndian, size)
	if err != nil {
		return err
	}

	if size == 0 {
		return ErrNil
	}

	body := make([]byte, size)
	err = io.ReadFull(r, body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(this, body)
	if err != nil {
		return err
	}

	return nil
}

type Protocol interface {
	SetConn(net.Conn)
	Read() (*DataPkg, error)
	Write(*DataPkg) error
	Close()
	Id() int32
	User() []byte
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
