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
		var incoming struct {
			To   string `json:"to"`
			Text string `json:"text"`
		}
		_, message, err := c.conn.ReadMessage()
		err = json.Unmarshal(message, &incoming)

		if err != nil {
			log.Println("Ошибка чтения или клиент отключился: ", err)
			c.hub.Unregister <- c
			break
		}
		log.Printf("recv: to=%s  text=%s\n", incoming.To, incoming.Text)
		c.hub.Broadcast <- Message{Type: "message", From: c.username, To: incoming.To, Text: incoming.Text}
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
