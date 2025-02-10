package entity

import (
	"log"
	"sync"
)

type Hub struct {
	Clients    map[*Client]bool
	Rooms      sync.Map
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
		History:    make(map[string][]Message),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case regMessage := <-h.Register:
			client, roomNumber := regMessage.Client, regMessage.RoomNumber
			h.Clients[client] = true

			clientsInRoom, _ := h.Rooms.LoadOrStore(roomNumber, []*Client{})
			clientsInRoom = append(clientsInRoom.([]*Client), client)

			h.Rooms.Store(roomNumber, clientsInRoom)

			if history, exists := h.History[roomNumber]; exists {
				for _, msg := range history {
					client.SendCh <- msg
				}
			}
			log.Println("Client registered into the room", roomNumber)
		case unregMessage := <-h.Unregister:
			client, roomNumber := unregMessage.Client, unregMessage.RoomNumber
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)

				if roomClients, ok := h.Rooms.Load(roomNumber); ok {
					clientsInRoom := roomClients.([]*Client)
					for i, c := range clientsInRoom {
						if c == client {
							clientsInRoom = append(clientsInRoom[:i], clientsInRoom[i+1:]...)
							h.Rooms.Store(roomNumber, clientsInRoom)
							break
						}
					}
					if len(clientsInRoom) == 0 {
						h.Rooms.Delete(roomNumber)
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
			if roomClients, ok := h.Rooms.Load(roomNumber); ok {
				for _, client := range roomClients.([]*Client) {
					select {
					case client.SendCh <- message:
						log.Println("Message sent to client:", client.Name, " - " , message, " Room - ", roomNumber)
					default:
						close(client.SendCh)
						delete(h.Clients, client)
					}
				}
			}
		}
	}
}
