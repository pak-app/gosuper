# Gosuper 🚀

**Gosuper** is a lightweight, high-performance process supervisor and daemon written in Go. Inspired by tools like `supervisord` and `PM2`, Gosuper allows you to manage, monitor, and control background services using YAML configuration files over a Unix Domain Socket HTTP API.

---

## 📋 Table of Contents
- [Architecture & Design](#-architecture--design)
- [Installation & Building](#-installation--building)
- [Configuration (`gosuper.yaml`)](#-configuration-gosuperyaml)
- [CLI Command Usage](#-cli-command-usage)
- [Unix Socket HTTP API](#-unix-socket-http-api)
- [Project Status & Roadmap](#-project-status--roadmap)
- [Development & Testing](#-development--testing)

---

## 🏗 Architecture & Design

Gosuper is structured into decoupled modules:

```
gosuper/
├── cmd/                # Cobra CLI command definitions (daemon, service, supervisor)
├── internal/
│   ├── client/         # Go HTTP client communicating over Unix Domain Socket
│   ├── config/         # YAML configuration loader & types
│   ├── core/           # Process management engine, Service & Supervisor state machines
│   ├── logging/        # Multi-writer logging setup (stdout + log files)
│   ├── server/         # HTTP Daemon server listening on Unix socket (tmp/gosuper.sock)
│   └── types/          # Shared API request/response types
├── gosuper.yaml        # Default configuration file
└── main.go             # Application entry point
```

---

## 🔨 Installation & Building

### Prerequisites
- **Go**: 1.22 or higher

### Build Binary
```bash
# Clone repository
git clone https://github.com/pak-app/gosuper.git
cd gosuper

# Build binary
go build -o gosuper main.go
```

---

## ⚙️ Configuration (`gosuper.yaml`)

Gosuper uses YAML configuration files to define supervisor settings and managed service specifications.

```yaml
# gosuper.yaml

# Global supervisor configuration
supervisor:
  name: 'dev-supervisor'
  log_dir: './log/gosuper'
  restart_delay: 2s
  stop_timeout: 5s
  env:
    APP_ENV: development

# Defined services managed by the supervisor
services:
  dummy_worker:
    command: ["./dummy.sh"]
    dir: "./tmp"
    autostart: true
    autorestart: true
    restart_limit: 5
    restart_window: 30s
    stdout: "worker.out.log"
    stderr: "worker.err.log"
    env:
      DEBUG: "true"

  dummy_api:
    command: ["./api_server"]
    dir: "./tmp"
    autostart: true
    autorestart: true
    stdout: "api.out.log"
    stderr: "api.err.log"
```

---

## 💻 CLI Command Usage

The `gosuper` CLI provides commands to manage the daemon and running services.

### Daemon Management
- **Start Daemon in Background**:
  ```bash
  ./gosuper daemon start
  ```
- **Run Daemon in Foreground (Serve)**:
  ```bash
  ./gosuper daemon serve
  ```
- **Check Daemon Health**:
  ```bash
  ./gosuper daemon status
  ```
- **Stop Daemon**:
  ```bash
  ./gosuper daemon stop
  ```

### Service & Supervisor Management
- **Start Services (from config)**:
  ```bash
  ./gosuper service start --config gosuper.yaml
  ```
- **Check Service & Supervisor Status**:
  ```bash
  ./gosuper service status --supervisor-name dev-supervisor
  ```
- **Stop Services**:
  ```bash
  ./gosuper service stop --supervisor-name dev-supervisor
  ```

---

## 🔌 Unix Socket HTTP API

The daemon listens on `tmp/gosuper.sock` by default and exposes RESTful HTTP endpoints:

| Endpoint | Method | Description |
| :--- | :--- | :--- |
| `/daemon/status` | `GET` | Get daemon health status and active supervisor count |
| `/daemon/stop` | `POST` | Gracefully shut down the HTTP server and remove socket |
| `/service/start` | `POST` | Load configuration payload and start supervisor services |
| `/service/status` | `GET` | Get status of all supervisors or filter by `?supervisor_name=` |
| `/service/stop` | `POST` | Stop all services under a supervisor specified by `?supervisor_name=` |

---

## 📊 Project Status & Roadmap

### ✅ Implemented Features
- [x] Concurrent process execution engine (`internal/core`)
- [x] Process state tracking (`Starting`, `Running`, `Stopping`, `Stopped`, `Failed`)
- [x] Uptime calculation & PID tracking
- [x] Unix Domain Socket HTTP server & Client library
- [x] YAML Configuration loading (`internal/config`)
- [x] Cobra CLI interface (`cmd/`)
- [x] Multi-writer logging (`internal/logging`)

### 🚧 In Progress / Planned Features
- [ ] **Auto-restart on Crash**: Implement automatic retry backoff based on `restart_limit` & `restart_window`.
- [ ] **Log File Redirection**: Connect service `stdout`/`stderr` configs to file loggers instead of inheriting daemon stdout.
- [ ] **Environment Variable Merging**: Inject `supervisor.env` and `service.env` into child process execution environments.
- [ ] **Single Service Control**: Support starting/stopping individual services within a supervisor group.
- [ ] **Configuration Schema Validation**: Add validation logic in `internal/config/validation.go`.
- [ ] **Log Streaming API**: Implement `/log` endpoint for real-time log tailing.

---

## 🧪 Development & Testing

Run unit tests across all packages:
```bash
go test ./...
```
