package chat

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	username string
	conn     *websocket.Conn
	send     chan []byte
	hub      *Hub
}

func (c *Client) SendGreeting() {
	c.send <- []byte("Welcome, " + c.username)
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
		c.hub.Broadcast <- message
	}
}

func (c *Client) WritePump() {
	for {
		message, ok := <-c.send
		if ok {
			err := c.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Println(err)
				break
			}
		} else {
			break
		}

	}

}

func NewClient(username string, conn *websocket.Conn, hub *Hub) *Client {
	send := make(chan []byte, 16)
	return &Client{username, conn, send, hub}
}
