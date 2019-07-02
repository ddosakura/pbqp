package main

import (
	"log"
	"net"
	"os"
	"time"

	"github.com/ddosakura/pbqp"
	"github.com/ddosakura/pbqp/tests/proto/msg"
	"github.com/golang/protobuf/proto"
	"github.com/kr/pretty"
)

func main() {
	conn, err := net.Dial("unix", os.Args[1])
	if err != nil {
		panic(err)
	}
	s, err := pbqp.Client(conn)
	if err != nil {
		panic(err)
	}

	go func() {
		defer func() {
			log.Println("Stop Accept")
		}()
		for {
			c := s.Accept()
			if err := c.Test(); err != nil {
				log.Println("accept", err)
				break
			}
			go func() {
				defer func() {
					err := recover()
					if err != nil {
						log.Println("(Accept) session closed")
					}
				}()
				if err := c.
					Then(pbqp.ReadOnce(new(msg.Data))).
					Then(func(pb proto.Message, conn net.Conn) (proto.Message, proto.Message, error) {
						println("get pb success")
						pretty.Println(pb)
						return nil, new(msg.Data), nil
					}).
					Then(func(pb proto.Message, conn net.Conn) (proto.Message, proto.Message, error) {
						println("get pb success")
						pretty.Println(pb)
						return nil, nil, nil
					}).
					Finally(); err != nil {
					if err == pbqp.ErrIsClosed {
						panic(err)
					}
					log.Println("chain", err)
				}
			}()
		}
	}()

	i := 0
	for {
		time.Sleep(time.Second * 2)
		i++
		d := &msg.Data{
			Ver:     1,
			Payload: foo(i),
		}
		if err := s.
			Open().
			Then(pbqp.WriteOnce(d, nil)).
			Finally(); err != nil {
			if err == pbqp.ErrIsClosed {
				return
			}
			log.Println("o-chain", err)
		}
	}
}

func foo(n int) string {
	s := ""
	for i := 0; i < n; i++ {
		s += "="
	}
	return s
}
