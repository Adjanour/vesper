# Vesper (Time Block Planner)

[![Go Version](https://img.shields.io/badge/Go-1.24%2B-00ADD8?logo=go)](https://golang.org/dl/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)

Vesper is a lightweight backend service for a **time-block planning workflow**.
Each night, the system (planned) will send you a link to a calendar UI where you can create time blocks.
Those blocks are then synced to your calendar (**Google Calendar integration planned**).

This repository contains the Go backend that stores and manages time blocks (tasks) and exposes a simple HTTP API.

---

## 📚 Documentation

- **[Installation Guide](INSTALL.md)** - Detailed setup instructions
- **[API Documentation](API.md)** - Complete API reference with examples
- **[Contributing Guide](CONTRIBUTING.md)** - How to contribute to the project
- **[Changelog](CHANGELOG.md)** - Version history and changes

---

## Table of Contents

* [Project Overview](#project-overview)
* [What Is Implemented Today](#what-is-implemented-today)
* [Quick Start](#quick-start)
* [Development](#development)
* [Docker](#docker)
* [Roadmap / Next Steps](#roadmap--next-steps)
* [Contributing](#contributing)
* [License](#license)

---

## Project Overview

Vesper’s goal is to make **nightly time-block planning frictionless**:

1. Each night, the user receives a link to a planning page for the next day.
2. The user arranges time blocks in a calendar-like UI.
3. The selected blocks are synced to the user’s primary calendar (Google Calendar).

This repo implements the **backend storage and API** for tasks (time blocks).
It is intentionally small and opinionated so the frontend and integrations can evolve separately.


## What Is Implemented Today

✅ **Features implemented:**

* HTTP server that listens on `:8080` and exposes a JSON API
* SQLite-based persistence stored at `./data/tasks.db`
* Core task operations:
  * Create a task (with overlap check)
  * Get a single task
  * Delete a task
* Database migrations with automated migration runner
* Simple CORS setup for browser-based UIs
* Docker support with multi-stage builds
* Comprehensive documentation and setup guides

🚧 **Not yet implemented:**

* Google Calendar OAuth + sync
* Nightly email with planning link
* Frontend planning UI
* Authentication & multi-user support (basic `user_id` exists but no auth)
* Background worker / nightly scheduler

---

## Quick Start

### Prerequisites

- Go 1.24 or higher
- SQLite (embedded via Go driver - no separate installation needed)

### Installation

```bash
# Clone the repository
git clone https://github.com/Adjanour/vesper.git
cd vesper

# Complete setup (install dependencies + run migrations)
make setup

# Run the server
make run
```

The server starts on **port 8080**. Test it:

```bash
curl http://localhost:8080/api/health
```

**For detailed installation instructions, see [INSTALL.md](INSTALL.md).**

### Quick API Test

```bash
# Create a task
curl -X POST http://localhost:8080/api/tasks/ \
  -H "Content-Type: application/json" \
  -d '{
    "id": "task-001",
    "title": "Morning Review",
    "start": "2026-02-08T09:00:00Z",
    "end": "2026-02-08T10:00:00Z",
    "user_id": "1",
    "status": "scheduled"
  }'

# Get the task
curl http://localhost:8080/api/tasks/task-001
```

**For complete API documentation, see [API.md](API.md).**

---

## Development

### Available Make Commands

```bash
make help          # Show all available commands
make build         # Build the binary
make run           # Build and run
make clean         # Clean build artifacts
make test          # Run tests
make migrate       # Run database migrations
make migrate-down  # Rollback migrations
make dev           # Run with hot reload (requires Air)
make fmt           # Format code
make lint          # Run linters
make setup         # Complete project setup
```

### Development Mode (Hot Reload)

Install [Air](https://github.com/cosmtrek/air) for automatic reloading:

```bash
go install github.com/cosmtrek/air@latest
make dev
```

**For contributing guidelines, see [CONTRIBUTING.md](CONTRIBUTING.md).**

---

## Docker

### Build and Run

```bash
# Using Make
make docker-build
make docker-run

# Or using Docker directly
docker build -t vesper:latest .
docker run -p 8080:8080 -v $(pwd)/data:/data vesper:latest
```

The Docker image:
- Uses multi-stage builds for a small image size (~20MB)
- Runs as a non-root user
- Mounts `/data` as a volume for database persistence
- Exposes port 8080

---

## Roadmap / Next Steps

### Short Term

* [ ] Google OAuth 2.0 integration & background sync job
* [ ] Nightly job sending users planning links via email
* [ ] Browser-based planning UI consuming `/api`
* [ ] Authentication & true multi-user separation
* [ ] Comprehensive test suite

### Medium Term

* [ ] Recurring blocks & conflict resolution suggestions
* [ ] iCal import/export support
* [ ] Rate limiting, logging, and metrics
* [ ] CI/CD pipeline with automated tests
* [ ] API versioning

### Long Term

* [ ] Team scheduling & shared calendars
* [ ] Advanced conflict resolution heuristics
* [ ] Mobile applications (iOS/Android)
* [ ] Third-party calendar integrations (Outlook, Apple Calendar)

---

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details on:

- Development setup
- Coding standards
- Pull request process
- Reporting bugs and suggesting features

Quick start for contributors:

```bash
git checkout -b feat/your-feature
# Make your changes
make test
make lint
git commit -m "feat: add your feature"
git push origin feat/your-feature
```

---

## Project Structure

```
vesper/
├── cmd/
│   └── server/              # Main application entry point
├── internal/
│   ├── api/                 # HTTP handlers and routing
│   ├── database/           # Database operations and migrations
│   │   └── migrate/        # Migration runner
│   └── models/             # Data models
├── data/                   # SQLite database storage (gitignored)
├── API.md                  # API documentation
├── CONTRIBUTING.md         # Contributing guidelines
├── INSTALL.md              # Installation guide
├── CHANGELOG.md            # Version history
├── Makefile               # Build and development tasks
├── Dockerfile             # Docker configuration
└── README.md              # This file
```

---

## License

MIT License - see [LICENSE](LICENSE) file for details.

---

## Author

Created by [Adjanour](https://github.com/Adjanour)

For questions about design decisions and roadmap, please open an issue or discussion.

---

**Built with ❤️ using Go**
