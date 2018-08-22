package test

import (
	"bufio"
	"fmt"
	"net"
)

func StartServer() {
	server, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic("couldn't start listening: ")
	}
	conns := clientConns(server)
	for {
		go handleConn(<-conns)
	}

}
func clientConns(listener net.Listener) chan net.Conn {
	ch := make(chan net.Conn)
	i := 0
	go func() {
		for {
			client, err := listener.Accept()
			if err != nil {
				fmt.Printf("couldn't accept ", err)
				continue
			}
			i++
			fmt.Printf("%d: %v <-> %v\n", i, client.LocalAddr(), client.RemoteAddr())
			ch <- client
		}
	}()
	return ch
}

func handleConn(client net.Conn) {
	b := bufio.NewReader(client)
	for {
		line, err := b.ReadBytes('\n')
		if err != nil { // EOF, or worse
			break
		}

		client.Write(line)
	}
}
