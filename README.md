# Vesper (Time Block Planner)

Vesper is a lightweight backend service for a **time-block planning workflow**.
Each night, the system (planned) will send you a link to a calendar UI where you can create time blocks.
Those blocks are then synced to your calendar (**Google Calendar integration planned**).

This repository contains the Go backend that stores and manages time blocks (tasks) and exposes a simple HTTP API.


**Status:** Early stage 🚧

---

## Table of Contents

* [Project Overview](#project-overview)
* [What Is Implemented Today](#what-is-implemented-today)
* [Quick Start (Build & Run)](#quick-start--build--run)
* [HTTP API (Endpoints + Examples)](#http-api)
* [Data Model & Database](#data-model--database)
* [Docker (Suggested Dockerfile)](#docker-suggested-dockerfile)
* [Roadmap / Next Steps](#roadmap--next-steps)
* [Contributing](#contributing)
* [Notes and Implementation Details](#notes-and-implementation-details)
* [Contacts / Author](#contacts--author)
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

✅ Features implemented:

* HTTP server that listens on `:8080` and exposes a small JSON API.
* SQLite-based persistence stored at `./data/tasks.db`.
* Core task operations:

  * Create a task (with overlap check)
  * Get a single task
  * Delete a task
* Database migration files under `internal/database/migrations`.
* Simple CORS setup for browser-based UIs.

🚧 Not yet implemented:

* Google Calendar OAuth + sync
* Nightly email with planning link
* Frontend planning UI
* Authentication & multi-user support (basic `user_id` exists but no auth)
* Background worker / nightly scheduler


## Quick Start — Build & Run

### Prerequisites

* Go **1.24+**
* SQLite (handled via embedded Go driver `modernc.org/sqlite` — no external DB required)

### Build Locally

```bash
go build -o vesper ./cmd/server
```

### Run

```bash
# Run directly
./vesper

# Or with go run
go run ./cmd/server
```

The server starts on **port 8080**:

* Health check → [http://localhost:8080/api/health](http://localhost:8080/api/health)
* Task endpoints → `/api/tasks`


## HTTP API

**Base URL:** `http://localhost:8080/api`

### Health Check

**GET** `/api/health`

### Create Task

**POST** `/api/tasks/`
**Body:** JSON `Task` object (see [Data Model](#data-model--database))

If a conflicting block exists, returns **409 Conflict**.

### Get Task

**GET** `/api/tasks/{id}`

### Delete Task

**DELETE** `/api/tasks/{id}`


### Notes on the Current API

* `GetTasks` currently uses a **hardcoded `user_id = "1"`**.
  → Multi-user auth support is planned.
* Standard error codes:

  * `400` — Invalid JSON
  * `404` — Task not found
  * `409` — Task overlaps with existing one

## Roadmap / Next Steps

### Short Term

* [ ] Google OAuth 2.0 integration & background sync job
* [ ] Nightly job sending users planning links via email
* [ ] Browser-based planning UI consuming `/api`
* [ ] Authentication & true multi-user separation

### Medium Term

* [ ] Recurring blocks & conflict resolution suggestions
* [ ] iCal import/export support
* [ ] Rate limiting, logging, and metrics
* [ ] CI with migrations + tests

### Long Term

* [ ] Team scheduling & shared calendars
* [ ] Advanced conflict resolution heuristics


## Contributing

* Fork the repo & open a PR for new features or fixes.
* Add tests when possible.

Suggested workflow:

```bash
git checkout -b feat/your-feature
# implement feature
git commit -m "feat: add your-feature"
git push origin feat/your-feature
```

Then open a PR with a clear description.


## Notes and Implementation Details

* Codebase is **small and idiomatic Go**.
* Uses [`chi`](https://github.com/go-chi/chi) for routing.
* Uses [`modernc.org/sqlite`](https://pkg.go.dev/modernc.org/sqlite) as an embedded driver.
* Database operations live in `internal/database/` and return domain-level errors (`ErrTaskOverlap`, `ErrNotFound`, etc.).
* A `/tmp` directory may hold build artifacts — you can delete it for a clean repo.


## Contacts / Author

See the [repository owner](https://github.com/Adjanour) for questions about design decisions and roadmap.


## License
MIT License (Beerware Edition)
