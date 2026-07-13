# Minitor Architecture

## The big picture

Minitor splits things into three jobs: **collect data**, **ship it**, **show it**. The server owns collection and broadcasting. The TUI client is just a consumer — it doesn't touch the OS directly.

```
┌─────────────┐     WebSocket      ┌──────────────┐
│  TUI Client │ ◄──────────────── │    Server    │
│ view/terminal│                   │  transport/  │
└─────────────┘                    └──────┬───────┘
                                          │
                                   ┌──────▼───────┐
                                   │  collector/  │
                                   │  (gopsutil)  │
                                   └──────────────┘
```

## Layers

### 1. Collector (`collector/`)

Reads metrics from the OS. Each metric has its own worker on a ticker:

| Worker | Interval | Output |
|--------|----------|--------|
| `WorkerCpu` | 1s | `CpuMetric` |
| `WorkerRam` | 1s | `RamMetric` |
| `WorkerDisk` | 3s | `DiskMetric` |
| `WorkerNetwork` | 3s | `NetworkMetric` |
| `ProcessWorker` | 3s | `[]*ProcessMetric` |

Processes are built as a parent/child tree first, then flattened when sent over the API.

### 2. Transport (`transport/`)

#### HTTP (`transport/http/`)

- `http.Server` with graceful shutdown
- Routes: `/health` and `/ws`
- Starts the `Monitor` and wires it to `socket.Handler`

#### Socket (`transport/socket/`)

| Component | What it does |
|-----------|--------------|
| `Monitor` | runs collector workers, builds snapshots, broadcasts |
| `Hub` | tracks connected clients + keeps the last message |
| `Handler` | upgrades HTTP to WebSocket |
| `Client` | read/write loop per connection |
| `Router` | dispatches incoming client messages (`ping`, `processes`) |
| `Snapshot` | JSON model for broadcast metrics |
| `ProcessesPage` | paginated process list |

### 3. View (`view/terminal/`)

The TUI, built with Bubble Tea. It does **not** talk to collectors directly — only WebSocket:

| File | Role |
|------|------|
| `model.go` | state and init |
| `update.go` | message and keyboard handling |
| `view.go` | renders the UI |
| `ws.go` | WebSocket connection, snapshots, process fetching |

### 4. Config (`config/`)

Server, client, and socket settings. Loaded from JSON, env vars, and CLI flags, with validation.

## Data flow

### Broadcasting metrics

```
collector workers
       │
       ▼
   Monitor (select on channels)
       │
       ├── update Snapshot
       └── Hub.Broadcast()
                │
                ▼
         all Client.Send channels
                │
                ▼
         Client.WritingLoop → WebSocket
```

### New client connects

1. `Handler.ServeHTTP` accepts the WebSocket
2. Client gets registered in the `Hub`
3. Last `snapshot` is sent right away
4. `Client.Run` starts the read/write loop

### Process requests

The full process list isn't in every snapshot — it's too big. Clients request pages instead:

```
Client ──{"type":"processes","data":{"offset":0,"limit":50}}──► Router
                                                                    │
                                                                    ▼
                                                              Monitor.ProcessesPage()
                                                                    │
Client ◄──{"type":"processes","data":{...}}────────────────────────┘
```

The TUI client fetches all pages and rebuilds the tree so `view.go` doesn't need to know about pagination.

## WebSocket protocol

Shared envelope:

```json
{ "type": "string", "data": { ... } }
```

| type | direction | what |
|------|-----------|------|
| `snapshot` | server → client | CPU, RAM, Disk, Network + process_count |
| `processes` | both ways | paginated process list request/response |
| `ping` | client → server | health check |
| `pong` | server → client | ping reply |

## Graceful shutdown

```
SIGINT / SIGTERM
      │
      ▼
signal.NotifyContext (main)
      │
      ▼
http.Server.Shutdown (timeout from config)
      │
      ├── stop accepting connections
      ├── cancel active request contexts
      └── close WebSocket connections
      │
      ▼
Monitor.Run → ctx.Done() → exit
```

## Entry points

| Binary | Path | Role |
|--------|------|------|
| Server | `main.go` | collectors + HTTP/WS |
| Client | `cmd/client/main.go` | TUI |

## Key dependencies

| Package | Used for |
|---------|----------|
| `gopsutil` | reading OS metrics |
| `coder/websocket` | WebSocket |
| `bubbletea` + `lipgloss` | TUI |

## Why it's built this way

**Why client/server?**
Keeps collection and display separate. Multiple clients can hook into one server.

**Why paginate processes?**
There are a lot of them. Sending the full list in every snapshot makes messages huge and connections drop.

**Why is the TUI separate from collectors?**
`view/terminal` is just a consumer. You can change collectors or transport without touching the UI — only the WebSocket protocol matters.

**Why a separate Router?**
Each incoming message type gets its own handler. Adding a new route doesn't mess with connection lifecycle.
