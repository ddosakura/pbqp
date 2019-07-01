package pbqp

import (
	"bytes"
	"io"

	"github.com/golang/protobuf/proto"
)

// Read from conn
func Read(r io.ReadWriter, pb proto.Message) error {
	var buf bytes.Buffer
	_, err := io.Copy(&buf, r)
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
func Write(w io.ReadWriter, pb proto.Message) error {
	bs, err := proto.Marshal(pb)
	if err != nil {
		return err
	}
	if _, err = w.Write(bs); err != nil {
		return RootErr(err)
	}
	return nil
}
