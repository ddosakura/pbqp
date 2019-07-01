package main

import (
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/ddosakura/pbqp"
	"github.com/ddosakura/pbqp/tests/proto/msg"
	"github.com/golang/protobuf/proto"
	"github.com/kr/pretty"
)

func main() {
	l, err := net.Listen("unix", os.Args[1])
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
				if err := s.
					Open().
					Then(pbqp.WriteOnce(d, nil)).
					Finally(); err != nil {
					if err == pbqp.ErrIsClosed {
						panic(err)
					}
					log.Println("o-chain", err)
				}
			}
		}()
	}
}
