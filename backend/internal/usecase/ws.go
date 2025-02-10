package usecase

import (
	"log"
	"net/http"
	"websocket/backend/internal/entity"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Handler(hub *entity.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "Failed to upgrade to WebSocket", http.StatusInternalServerError)
			return
		}

		room := r.URL.Query().Get("room")
		if room == "" {
			http.Error(w, "Room parameter is missing", http.StatusBadRequest)
			return
		}

		name := r.URL.Query().Get("name")
		if name == "" {
			http.Error(w, "Name parameter is missing", http.StatusBadRequest)
			return
		}

		log.Printf("%s joined room %s", name, room)

		newClient := &entity.Client{
			Connection: conn,
			SendCh:     make(chan entity.Message, 256),
			Hub:        hub,
			RoomNumber: room,
			Name:       name,
		}

		hub.Register <- entity.RegisterMessage{Client: newClient, RoomNumber: room}

		go newClient.Read()
		go newClient.Write()
	}
}

func CountClientsPerRoom(hub *entity.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "Failed to upgrade to WebSocket", http.StatusInternalServerError)
			return
		}
		defer conn.Close()

		room := r.URL.Query().Get("room")
		if room == "" {
			http.Error(w, "Room parameter is missing", http.StatusBadRequest)
			return
		}

		for {
			clients := []string{}
		
			if roomClients, ok := hub.Rooms.Load(room); ok {
				for _, client := range roomClients.([]*entity.Client) {
					clients = append(clients, client.Name)
				}
			}
		
			response := map[string]interface{}{
				"room":    room,
				"clients": clients,
				"count":   len(clients),
			}
		
			if err := conn.WriteJSON(response); err != nil {
				log.Println("Failed to send the message:", err)
				break
			}
		}
		
	}
}

func Room(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "frontend/views/room.html")
}
