package handler

import "net/http"

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("Ok! Server is up and running!"))
	if err != nil {
		return
	}

}
