package pbqp

import (
	"io"
	"net"

	"github.com/xtaci/smux"
)

// Builder func of smux
type Builder func(conn io.ReadWriteCloser, config *smux.Config) (*smux.Session, error)

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
func (t *T) Accept() *Chain {
	s, err := t.sess.AcceptStream()
	err = RootErr(err)
	if err != nil {
		if err == io.EOF {
			err = ErrIsClosed
		}
		return &Chain{
			err: err,
		}
	}

	// dieCh := s.GetDieCh()
	// go func() {
	// 	select {
	// 	case <-dieCh:
	// 		println("accept stream die")
	// 	}
	// }()

	return &Chain{stream: s}
}

// Open other-side
func (t *T) Open() *Chain {
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

	return &Chain{stream: s}
}
