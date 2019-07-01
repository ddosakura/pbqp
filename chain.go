package pbqp

import (
	"net"

	"github.com/golang/protobuf/proto"
	"github.com/xtaci/smux"
)

// Handler of Chain
//   func(readbuf, conn) (writebuf, next_readbuf, err)
type Handler func(proto.Message, net.Conn) (proto.Message, proto.Message, error)

// Chain for SMUX Action
type Chain struct {
	stream *smux.Stream

	pb  proto.Message
	err error
}

// Then Action
func (c *Chain) Then(fn Handler) *Chain {
	if c.err == nil {
		pb, next, err := fn(c.pb, c.stream)
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
		if next == nil {
			c.pb = nil
		} else {
			c.pb = next
			if err = Read(c.stream, c.pb); err != nil {
				c.err = err
				return c
			}
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
func ReadOnce(next proto.Message) Handler {
	return func(proto.Message, net.Conn) (proto.Message, proto.Message, error) {
		return nil, next, nil
	}
}

// WriteOnce Action
func WriteOnce(pb proto.Message, next proto.Message) Handler {
	return func(proto.Message, net.Conn) (proto.Message, proto.Message, error) {
		return pb, next, nil
	}
}
