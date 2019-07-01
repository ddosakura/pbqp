package main

import (
	"log"
	"net"
	"strconv"
	"time"

	"github.com/ddosakura/pbqp"
	"github.com/ddosakura/pbqp/tests/proto/msg"
	"github.com/golang/protobuf/proto"
	"github.com/kr/pretty"
)

func gen() proto.Message {
	return new(msg.Data)
}

func main() {
	conn, err := net.Dial("unix", "../swap.sock")
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
			c := s.Accept(gen)
			if err := c.Test(); err != nil {
				log.Println(err)
				break
			}
			go func() {
				defer func() {
					err := recover()
					if err != nil {
						log.Println("(Accept) session closed")
					}
				}()
				if err := c.Then(pbqp.ReadOnce).Then(func(pb proto.Message, conn net.Conn) (proto.Message, bool, error) {
					println("get pb success")
					pretty.Println(pb)
					return nil, false, nil
				}).Finally(); err != nil {
					if err == pbqp.ErrIsClosed {
						panic(err)
					}
					log.Println(err)
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
			Payload: "No. " + strconv.Itoa(i),
		}
		if err := s.Send(gen).Then(pbqp.WriteOnce(d, false)).Finally(); err != nil {
			if err == pbqp.ErrIsClosed {
				return
			}
			log.Println(err)
		}
	}
}
