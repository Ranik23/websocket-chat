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

func CountClientsPerRoom(w http.ResponseWriter, r *http.Request) {
	_, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to upgrade to WebSocket", http.StatusInternalServerError)
		return
	}

	roomNumber := r.URL.Query().Get("room")
	if roomNumber == "" {
		http.Error(w, "Room parameter is missing", http.StatusBadRequest)
		return
	}

	// для вывода количества клиентов в комнате
	
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
