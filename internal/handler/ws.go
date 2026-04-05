package handler

import (
	"log"
	"net/http"
	"strings"

	"github.com/EbiSKaeSBI/chatGo/internal/chat"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func (h *Handler) WebSocket(w http.ResponseWriter, r *http.Request) {

	q := r.URL.Query()
	username := strings.TrimSpace(q.Get("username"))
	password := strings.TrimSpace(q.Get("password"))

	userId, err := h.repo.FindUser(username)
	if err != nil {
		err := h.repo.CreateUser(username, password)
		userId, _ = h.repo.FindUser(username)
		if err != nil {
			log.Println(err)
		}
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	//if username == "" {
	//	http.Error(w, "Missing username", http.StatusBadRequest)
	//	return
	//}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	log.Println("websocket connected:", username)

	client := chat.NewClient(userId, username, conn, h.hub, h.repo)
	h.hub.Register <- client
	client.SendGreeting()
	go client.WritePump()
	client.ReadPump()

}
