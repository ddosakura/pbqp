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
	l, err := net.Listen("unix", "../swap.sock")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			println(err)
			log.Println(err)
			pretty.Println("outside", err)
			continue
		}
		s, err := pbqp.Server(conn)
		if err != nil {
			log.Println(err)
			continue
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

		go func() {
			defer func() {
				err := recover()
				if err != nil {
					log.Println("(Send) session closed")
				}
			}()
			i := 0
			for {
				time.Sleep(time.Second * 4)
				i++
				d := &msg.Data{
					Ver:     1,
					Payload: "No. " + strconv.Itoa(i),
				}
				if err := s.Send(gen).Then(pbqp.WriteOnce(d, false)).Finally(); err != nil {
					if err == pbqp.ErrIsClosed {
						panic(err)
					}
					log.Println(err)
				}
			}
		}()
	}
}
