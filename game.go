package main

import "net"

type Player struct {
	User string   `json:"user"`
	Conn net.Conn `json:"-"`
}

func NewPlayer(username string, conn net.Conn) *Player {
	return &Player{username, conn}
}

type Game struct {
	Name  string  `json:"name"`
	White *Player `json:"white"`
	Black *Player `json:"black"`
}

func NewGame() *Game {
	m := new(Game)
	m.Name = RandomString(10)
	return m
}
