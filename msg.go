package main

import (
	"encoding/json"
	"fmt"
	"net"
)

const (
	CREATE = "create"
	JOIN   = "join"
	MOVE   = "move"
)

type Msg struct {
	Cmd  string
	User string
	Args string
}

func NewMsg(cmd string, username string) *Msg {
	return &Msg{cmd, username, ""}
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

func (m *Msg) String() string {
	return fmt.Sprintf("[%s] %s %s\n", m.User, m.Cmd, m.Args)
}
