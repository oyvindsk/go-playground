package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
)

type room struct {
	toAll chan message

	// a channel for ppl whishing to join the room
	join chan *websocket.Conn

	// a channel for ppl wanting to leave the room
	//leave chan *client

	// all clients currently in the room
	clients map[*websocket.Conn]bool
}

func startRoom() *room {
	log.Println("StartRoom")
	r := new(room)

	r.toAll = make(chan message)
	r.join = make(chan *websocket.Conn)
	r.clients = make(map[*websocket.Conn]bool)

	go r.run()
	return r
}

func (r room) run() {
	log.Println("Running room")

	for {
		select {
		case c := <-r.join:
			log.Println("room run() join!")
			r.clients[c] = true
		case msg := <-r.toAll:
			log.Println("room run() message toAll: ", msg)
			for c := range r.clients {
				log.Println("room run() sending..")
				c.WriteMessage(websocket.TextMessage, msg.Msg)
			}
		}
	}
}
