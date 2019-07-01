package pbqp

import (
	"net"

	"github.com/golang/protobuf/proto"
	"github.com/xtaci/smux"
)

// Handler of Chain
type Handler func(req proto.Message, conn net.Conn) (res proto.Message, waitReq bool, e error)

// Chain for SMUX Action
type Chain struct {
	stream *smux.Stream
	gen    Generator

	pb  proto.Message
	err error
}

// Then Action
func (c *Chain) Then(fn Handler) *Chain {
	if c.err == nil {
		pb, w, err := fn(c.pb, c.stream)
		if err != nil {
			c.err = err
			return c
		}
		if pb != nil {
			if err = Write(c.stream, pb); err != nil {
				c.err = err
				return c
			}
		}
		if w {
			c.pb = c.gen()
			if err = Read(c.stream, c.pb); err != nil {
				c.err = err
				return c
			}
		} else {
			c.pb = nil
		}
	}
	return c
}

// Test return error in chain
func (c *Chain) Test() error {
	return c.err
}

// Finally Action
func (c *Chain) Finally() error {
	if c.stream != nil {
		c.stream.Close()
	}
	return c.err
}

// --- Quick Action(s) ---

// ReadOnce Action
func ReadOnce(req proto.Message, conn net.Conn) (res proto.Message, waitReq bool, e error) {
	return nil, true, nil
}

// WriteOnce Action
func WriteOnce(pb proto.Message, wait bool) Handler {
	return func(req proto.Message, conn net.Conn) (proto.Message, bool, error) {
		return pb, wait, nil
	}
}
