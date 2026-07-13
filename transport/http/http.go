package http

import (
	"encoding/json"
	"log"
	"minitor/transport/socket"
	"net/http"
)

type Http struct {
}

func NewHttp() *Http {
	return &Http{}
}

func (s *Http) Run() {
	mux := http.NewServeMux()
	hub := socket.NewHub()
	monitor := socket.NewMonitor(hub)

	go monitor.Run()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(map[string]any{
			"status": "ok",
		}); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	})

	mux.Handle("/ws", socket.NewHandler(hub, monitor))

	log.Fatal(http.ListenAndServe(":8080", mux))
}
