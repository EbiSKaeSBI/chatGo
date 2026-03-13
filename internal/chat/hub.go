package chat

type Hub struct {
	clients    map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan []byte
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.clients[client] = true

		case client := <-h.Unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.Broadcast:
			for client := range h.clients {
				client.send <- message
			}
		}
	}
}

func NewHub() *Hub {
	clients := make(map[*Client]bool)
	register := make(chan *Client)
	unregister := make(chan *Client)
	broadcast := make(chan []byte)

	return &Hub{clients: clients, Broadcast: broadcast, Register: register, Unregister: unregister}
}
