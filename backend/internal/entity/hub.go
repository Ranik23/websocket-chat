package entity

import (
	"log"
)

type Hub struct {
	Clients    map[*Client]bool
	Rooms      map[string][]*Client
	Broadcast  chan Message
	Register   chan RegisterMessage
	History    map[string][]Message
	Unregister chan RegisterMessage
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
		Register:   make(chan RegisterMessage),
		Unregister: make(chan RegisterMessage),
		Rooms:      make(map[string][]*Client),
		History:    make(map[string][]Message),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case regMessage := <-h.Register:
			client, roomNumber := regMessage.Client, regMessage.RoomNumber
			h.Clients[client] = true
			h.Rooms[roomNumber] = append(h.Rooms[roomNumber], client)

			if history, exists := h.History[roomNumber]; exists {
				for _, msg := range history {
					client.SendCh <- msg
				}
			}

			// h.broadcast <- Message{
			// 	SenderName: "System",
			// 	Text: client.name + " присоединился к комнате",
			// 	Room: roomNumber,
			// }

			log.Println("Client registered into the room", roomNumber)
		case unregMessage := <-h.Unregister:
			client, roomNumber := unregMessage.Client, unregMessage.RoomNumber
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)

				if roomClients, ok := h.Rooms[roomNumber]; ok {
					for i, c := range roomClients {
						if c == client {
							h.Rooms[roomNumber] = append(roomClients[:i], roomClients[i+1:]...)
							break
						}
					}
					if len(h.Rooms[roomNumber]) == 0 {
						delete(h.Rooms, roomNumber)
					}
				}

				close(client.SendCh)

				h.Broadcast <- Message{
					SenderName: "System",
					Text:       client.Name + " покинул комнату",
					Room:       roomNumber,
				}

				log.Println("Client unregistered from room:", roomNumber)
			}

		case message := <-h.Broadcast:
			roomNumber := message.Room
			for _, client := range h.Rooms[roomNumber] {
				select {
				case client.SendCh <- message:
					log.Println("Message sent to client - ", message, "room - ", roomNumber)
				default:
					close(client.SendCh)
					delete(h.Clients, client)
				}
			}
		}
	}
}
