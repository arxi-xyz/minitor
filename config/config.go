package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

const envPrefix = "MINITOR_"

type Config struct {
	Server Server `json:"server"`
	Client Client `json:"client"`
	Socket Socket `json:"socket"`
}

type Server struct {
	Addr            string `json:"addr"`
	ShutdownTimeout string `json:"shutdown_timeout"`
}

type Client struct {
	WSURL          string `json:"ws_url"`
	ReconnectDelay string `json:"reconnect_delay"`
}

type Socket struct {
	DefaultProcessLimit int `json:"default_process_limit"`
	MaxProcessLimit     int `json:"max_process_limit"`
}

type LoadOptions struct {
	Path string

	ServerAddr string
	WSURL      string
}

func Default() Config {
	return Config{
		Server: Server{
			Addr:            ":8080",
			ShutdownTimeout: "10s",
		},
		Client: Client{
			WSURL:          "ws://127.0.0.1:8080/ws",
			ReconnectDelay: "3s",
		},
		Socket: Socket{
			DefaultProcessLimit: 50,
			MaxProcessLimit:     200,
		},
	}
}

func Load(opts LoadOptions) (Config, error) {
	cfg := Default()

	path := opts.Path
	if path == "" {
		path = os.Getenv(envPrefix + "CONFIG")
	}

	if path != "" {
		data, err := os.ReadFile(path)
		if err != nil {
			return Config{}, fmt.Errorf("read config: %w", err)
		}

		if err := json.Unmarshal(data, &cfg); err != nil {
			return Config{}, fmt.Errorf("parse config: %w", err)
		}
	}

	applyEnv(&cfg)

	if opts.ServerAddr != "" {
		cfg.Server.Addr = opts.ServerAddr
	}
	if opts.WSURL != "" {
		cfg.Client.WSURL = opts.WSURL
	}

	return cfg, cfg.Validate()
}

func (c Config) Validate() error {
	if c.Server.Addr == "" {
		return errors.New("server.addr is required")
	}

	if _, err := c.Server.ShutdownTimeoutDuration(); err != nil {
		return fmt.Errorf("server.shutdown_timeout: %w", err)
	}

	if c.Client.WSURL == "" {
		return errors.New("client.ws_url is required")
	}

	if _, err := c.Client.ReconnectDelayDuration(); err != nil {
		return fmt.Errorf("client.reconnect_delay: %w", err)
	}

	if c.Socket.DefaultProcessLimit <= 0 {
		return errors.New("socket.default_process_limit must be positive")
	}

	if c.Socket.MaxProcessLimit < c.Socket.DefaultProcessLimit {
		return errors.New("socket.max_process_limit must be >= default_process_limit")
	}

	return nil
}

func (s Server) ShutdownTimeoutDuration() (time.Duration, error) {
	return time.ParseDuration(s.ShutdownTimeout)
}

func (c Client) ReconnectDelayDuration() (time.Duration, error) {
	return time.ParseDuration(c.ReconnectDelay)
}

func applyEnv(cfg *Config) {
	if v := os.Getenv(envPrefix + "SERVER_ADDR"); v != "" {
		cfg.Server.Addr = v
	}
	if v := os.Getenv(envPrefix + "SERVER_SHUTDOWN_TIMEOUT"); v != "" {
		cfg.Server.ShutdownTimeout = v
	}
	if v := os.Getenv(envPrefix + "CLIENT_WS_URL"); v != "" {
		cfg.Client.WSURL = v
	}
	if v := os.Getenv(envPrefix + "CLIENT_RECONNECT_DELAY"); v != "" {
		cfg.Client.ReconnectDelay = v
	}
	if v, ok := envInt(envPrefix + "SOCKET_DEFAULT_PROCESS_LIMIT"); ok {
		cfg.Socket.DefaultProcessLimit = v
	}
	if v, ok := envInt(envPrefix + "SOCKET_MAX_PROCESS_LIMIT"); ok {
		cfg.Socket.MaxProcessLimit = v
	}
}

func envInt(key string) (int, bool) {
	v := os.Getenv(key)
	if v == "" {
		return 0, false
	}

	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, false
	}

	return n, true
}
