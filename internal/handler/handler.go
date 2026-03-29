package handler

import (
	"github.com/EbiSKaeSBI/chatGo/internal/chat"
	"github.com/EbiSKaeSBI/chatGo/internal/repository"
)

type Handler struct {
	hub  *chat.Hub
	repo *repository.Repository
}

func NewHandler(hub *chat.Hub, repo *repository.Repository) *Handler {
	return &Handler{
		hub:  hub,
		repo: repo,
	}
}
