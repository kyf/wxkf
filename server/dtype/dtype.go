package dtype

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
)

var (
	ErrNil = errors.New("StatusPkg body is nil")
)

type StatusPkg struct {
	IP      string
	Name    string
	ConnNum int32
}

func (this *StatusPkg) Encode(w io.Writer) error {
	body, err := json.Marshal(this)
	if err != nil {
		return err
	}
	bodySize := int64(len(body))

	err = binary.Write(w, binary.BigEndian, bodySize)
	if err != nil {
		return err
	}
	n, err := w.Write(body)
	if err != nil {
		return err
	}

	return nil
}

func (this *StatusPkg) Decode(r io.Reader) error {
	var bodySize int64
	err := binary.Read(r, binary.BigEndian, bodySize)
	if err != nil {
		return err
	}

	if bodySize == 0 {
		return ErrNil
	}

	body := make([]byte, bodySize)
	err = io.ReadFull(r, body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, this)
	if err != nil {
		return err
	}
	return nil
}
