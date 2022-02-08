package main

import "net"

type Player struct {
	User string   `json:"user"`
	Conn net.Conn `json:"-"`
}

func NewPlayer(username string, conn net.Conn) *Player {
	return &Player{username, conn}
}

type Match struct {
	Name  string  `json:"name"`
	White *Player `json:"white"`
	Black *Player `json:"black"`
}

func NewMatch() *Match {
	m := new(Match)
	m.Name = RandomString(10)
	return m
}
