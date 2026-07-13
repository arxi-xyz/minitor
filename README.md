# Minitor

A small system monitor with a client/server setup. The server collects OS metrics and pushes them over WebSocket. The TUI client just listens and shows you what's going on.

## What you get

- **CPU**, **RAM**, **Disk**, and **Network** stats
- **Process list** with a tree view and scroll
- Real-time updates over **WebSocket**
- A terminal UI built with [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- **Graceful shutdown** on `SIGINT` / `SIGTERM`
- **Config** via JSON file, env vars, or CLI flags

## Requirements

- Go 1.24+
- Linux (collectors lean on `/proc` and gopsutil)

## Setup

```bash
git clone <repo-url>
cd minitor
go mod download
```

## Running it

### Server

```bash
go run main.go
```

That starts the server on `:8080`:

| Endpoint | What it does |
|----------|--------------|
| `GET /health` | health check |
| `GET /ws` | WebSocket |

### TUI client

Open a second terminal:

```bash
go run ./cmd/client
```

Keyboard shortcuts:

| Key | Action |
|-----|--------|
| `j` / `k` | scroll the process list |
| `pgup` / `pgdown` | scroll faster |
| `q` / `Esc` / `Ctrl+C` | quit |

## Config

There's an example file at [`config.example.json`](config.example.json).

```bash
# server
go run main.go -config config.example.json
go run main.go -addr :9090

# client
go run ./cmd/client -config config.example.json
go run ./cmd/client -url ws://localhost:9090/ws
```

### Env vars

| Variable | Default |
|----------|---------|
| `MINITOR_CONFIG` | — (path to JSON file) |
| `MINITOR_SERVER_ADDR` | `:8080` |
| `MINITOR_SERVER_SHUTDOWN_TIMEOUT` | `10s` |
| `MINITOR_CLIENT_WS_URL` | `ws://127.0.0.1:8080/ws` |
| `MINITOR_CLIENT_RECONNECT_DELAY` | `3s` |
| `MINITOR_SOCKET_DEFAULT_PROCESS_LIMIT` | `50` |
| `MINITOR_SOCKET_MAX_PROCESS_LIMIT` | `200` |

Load order: **defaults → file → env → flags**

## WebSocket API

Every message is JSON shaped like this:

```json
{ "type": "<name>", "data": { ... } }
```

### Server → client

**`snapshot`** — system metrics (no full process list):

```json
{
  "type": "snapshot",
  "data": {
    "cpu": { "usage": 12.5, "core_usage": [...], ... },
    "ram": { ... },
    "disk": { ... },
    "network": { ... },
    "process_count": 342
  }
}
```

**`processes`** — paginated process list (response to a request)

**`pong`** — reply to ping

### Client → server

**`ping`**

```json
{ "type": "ping" }
```

**`processes`**

```json
{ "type": "processes", "data": { "offset": 0, "limit": 50 } }
```

## Build

```bash
task build                              # server → minitor
go build -o minitor-client ./cmd/client
```

## Development

```bash
task lint
task test
```

CI runs on push/PR to `main`: lint, build, and test.

## Architecture

Want the full picture? Check out [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md).

## Project layout

```
minitor/
├── main.go              # server entry point
├── cmd/client/          # TUI client entry point
├── config/              # config loading
├── collector/           # OS metric collectors
├── transport/           # HTTP + WebSocket server
│   ├── http/
│   └── socket/
├── view/terminal/       # TUI (Bubble Tea)
└── helper/              # small utils
```
