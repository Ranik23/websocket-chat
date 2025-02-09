package main

import (
	"flag"
	"log"
	"net/http"
	"websocket/backend/internal/entity"
	"websocket/backend/internal/usecase"
)

func main() {
	addr := flag.String("addr", ":8080", "адрес сервера (например, :8080)")
	flag.Parse()

	hub := entity.NewHub()
	go hub.Run()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("frontend/static"))))

	http.HandleFunc("/", usecase.Room)
	http.HandleFunc("/ws", usecase.Handler(hub))
	http.HandleFunc("/wscount", usecase.CountClientsPerRoom(hub))

	log.Println("Сервер запущен на", *addr)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
