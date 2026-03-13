package chat

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	username string
	conn     *websocket.Conn
	send     chan Message
	hub      *Hub
}

func (c *Client) SendGreeting() {
	c.send <- Message{
		Type: "system",
		Text: "Welcome, " + c.username,
	}
}

func (c *Client) ReadPump() {
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("Ошибка чтения или клиент отключился: ", err)
			c.hub.Unregister <- c
			break
		}
		log.Printf("recv: %s\n", string(message))
		c.hub.Broadcast <- Message{Type: "message", Username: c.username, Text: string(message)}
	}
}

func (c *Client) WritePump() {
	for {
		message, ok := <-c.send

		if ok {
			data, err := json.Marshal(message)
			if err != nil {
				log.Println(err)
				break
			}
			err = c.conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Println(err)
				c.hub.Unregister <- c
				break
			}
		} else {
			break
		}

	}

}

func NewClient(username string, conn *websocket.Conn, hub *Hub) *Client {
	send := make(chan Message, 16)
	return &Client{username, conn, send, hub}
}
