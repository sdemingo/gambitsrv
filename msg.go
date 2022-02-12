package main

import (
	"encoding/json"
	"net"
)

const (
	ERROR  = 0
	OK     = 1
	CREATE = 2
	JOIN   = 3
	MOVE   = 4
	END    = 5
)

type Msg struct {
	Cmd  byte
	User string
	Args string
}

func NewMsg(cmd byte, username string, args string) *Msg {
	return &Msg{cmd, username, args}
}

func UnpackMsg(conn net.Conn) (*Msg, error) {
	d := json.NewDecoder(conn)
	var msg Msg
	err := d.Decode(&msg)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

func (m *Msg) PackMsg() []byte {
	b, _ := json.Marshal(m)
	return b
}
