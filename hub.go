package main

import "log"

type Hub struct {
	clients map[*Client]bool

	rooms map[string][]*Client // TO DO
	
	broadcast chan []byte

	register chan *Client

	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[*Client]bool),
		broadcast: make(chan []byte),
		register: make(chan *Client),
		unregister: make(chan *Client),
	}
}


func (h *Hub) Run() {
	for {
		select {
			case client := <-h.register:
				h.clients[client] = true
				log.Println("Client registered")
			case client := <-h.unregister:
				if _, ok := h.clients[client]; ok {
					delete(h.clients, client)
					close(client.sendCh)
					log.Println("Client unregistered")
				}
			case message := <- h.broadcast:
				log.Println("message sent")
				for client := range h.clients {
					select {
						case client.sendCh <- message:
							log.Println("message sent to client")
						default:
							close(client.sendCh)
							delete(h.clients, client)
					}
				}
		}
	}
}