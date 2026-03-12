package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

func main() {

	mux := http.NewServeMux()

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("Ok! Server is up and running!"))
		if err != nil {
			return
		}

	})

	fileServer := http.FileServer(http.Dir("./web/styles"))
	mux.Handle("/styles/", http.StripPrefix("/styles", fileServer))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		http.ServeFile(w, r, "./web/index.html")
	})

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		username := strings.TrimSpace(q.Get("username"))

		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if username == "" {
			http.Error(w, "Missing username", http.StatusBadRequest)
			return
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
		}
		defer func(conn *websocket.Conn) {
			err := conn.Close()
			if err != nil {

			}
		}(conn)
		log.Println("websocket connected:", username)

		err = conn.WriteMessage(1, []byte("Welcome, "+username))
		if err != nil {
			log.Println(err)
		}

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Ошибка чтения или клиент отключился: ", err)
				break
			}
			log.Printf("recv: %s\n", string(message))

			if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Println("Ошибка записи: ", err)
				break
			}
		}

	})

	addr := ":8080"

	log.Println("Listening on", addr)
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		log.Fatal(err)
	}
}
