package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
)

type Client struct {
	sendCh     chan Message
	hub        *Hub
	connection *websocket.Conn
	roomNumber string
	name       string `json:"name"`
}

func (c *Client) Read() {
	defer func() {
		c.hub.unregister <- RegisterMessage{
			client:     c,
			roomNumber: c.roomNumber,
		}
		c.connection.Close()
	}()
	for {
		_, message, err := c.connection.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}

		var receivedMsg Message
		err = json.Unmarshal(message, &receivedMsg) // Парсим JSON
		if err != nil {
			log.Println("Ошибка парсинга JSON:", err)
			continue
		}

		receivedMsg.SenderName = c.name // Подставляем имя клиента
		receivedMsg.Room = c.roomNumber

		c.hub.broadcast <- receivedMsg
	}
}

func (c *Client) Write() {
	defer c.connection.Close()
	for message := range c.sendCh {
		log.Println("text to be sent - ", message.Text)
		if err := c.connection.WriteJSON(message); err != nil {
			log.Println(err)
			break
		}
	}
}
