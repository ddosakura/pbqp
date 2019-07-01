package pbqp

import (
	"net"

	"github.com/golang/protobuf/proto"
)

// Read from conn
func Read(conn net.Conn, pb proto.Message) error {
	// TODO: get full data
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return RootErr(err)
	}
	if err = proto.Unmarshal(buf[:n], pb); err != nil {
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
