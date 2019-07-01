package pbqp

import (
	"io"
	"net"

	"github.com/golang/protobuf/proto"
	"github.com/xtaci/smux"
)

// Builder func of smux
type Builder func(conn io.ReadWriteCloser, config *smux.Config) (*smux.Session, error)

// Generator of protobuf
type Generator func() proto.Message

// T Tunnel
type T struct {
	sess *smux.Session
}

func build(conn net.Conn, fn Builder) (*T, error) {
	sess, err := fn(conn, nil)
	if err != nil {
		return nil, err
	}
	return &T{
		sess: sess,
	}, nil
}

// Server Side
func Server(conn net.Conn) (*T, error) {
	return build(conn, smux.Server)
}

// Client Side
func Client(conn net.Conn) (*T, error) {
	return build(conn, smux.Client)
}

// Accept other-side
func (t *T) Accept(gen Generator) *Chain {
	s, err := t.sess.AcceptStream()
	if err != nil {
		return &Chain{
			err: RootErr(err),
		}
	}

	// dieCh := s.GetDieCh()
	// go func() {
	// 	select {
	// 	case <-dieCh:
	// 		println("accept stream die")
	// 	}
	// }()

	return &Chain{stream: s, gen: gen}
}

// Send to other-side
func (t *T) Send(gen Generator) *Chain {
	if t.sess.IsClosed() {
		return &Chain{
			err: ErrIsClosed,
		}
	}
	s, err := t.sess.OpenStream()
	if err != nil {
		return &Chain{
			err: RootErr(err),
		}
	}

	// dieCh := s.GetDieCh()
	// go func() {
	// 	select {
	// 	case <-dieCh:
	// 		println("send stream die")
	// 	}
	// }()

	return &Chain{stream: s, gen: gen}
}
