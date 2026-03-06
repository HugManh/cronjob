<div align="center">

# 🕒 CronJob Manager

**A robust, production-ready Go service for managing, monitoring, and automating scheduled tasks.**

[![CI](https://github.com/HugManh/cronjob/actions/workflows/ci.yml/badge.svg)](https://github.com/HugManh/cronjob/actions/workflows/ci.yml)
[![Docker](https://github.com/HugManh/cronjob/actions/workflows/docker.yml/badge.svg)](https://github.com/HugManh/cronjob/actions/workflows/docker.yml)
[![Go Version](https://img.shields.io/badge/Go-1.24-00ADD8?logo=go)](go.mod)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

</div>

---

## ✨ Features

| Feature | Description |
|---|---|
| 🗓️ **Task Scheduling** | Create and manage cron jobs with full CRUD support |
| ▶️ **Live Controls** | Enable / disable tasks instantly via UI or API |
| 🔔 **Slack Alerts** | Automated notifications to your Slack workspace |
| 🗄️ **PostgreSQL** | Persistent task storage powered by GORM |
| 🌐 **Web Dashboard** | Clean Vanilla JS UI for real-time monitoring |
| 🐳 **Docker Ready** | Multi-stage image, docker compose for local setup |

---

## 🛠️ Tech Stack

- **Language**: Go 1.24
- **Web Framework**: [Gin](https://github.com/gin-gonic/gin)
- **ORM**: [GORM](https://gorm.io) + PostgreSQL 16
- **Scheduler**: [robfig/cron v3](https://github.com/robfig/cron)
- **Notifications**: [slack-go/slack](https://github.com/slack-go/slack)
- **Logging**: [Logrus](https://github.com/sirupsen/logrus)

---

## 🚀 Quick Start (Docker)

> **Recommended:** The fastest way to run the full stack locally.

**1. Copy the environment template:**
```bash
cp .env.template .env
# Edit .env and fill in any values you want to override
```

**2. Start the app and database:**
```bash
docker compose up -d
```

**3. Open the dashboard:**
```
http://localhost:8080
```

**Stop everything:**
```bash
docker compose down
```

---

## 🔧 Local Development

### Prerequisites
- [Go 1.24+](https://go.dev/dl/)
- [PostgreSQL 16+](https://www.postgresql.org/download/)

### Setup
```bash
# 1. Clone the repository
git clone https://github.com/HugManh/cronjob.git && cd cronjob

# 2. Copy environment config
cp .env.template .env

# 3. Download dependencies
go mod tidy

# 4. Run the server
make run
# or: go run ./cmd/api
```

---

## ⚙️ Configuration

All configuration is done via environment variables (loaded from `.env`):

| Variable | Description | Default |
|---|---|---|
| `SERVER_PORT` | HTTP server port | `8080` |
| `ENVIRONMENT` | Runtime env (`local`, `production`) | `production` |
| `LOG_LEVEL` | Log verbosity (`debug`, `info`, `error`) | `debug` |
| `DB_HOST` | PostgreSQL host | `localhost` |
| `DB_PORT` | PostgreSQL port | `5432` |
| `DB_USER` | Database username | `postgres` |
| `DB_PASSWORD` | Database password | *(required)* |
| `DB_DATABASE` | Database name | `postgres` |
| `DB_SSL` | Enable SSL for DB connection | `false` |

---

## 🔌 API Reference

Base URL: `/api/v1`

### Tasks

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/tasks` | List all tasks |
| `POST` | `/tasks` | Create a new task |
| `GET` | `/tasks/:id` | Get task by ID |
| `PUT` | `/tasks/:id` | Update a task |
| `DELETE` | `/tasks/:id` | Delete a task |
| `POST` | `/tasks/:id/active` | Toggle task enabled/disabled |

### Slack Channels

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/slacks` | List Slack integrations |
| `POST` | `/slacks` | Register a Slack webhook |
| `GET` | `/slacks/:id` | Get Slack channel by ID |

---

## 🏗️ Build

```bash
# Build native binary
make build            # Output: bin/cronjob

# Cross-compile for all platforms
make build-all        # Output: bin/cronjob-<os>-<arch>

# Build & tag Docker image
make docker-build

# View all available targets
make help
```

---

## 🔁 CI / CD

| Workflow | Trigger | What it does |
|---|---|---|
| **CI** (`ci.yml`) | Push / PR → `main` | `go vet`, `golangci-lint`, `go test -race` |
| **Docker** (`docker.yml`) | Push → `main`, semver tag | Build & push image to `ghcr.io` |
| **Release** (`release.yml`) | Push tag `v*.*.*` | Build cross-platform binaries, create GitHub Release |

Dependency updates are automated via [Dependabot](.github/dependabot.yml) for Go modules, GitHub Actions, and Docker base images.

---

## 📂 Project Structure

```text
.
├── .github/
│   ├── workflows/         # CI/CD pipelines (ci, docker, release)
│   └── ISSUE_TEMPLATE/    # Bug report & feature request templates
├── cmd/api/               # Application entry point (main.go)
├── configs/               # Environment config loading
├── internal/
│   ├── dto/               # Data Transfer Objects
│   ├── handler/           # HTTP handlers (Gin)
│   ├── model/             # GORM database models
│   ├── repository/        # Database layer
│   ├── routing/           # Route registration
│   ├── service/           # Business logic & cron management
│   ├── startup/           # Server initialization
│   └── web/               # View-specific logic
├── pkg/
│   ├── db/postgres/       # Database connection & helpers
│   └── logger/            # Logrus logger setup
├── public/                # Static frontend (HTML, CSS, JS)
├── views/                 # Server-side templates
├── Dockerfile             # Multi-stage Docker build
├── docker-compose.yml     # Local dev compose (app + db)
├── Makefile               # Build, test, lint, docker targets
└── .env.template          # Environment variable template
```

---

## 📝 License

Distributed under the **MIT License**. See [LICENSE](LICENSE) for details.
