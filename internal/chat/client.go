package chat

import (
	"encoding/json"
	"log"

	"github.com/EbiSKaeSBI/chatGo/internal/repository"
	"github.com/gorilla/websocket"
)

type Client struct {
	userId   int64
	username string
	conn     *websocket.Conn
	send     chan Message
	hub      *Hub
	repo     *repository.Repository
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
		errEncodingMessage := json.Unmarshal(message, &incoming)
		if errEncodingMessage != nil {
			log.Println(errEncodingMessage)
		}

		if err != nil {
			log.Println("Ошибка чтения или клиент отключился: ", err)
			c.hub.Unregister <- c
			break
		}
		log.Printf("recv: to=%s  text=%s\n", incoming.To, incoming.Text)
		errSaveMessage := c.repo.SaveMessage(int(c.userId), incoming.Text)
		if errSaveMessage != nil {
			log.Println(errSaveMessage)
		}
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

func NewClient(userId int64, username string, conn *websocket.Conn, hub *Hub, repo *repository.Repository) *Client {
	send := make(chan Message, 16)
	return &Client{userId, username, conn, send, hub, repo}
}
