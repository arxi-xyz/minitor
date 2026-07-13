package http

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"minitor/transport/socket"
	"net/http"
	"time"
)

const shutdownTimeout = 10 * time.Second

type Http struct {
	addr string
}

func NewHttp() *Http {
	return &Http{addr: ":8080"}
}

func (s *Http) Run(ctx context.Context) error {
	mux := http.NewServeMux()
	hub := socket.NewHub()
	monitor := socket.NewMonitor(hub)

	go monitor.Run(ctx)

	mux.HandleFunc("/health", healthHandler)
	mux.Handle("/ws", socket.NewHandler(hub, monitor))

	srv := &http.Server{
		Addr:    s.addr,
		Handler: mux,
	}

	errCh := make(chan error, 1)
	go func() {
		log.Printf("listening on %s", s.addr)
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		log.Println("shutting down...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		return srv.Shutdown(shutdownCtx)
	case err := <-errCh:
		return err
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
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
}
