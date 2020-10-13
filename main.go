package main

import (
	"bufio"
	"context"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"
)

var text = []byte(`HTTP/1.1 200 OK
Server: HA-Counter (github.com/kiselev-nikolay/ha-counter)
Connection: keep-alive
Cache-Control: no-store
Access-Control-Allow-Origin: *
Content-Type: text/plain
Content-Length: 2

ok`)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide a port number!")
		return
	}

	PORT := ":" + arguments[1]
	l, err := net.Listen("tcp4", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	rand.Seed(time.Now().Unix())

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		defer c.Close()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go func() {
			select {
			default:
				buff := bufio.NewReader(c)
				netData, err := buff.ReadString('\n')
				if err != nil {
					return
				}
				request := strings.Fields(netData)
				fmt.Println(request[1])
				c.Write(text)
				return
			case <-ctx.Done():
				return
			}
		}()

		go func() {
			<-time.After(500 * time.Millisecond)
			return
		}()
	}
}
