package main

import (
	"crypto/rand"
	"fmt"
	"net"
	"strings"
)

const PORT = 22022

var clients []net.Conn
var game *Game
var gameTable map[string] *Game

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
	
	for {
		req, err := UnpackMsg(client)
		if err != nil {
			break
		}

		if req.Cmd == CREATE {
			game = NewGame()
			game.White = NewPlayer(req.User, client)
			client.Write([]byte(game.Name + "\n"))
			clientGame=game
			clientName=req.User
			fmt.Printf("Game created [%s]. White player[%s]\n", game.Name, game.White.User)
		}

		if req.Cmd == JOIN {
			// we need search the game in the gametable!! Now only have one game
			if req.Args == game.Name {
				game.Black = NewPlayer(req.User, client)
				client.Write([]byte(game.White.User + "\n"))
				game.White.Conn.Write([]byte(game.Black.User + "\n"))
				clientGame=game
				clientName=req.User
				fmt.Printf("Game start [%s]. Black player[%s]\n", game.Name, game.Black.User)
			} else {
				fmt.Printf("Game not found [%s]\n", req.Args)
			}
		}

		if req.Cmd == MOVE {
			args := strings.Split(req.Args, ":")
			if len(args) == 2 {
				// we need search the game in the gametable!! Now only have one game
				if req.User == game.Black.User {
					game.White.Conn.Write(req.PackMsg())
					fmt.Printf("Send move from [%s] to [%s]\n", game.Black.User, game.White.User)
				} else {
					game.Black.Conn.Write(req.PackMsg())
					fmt.Printf("Send move from [%s] to [%s]\n", game.White.User, game.Black.User)
				}
			}
		}

	}

	fmt.Printf("End connection from %s\n", client.RemoteAddr())
	if clientGame!=nil{
		msg:=NewMsg(END,fmt.Sprintf("El jugador %s ha abandonado la partida",clientName))
		if clientGame.White !=nil && clientGame.White.Conn !=nil{
			clientGame.White.Conn.Write(msg.PackMsg())
		}
		if clientGame.Black !=nil && clientGame.Black.Conn !=nil{
			clientGame.Black.Conn.Write(msg.PackMsg())
		}
	}
		
}
