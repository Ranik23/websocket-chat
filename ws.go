package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Handler(hub *Hub) http.HandlerFunc {
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

		newClient := &Client{
			connection: conn,
			sendCh:     make(chan Message, 256),
			hub:        hub,
			roomNumber: room,
			name:       name,
		}

		hub.register <- RegisterMessage{client: newClient, roomNumber: room}

		go newClient.Read()
		go newClient.Write()
	}
}
// убрать отсюда, тк нужна либо синк мапа либо паника
func CountClientsPerRoom(hub *Hub) http.HandlerFunc {
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
			clients := make([]string, 0, len(hub.rooms[room]))
			for _, client := range hub.rooms[room] {
				clients = append(clients, client.name)
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
