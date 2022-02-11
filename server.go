package main

import (
	"crypto/rand"
	"fmt"
	"net"
	"strings"
)

const PORT = 22022

var clients []net.Conn

//var game *Game
var gameTable map[string]*Game

func RandomString(len int) string {
	b := make([]byte, len)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func InitServer() {
	srv, _ := net.Listen("tcp", ":"+fmt.Sprintf("%d", PORT))
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
			client.Write([]byte(clientGame.Name + "\n"))
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
				client.Write([]byte("\n")) // error jugadores con el mismo nombre
			}

			clientGame.Black = NewPlayer(req.User, client)
			client.Write([]byte(clientGame.White.User + "\n"))
			clientGame.White.Conn.Write([]byte(clientGame.Black.User + "\n"))
			clientName = req.User
			fmt.Printf("Game start [%s]. Black player[%s]\n", clientGame.Name, clientGame.Black.User)

		}

		if req.Cmd == MOVE {
			args := strings.Split(req.Args, ":")
			if len(args) == 2 {
				if clientGame, ok = gameTable[args[0]]; !ok {
					fmt.Printf("Game not found [%s]\n", args[0])
					return
				}

				// we need search the game in the gametable!! Now only have one game
				if req.User == clientGame.Black.User {
					clientGame.White.Conn.Write(req.PackMsg())
					fmt.Printf("Send move from [%s] to [%s]\n", clientGame.Black.User, clientGame.White.User)
				} else {
					clientGame.Black.Conn.Write(req.PackMsg())
					fmt.Printf("Send move from [%s] to [%s]\n", clientGame.White.User, clientGame.Black.User)
				}
			}
		}

	}

	fmt.Printf("End connection from %s\n", client.RemoteAddr())
	if clientGame != nil {
		msg := NewMsg(END, fmt.Sprintf("El jugador %s ha abandonado la partida", clientName))
		if clientGame.White != nil && clientGame.White.Conn != nil {
			clientGame.White.Conn.Write(msg.PackMsg())
		}
		if clientGame.Black != nil && clientGame.Black.Conn != nil {
			clientGame.Black.Conn.Write(msg.PackMsg())
		}
	}

}
