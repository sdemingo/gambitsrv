package main

import (
	"crypto/rand"
	"fmt"
	"net"
	"os"
)


var clients []net.Conn
var gameTable map[string]*Game

func RandomString(len int) string {
	b := make([]byte, len)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func InitServer(port string) {
	srv, err := net.Listen("tcp", ":"+port)
	if err!=nil{
		fmt.Println(err)
		os.Exit(1)
	}
	conns := clientConns(srv)

	clients = make([]net.Conn, 0)
	gameTable = make(map[string]*Game)

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
	var clientGame *Game = nil
	var clientName string
	var ok bool

	for {
		req, err := UnpackMsg(client)
		if err != nil {
			break
		}

		if req.Cmd == CREATE {
			clientGame = NewGame()
			clientGame.White = NewPlayer(req.User, client)
			client.Write(NewMsg(OK, req.User, clientGame.Name).PackMsg())
			clientName = req.User
			gameTable[clientGame.Name] = clientGame
			fmt.Printf("Game created [%s]. White player[%s]\n", clientGame.Name, clientGame.White.User)
		}

		if req.Cmd == JOIN {
			if clientGame, ok = gameTable[req.Args]; !ok {
				fmt.Printf("Game not found [%s]\n", req.Args)
				return
			}

			if req.User == clientGame.White.User {
				client.Write(NewMsg(ERROR, req.User,
					"No se puede usar el mismo nombre que el oponente").PackMsg())
			}

			clientGame.Black = NewPlayer(req.User, client)
			client.Write(NewMsg(OK, req.User, clientGame.White.User).PackMsg())
			clientGame.White.Conn.Write(NewMsg(OK, req.User, clientGame.Black.User).PackMsg())
			clientName = req.User
			fmt.Printf("Game start [%s]. Black player[%s]\n", clientGame.Name, clientGame.Black.User)

		}

		if req.Cmd == MOVE {
			if req.User == clientGame.Black.User {
				clientGame.White.Conn.Write(req.PackMsg())
				fmt.Printf("Send move from [%s] to [%s]\n", clientGame.Black.User, clientGame.White.User)
			} else {
				clientGame.Black.Conn.Write(req.PackMsg())
				fmt.Printf("Send move from [%s] to [%s]\n", clientGame.White.User, clientGame.Black.User)
			}
		}

		if req.Cmd == END {
			if req.User == clientGame.Black.User {
				clientGame.White.Conn.Write(NewMsg(END, req.User,
					fmt.Sprintf("%s ha ganado la partida", clientGame.White.User)).PackMsg())
				fmt.Printf("Send checkmate. [%s] win the game\n", clientGame.White.User)
			} else {
				clientGame.Black.Conn.Write(NewMsg(END, req.User,
					fmt.Sprintf("%s ha ganado la partida", clientGame.Black.User)).PackMsg())
				fmt.Printf("Send checkmate. [%s] win the game\n", clientGame.Black.User)
			}
		}

	}

	fmt.Printf("End connection from %s\n", client.RemoteAddr())
	if clientGame != nil {
		msg := NewMsg(END, clientName, fmt.Sprintf("El jugador %s ha abandonado la partida", clientName))
		if clientGame.White != nil && clientGame.White.Conn != nil {
			clientGame.White.Conn.Write(msg.PackMsg())
		}
		if clientGame.Black != nil && clientGame.Black.Conn != nil {
			clientGame.Black.Conn.Write(msg.PackMsg())
		}
	}

}
