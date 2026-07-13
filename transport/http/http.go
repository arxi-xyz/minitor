package http

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"minitor/config"
	"minitor/transport/socket"
	"net/http"
)

type Http struct {
	cfg config.Config
}

func NewHttp(cfg config.Config) *Http {
	return &Http{cfg: cfg}
}

func (s *Http) Run(ctx context.Context) error {
	mux := http.NewServeMux()
	hub := socket.NewHub()
	monitor := socket.NewMonitor(hub, socket.Settings{
		DefaultProcessLimit: s.cfg.Socket.DefaultProcessLimit,
		MaxProcessLimit:     s.cfg.Socket.MaxProcessLimit,
	})

	go monitor.Run(ctx)

	mux.HandleFunc("/health", healthHandler)
	mux.Handle("/ws", socket.NewHandler(hub, monitor))

	srv := &http.Server{
		Addr:    s.cfg.Server.Addr,
		Handler: mux,
	}

	shutdownTimeout, err := s.cfg.Server.ShutdownTimeoutDuration()
	if err != nil {
		return err
	}

	errCh := make(chan error, 1)
	go func() {
		log.Printf("listening on %s", s.cfg.Server.Addr)
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
