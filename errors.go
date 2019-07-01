package pbqp

import (
	"errors"
	"net"

	errs "github.com/pkg/errors"
)

// error(s)
var (
	ErrIsClosed = errors.New("Mux Tunnel is Closed")
)

// RootErr resolve
func RootErr(e error) error {
	if e == nil {
		return nil
	}
	e = errs.Cause(e)
	err, ok := e.(*net.OpError)
	if !ok {
		return e
	}
	if err.Temporary() {
		return err
	}
	return ErrIsClosed
}
