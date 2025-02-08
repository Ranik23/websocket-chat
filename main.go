package main

import (
	"flag"
	"log"
	"net/http"
)


func main() {
	addr := flag.String("addr", ":8080", "адрес сервера (например, :8080)")
	flag.Parse()

	hub := NewHub()
	go hub.Run()

	http.HandleFunc("/", Home)
	http.HandleFunc("/ws", Handler(hub))

	log.Println("Сервер запущен на", *addr)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}