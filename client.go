package main

import (
	"log"
	"github.com/gorilla/websocket"
)

type Client struct {
	sendCh 		chan Message
	hub 		*Hub
	connection 	*websocket.Conn
	roomNumber 	string
	name		string
}

func (c *Client) Read() {
	defer func() {
		c.hub.unregister <- RegisterMessage{
			client: c,
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
		Message := Message{
			SenderName: c.name,
			Text: string(message),
			Room: c.roomNumber,
		}
		c.hub.broadcast <- Message
	}
}

func (c *Client) Write() {
	defer c.connection.Close()
	for message := range c.sendCh {
		if err := c.connection.WriteJSON(message); err != nil {
			log.Println(err)
			break
		}
	}
}