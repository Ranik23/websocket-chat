package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Handler(hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		connection, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "Failed to upgrade to WebSocket", http.StatusInternalServerError)
			return
		}

		roomNumber := r.URL.Query().Get("room")
		if roomNumber == "" {
			http.Error(w, "Room parameter is missing", http.StatusBadRequest)
			return
		}

		name := r.URL.Query().Get("name")
		if name == "" {
			http.Error(w, "Name parameter is missing", http.StatusBadRequest)
			return
		}

		log.Println("Connected to room", roomNumber)

		client := &Client{
			connection: connection,
			sendCh:     make(chan Message, 256),
			hub:        hub,
			roomNumber: string(roomNumber),
			name:       name,
		}

		msg := RegisterMessage{
			client:     client,
			roomNumber: roomNumber,
		}

		hub.register <- msg

		go client.Read()
		go client.Write()
	}
}

func CountClientsPerRoom(hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		connection, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "Failed to upgrade to WebSocket", http.StatusInternalServerError)
			return
		}

		defer connection.Close()

		roomNumber := r.URL.Query().Get("room")
		if roomNumber == "" {
			http.Error(w, "Room parameter is missing", http.StatusBadRequest)
			return
		}

		for {
			time.Sleep(1 * time.Second)

			var clients []Client

			for _, client := range hub.rooms[roomNumber] {
				clients = append(clients, *client)
			}

			response := map[string]interface{}{
				"room":    roomNumber,
				"clients": clients,
				"count":   len(clients),
			}
			if err := connection.WriteJSON(response); err != nil {
				log.Println("Failed to send the message", err)
				break
			}
		}
	}
}

func Room(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "room.html")
}

func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}
