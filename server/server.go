package server

import (
	hub "ChatDemo/Hub"
	. "ChatDemo/user"
	"fmt"
	"log"
	"net"
	"strings"
)

type server struct {
	ip   string
	port string
	Hub  *hub.Hub
}

// NewServer return a server struct
func NewServer(ip, port string) *server {
	return &server{
		ip:   ip,
		port: port,
		Hub:  hub.NewHub(),
	}
}

// StartServer and listen connections
func (s *server) StartServer() error {
	//socket listener
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", s.ip, s.port))
	if err != nil {
		return err
	}

	// defer close conn
	defer listener.Close()

	// prepare to accept
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("%v\n", fmt.Errorf("listen conn failed"))
			continue
		}
		fmt.Println("accept: ", conn.RemoteAddr())
		// handle conn
		go s.handleConn(conn)
	}

}

// handleConn
func (s *server) handleConn(conn net.Conn) {
	// register
	newUser := &User{Conn: conn, Ip: conn.RemoteAddr().String(), UserName: conn.RemoteAddr().String()}
	s.Hub.Register <- newUser
	s.readClientMessage(newUser)
}

// 读取客户端发送到服务端的消息
func (s *server) readClientMessage(user *User) {
	s.Hub.Broadcast <- formatMessage(user.Ip, "<---用户上线了--->\n")
	for {
		var buf [128]byte
		n, err := user.Conn.Read(buf[:])
		if err != nil {
			log.Println("Read from tcp server failed,err:", err)
			s.Hub.Unregister <- user
			break
		}
		// 查询当前在线人数
		msg := string(buf[:n])
		if msg == "who" {
			user.Conn.Write([]byte(s.CheckAllOnlineUser()))
		} else {
			msg = formatMessage(user.Ip, msg+"\n")
			s.Hub.Broadcast <- msg
		}

	}
}

func formatMessage(ip, msg string) string {
	return fmt.Sprintf("[%s]: %s", ip, msg)
}

func (s *server) CheckAllOnlineUser() string {
	sb := strings.Builder{}
	for u, _ := range s.Hub.UserMap {
		sb.WriteString("[" + u.Ip + "]: " + u.UserName + "在线\n")
	}
	return sb.String()
}
