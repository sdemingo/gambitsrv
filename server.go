package main

import (
	"crypto/rand"
	"fmt"
	"net"
)

const PORT = 22022

var clients []net.Conn
var match *Match

func RandomString(len int) string {
	b := make([]byte, len)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func InitServer() {
	srv, _ := net.Listen("tcp", ":"+fmt.Sprintf("%d", PORT))
	conns := clientConns(srv)

	clients = make([]net.Conn, 0)

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
			if client == nil {
				fmt.Printf("couldn't accept: " + fmt.Sprintf("%s", err))
				continue
			}
			i++
			fmt.Printf("%d: %v <-> %v\n", i, client.LocalAddr(), client.RemoteAddr())
			clients = append(clients, client)
			ch <- client
		}
	}()
	return ch
}

func handleConn(client net.Conn) {
	//b := bufio.NewReader(client)
	for {
		req, err := UnpackMsg(client)
		if err != nil {
			break
		}

		if req.Cmd == CREATE {
			match = NewMatch()
			client.Write([]byte(match.Name + "\n"))
		}
		/*message, _ := b.ReadString('\n')
		fields := strings.Split(message, ":")
		if len(fields) > 2 {
			user := fields[0]
			cmd := fields[1]
			//args := strings.Join(fields[1:], ":")

			if cmd == "/join" {
				if match == nil {
					// Creamos nueva partida y asignamos blancas
					match = NewMatch()

				}
				client.Write([]byte(match + "\n"))
			}
		}*/

	}

	fmt.Printf("End connection from %s\n", client.RemoteAddr())
}

/*
func sendToClients(message string) {
	for i := range clients {
		clients[i].Write([]byte(message))
	}
}
*/
