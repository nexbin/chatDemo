package user

import (
	"net"
)

type User struct {
	Conn     net.Conn
	UserName string
	Ip       string
}
