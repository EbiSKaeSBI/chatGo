package chat

type Hub struct {
	clients    map[string]*Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan Message
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.clients[client.username] = client
			msg := Message{Type: "system", Text: client.username + " joined the chat"}
			for _, client := range h.clients {
				client.send <- msg
			}
		case client := <-h.Unregister:
			if _, ok := h.clients[client.username]; ok {
				delete(h.clients, client.username)
				close(client.send)
				msg := Message{Type: "system", Text: client.username + " left the chat"}

				for _, client := range h.clients {
					client.send <- msg
				}
			}
		case message := <-h.Broadcast:
			sender, senderOk := h.clients[message.From]
			recipient, recipientOk := h.clients[message.To]
			if recipientOk {
				recipient.send <- message
			} else if senderOk {
				sender.send <- Message{Type: "system", Text: message.To + " is offline"}
			}

			if senderOk && message.From != message.To {
				sender.send <- message
			}
		}
	}
}

func NewHub() *Hub {
	clients := make(map[string]*Client)
	register := make(chan *Client)
	unregister := make(chan *Client)
	broadcast := make(chan Message)

	return &Hub{clients: clients, Broadcast: broadcast, Register: register, Unregister: unregister}
}
