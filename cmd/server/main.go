package main

import "net/http"

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		_, err := w.Write([]byte("Ok! Server is up and running!"))
		if err != nil {
			return
		}
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		http.ServeFile(w, r, "./index.html")
	})

	addr := ":8080"

	err := http.ListenAndServe(addr, mux)
	if err != nil {
		return
	}
}
