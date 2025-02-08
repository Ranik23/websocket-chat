package main

import (
	"log"
	"github.com/gorilla/websocket"
)

type Client struct {
	sendCh chan []byte
	hub *Hub
	connection *websocket.Conn
}

func (c *Client) Read() {
	defer func() {
		c.hub.unregister <- c
		c.connection.Close()
	}()
	for {
		_, message, err := c.connection.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
		c.hub.broadcast <- message
	}
}

func (c *Client) Write() {
	defer c.connection.Close()
	for message := range c.sendCh {
		log.Println("got the message:", message)
		if err := c.connection.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Println(err)
			break
		}
	}
}