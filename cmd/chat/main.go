package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/EbiSKaeSBI/chatGo/internal/chat"
	"github.com/EbiSKaeSBI/chatGo/internal/config"
	"github.com/EbiSKaeSBI/chatGo/internal/handler"
	"github.com/EbiSKaeSBI/chatGo/internal/repository"
	_ "github.com/lib/pq"
)

func main() {

	mux := http.NewServeMux()
	cfg := config.Load()
	db, err := sql.Open("postgres", cfg.DSN)
	repo := repository.NewRepository(db)
	log.Println(db.Ping())
	defer db.Close()

	fileServer := http.FileServer(http.Dir("./web/styles"))
	mux.Handle("/styles/", http.StripPrefix("/styles", fileServer))

	hub := chat.NewHub()
	go hub.Run()
	h := handler.NewHandler(hub, repo)

	mux.HandleFunc("/", h.GetWebPage)

	mux.HandleFunc("/health", h.HealthCheck)

	mux.HandleFunc("/ws", h.WebSocket)

	log.Println("Listening on", cfg.Port)
	err = http.ListenAndServe(":"+cfg.Port, mux)
	if err != nil {
		log.Fatal(err)
	}
}
