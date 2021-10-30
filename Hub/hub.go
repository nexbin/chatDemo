package hub

import (
	. "ChatDemo/user"
	"fmt"
	"io"
	"strings"
	"sync"
)

type Hub struct {
	mapMutex   *sync.RWMutex
	UserMap    map[*User]bool
	Register   chan *User
	Unregister chan *User
	Broadcast  chan string
}

func NewHub() *Hub {
	return &Hub{
		mapMutex:   &sync.RWMutex{},
		Broadcast:  make(chan string),
		UserMap:    make(map[*User]bool),
		Register:   make(chan *User),
		Unregister: make(chan *User),
	}
}

func (h *Hub) StartHub() {
	for {
		select {
		case newUser := <-h.Register:
			h.mapMutex.Lock()
			h.UserMap[newUser] = true
			h.mapMutex.Unlock()
			fmt.Println("人数: ", len(h.UserMap))

		case quitUser := <-h.Unregister:
			h.mapMutex.Lock()
			delete(h.UserMap, quitUser)
			h.mapMutex.Unlock()
			fmt.Println("人数: ", len(h.UserMap))

		case allMsg := <-h.Broadcast:
			h.mapMutex.Lock()
			for user := range h.UserMap {
				io.Copy(user.Conn, strings.NewReader(allMsg))
			}
			h.mapMutex.Unlock()
		}
	}
}
