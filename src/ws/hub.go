// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ws

import (
// "log"
)

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type hub struct {
	// Registered connections.
	connections map[*connection]bool

	// Inbound messages from the connections.
	broadcast chan []byte

	// Register requests from the connections.
	register chan *connection

	// Unregister requests from connections.
	unregister chan *connection
}

var hubMap = make(map[uint]*hub)

func makeHub(i uint) *hub {

	h := hub{
		broadcast:   make(chan []byte),
		register:    make(chan *connection),
		unregister:  make(chan *connection),
		connections: make(map[*connection]bool),
	}

	//log.Println("made hub")

	go h.run()
	hubMap[i] = &h

	//log.Println(hubMap)

	return &h
}

func Broadcast(id uint, b []byte) {
	if _, ok := hubMap[id]; ok {
		hubMap[id].broadcast <- b
		//log.Println("tried to broadcast")
	} else {
		//log.Println("im here?")
		//log.Println(id)
		//log.Println(hubMap)
	}
}

func (h *hub) run() {
	for {
		select {
		case c := <-h.register:
			h.connections[c] = true
		case c := <-h.unregister:
			if _, ok := h.connections[c]; ok {
				delete(h.connections, c)
				close(c.send)
			}
		case m := <-h.broadcast:
			for c := range h.connections {
				select {
				case c.send <- m:
					//log.Println("sending")
				default:
					//log.Println("closing")
					close(c.send)
					delete(h.connections, c)
				}
			}
		}
	}
}