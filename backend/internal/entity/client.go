package entity

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
)

type Client struct {
	SendCh     chan Message
	Hub        *Hub
	Connection *websocket.Conn
	RoomNumber string
	Name       string `json:"name"`
}

func (c *Client) Read() {
	defer func() {
		c.Hub.Unregister <- RegisterMessage{
			Client:     c,
			RoomNumber: c.RoomNumber,
		}
		c.Connection.Close()
	}()
	for {
		_, message, err := c.Connection.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}

		var receivedMsg Message
		err = json.Unmarshal(message, &receivedMsg)
		if err != nil {
			log.Println("Ошибка парсинга JSON:", err)
			continue
		}

		receivedMsg.SenderName = c.Name 
		receivedMsg.Room = c.RoomNumber

		c.Hub.Broadcast <- receivedMsg
	}
}

func (c *Client) Write() {
	defer c.Connection.Close()
	for message := range c.SendCh {
		log.Println("text to be sent - ", message.Text)
		if err := c.Connection.WriteJSON(message); err != nil {
			log.Println(err)
			break
		}
	}
}
