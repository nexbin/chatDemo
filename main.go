package main

import (
	Server "ChatDemo/server"
	_ "ChatDemo/user"
	"log"
)

func main() {
	// create server
	server := Server.NewServer("127.0.0.1", "8890")
	// start hub
	go server.Hub.StartHub()
	// start server
	if err := server.StartServer(); err != nil {
		log.Fatalln(err)
	}

}
