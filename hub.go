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
	clients 	map[*Client]bool
	rooms 		map[string][]*Client
	broadcast 	chan Message
	register 	chan RegisterMessage
	history 	map[string][]Message
	unregister 	chan RegisterMessage
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[*Client]bool),
		broadcast: make(chan Message),
		register: make(chan RegisterMessage),
		unregister: make(chan RegisterMessage),
		rooms: make(map[string][]*Client),
		history: make(map[string][]Message),
	}
}


func (h *Hub) Run() {
	for {
		select {
			case regMessage := <-h.register:
				client, roomNumber := regMessage.client, regMessage.roomNumber
				h.clients[client] = true
				h.rooms[roomNumber] = append(h.rooms[roomNumber], client)


				if history, exists := h.history[roomNumber]; exists {
					for _, msg := range history {
						client.sendCh <- msg
					}
				}

				// h.broadcast <- Message{
				// 	SenderName: "System",
				// 	Text: client.name + " присоединился к комнате",
				// 	Room: roomNumber,
				// }

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

					h.broadcast <- Message{
						SenderName: "System",
						Text:       client.name + " покинул комнату",
						Room:       roomNumber,
					}

					log.Println("Client unregistered from room:", roomNumber)
				}
			
			case message := <- h.broadcast:
				roomNumber := message.Room
				for _, client := range h.rooms[roomNumber] {
					select {
						case client.sendCh <- message:
							log.Println("Message sent to client - ", message, "room - ", roomNumber)
						default:
							close(client.sendCh)
							delete(h.clients, client)
					}
				}
		}
	}
}