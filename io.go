package pbqp

import (
	"bytes"
	"io"
	"net"

	"github.com/golang/protobuf/proto"
)

// Read from conn
func Read(conn net.Conn, pb proto.Message) error {
	var buf bytes.Buffer
	_, err := io.Copy(&buf, conn)
	err = RootErr(err)
	if err != nil && err != io.EOF {
		return err
	}
	if err = proto.Unmarshal(buf.Bytes(), pb); err != nil {
		return err
	}
	return nil
}

// Write to conn
func Write(conn net.Conn, pb proto.Message) error {
	bs, err := proto.Marshal(pb)
	if err != nil {
		return err
	}
	if _, err = conn.Write(bs); err != nil {
		return RootErr(err)
	}
	return nil
}
