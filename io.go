package pbqp

import (
	"bytes"
	"io"

	"github.com/golang/protobuf/proto"
)

// Read from conn
func Read(r io.ReadWriter, pb proto.Message) error {
	L := make([]byte, 8)
	r.Read(L)
	var buf bytes.Buffer
	_, err := io.CopyN(&buf, r, BStoI64(L))
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
	w.Write(I64toBS(int64(len(bs))))
	if _, err = w.Write(bs); err != nil {
		return RootErr(err)
	}
	return nil
}

// I64toBS Util
func I64toBS(n int64) []byte {
	bs := []byte{
		byte(n >> 56),
		byte(n >> 48),
		byte(n >> 40),
		byte(n >> 32),
		byte(n >> 24),
		byte(n >> 16),
		byte(n >> 8),
		byte(n),
	}
	return bs
}

// BStoI64 Util
func BStoI64(bs []byte) int64 {
	l := len(bs)
	if l < 8 {
		BS := make([]byte, 8-l, 8)
		bs = append(BS, bs...)
	} else if l > 8 {
		bs = bs[l-8:]
	}
	n := int64(bs[0])<<56 +
		int64(bs[1])<<48 +
		int64(bs[2])<<40 +
		int64(bs[3])<<32 +
		int64(bs[4])<<24 +
		int64(bs[5])<<16 +
		int64(bs[6])<<8 +
		int64(bs[7])
	return n
}
