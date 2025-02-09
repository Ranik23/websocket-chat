package main

import "log"


type Message struct {
	SenderName 	string `json:"sender"`
	Text 		string `json:"text"`
	Room 		string `json:"room"`
}

type RegisterMessage struct {
	client *Client
	roomNumber string
}


type Hub struct {

	clients map[*Client]bool

	rooms map[string][]*Client
	
	broadcast chan Message

	register chan RegisterMessage

	unregister chan RegisterMessage
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[*Client]bool),
		broadcast: make(chan Message),
		register: make(chan RegisterMessage),
		unregister: make(chan RegisterMessage),
		rooms: make(map[string][]*Client),
	}
}


func (h *Hub) Run() {
	for {
		select {
			case regMessage := <-h.register:
				client, roomNumber := regMessage.client, regMessage.roomNumber
				h.clients[client] = true
				h.rooms[roomNumber] = append(h.rooms[roomNumber], client)
				log.Println("Client registered into the room", roomNumber)
			case unregMessage := <-h.unregister:
				client, roomNumber := unregMessage.client, unregMessage.roomNumber
				if _, ok := h.clients[client]; ok {
					delete(h.clients, client)
					
					if roomClients, ok := h.rooms[roomNumber]; ok {
						for i, c := range roomClients {
							if c == client {
								h.rooms[roomNumber] = append(roomClients[:i], roomClients[i+1:]...)
								break
							}
						}
						if len(h.rooms[roomNumber]) == 0 {
							delete(h.rooms, roomNumber)
						}
					}
					
					close(client.sendCh)
					log.Println("Client unregistered from room:", roomNumber)
				}
			
			case message := <- h.broadcast:
				log.Println("Got the message")
				_, _, roomNumber := message.SenderName, message.Text, message.Room
				for _, client := range h.rooms[roomNumber] {
					select {
						case client.sendCh <- message:
							log.Println("Message sent to client; room - ", roomNumber)
						default:
							close(client.sendCh)
							delete(h.clients, client)
					}
				}
		}
	}
}