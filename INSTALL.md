# Installation Guide

This guide provides detailed instructions for setting up and running Vesper on your local machine or server.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Detailed Installation](#detailed-installation)
- [Running the Application](#running-the-application)
- [Docker Installation](#docker-installation)
- [Development Setup](#development-setup)
- [Troubleshooting](#troubleshooting)

## Prerequisites

Before installing Vesper, ensure you have the following:

- **Go 1.24 or higher** - [Download Go](https://golang.org/dl/)
- **Git** - [Install Git](https://git-scm.com/downloads)
- **SQLite** - Included via embedded Go driver (no separate installation needed)

### Verify Prerequisites

```bash
# Check Go version
go version
# Should output: go version go1.24.x or higher

# Check Git
git --version
```

## Quick Start

For the impatient, here's the fastest way to get started:

```bash
# Clone the repository
git clone https://github.com/Adjanour/vesper.git
cd vesper

# Set up and run
make setup
make run
```

The server will be available at `http://localhost:8080`.

## Detailed Installation

### Step 1: Clone the Repository

```bash
git clone https://github.com/Adjanour/vesper.git
cd vesper
```

### Step 2: Install Dependencies

```bash
make install
```

This will download all required Go modules.

### Step 3: Configure the Application (Optional)

Create a `.env` file for custom configuration:

```bash
cp .env.example .env
```

Edit `.env` to customize settings:

```env
PORT=8080
DATA_DIR=./data
DATABASE_PATH=./data/tasks.db
```

### Step 4: Run Database Migrations

Create the database and apply migrations:

```bash
make migrate
```

This creates the SQLite database at `./data/tasks.db` with the necessary tables.

### Step 5: Build the Application

```bash
make build
```

This creates a binary named `vesper` in the project root.

### Step 6: Run the Application

```bash
make run
```

Or run the binary directly:

```bash
./vesper
```

The server will start and listen on port 8080. You should see:

```
2026/02/07 19:17:00 Connected to database
2026/02/07 19:17:00 Server starting on :8080
```

### Step 7: Verify Installation

Test the health endpoint:

```bash
curl http://localhost:8080/api/health
```

Expected response:
```json
{"status":"ok"}
```

## Running the Application

### Using Make (Recommended)

```bash
# Build and run
make run

# Run in development mode with hot reload (requires Air)
make dev

# Run with custom options
PORT=9000 ./vesper
```

### Direct Execution

```bash
# Build first
go build -o vesper ./cmd/server

# Run
./vesper
```

### Background Execution

```bash
# Run in background (Linux/macOS)
./vesper &

# Or use nohup
nohup ./vesper > vesper.log 2>&1 &

# Check if running
ps aux | grep vesper
```

## Docker Installation

### Build and Run with Docker

```bash
# Build the Docker image
make docker-build

# Run the container
make docker-run
```

Or use Docker commands directly:

```bash
# Build image
docker build -t vesper:latest .

# Run container
docker run -p 8080:8080 -v $(pwd)/data:/data vesper:latest
```

### Docker Compose (Optional)

Create a `docker-compose.yml`:

```yaml
version: '3.8'

services:
  vesper:
    build: .
    ports:
      - "8080:8080"
    volumes:
      - ./data:/data
    environment:
      - PORT=8080
    restart: unless-stopped
```

Run with:

```bash
docker-compose up -d
```

## Development Setup

### Install Development Tools

1. **Air** (for hot reload):
   ```bash
   go install github.com/cosmtrek/air@latest
   ```

2. **golangci-lint** (optional, for linting):
   ```bash
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
   ```

### Run in Development Mode

```bash
make dev
```

This will:
- Watch for file changes
- Automatically rebuild and restart the server
- Show build errors in real-time

### Available Make Commands

```bash
make help          # Show all available commands
make build         # Build the binary
make run           # Build and run
make clean         # Clean build artifacts
make test          # Run tests
make migrate       # Run database migrations
make migrate-down  # Rollback migrations
make fmt           # Format code
make vet           # Run go vet
make lint          # Run all linters
make setup         # Complete project setup
```

## Troubleshooting

### Port Already in Use

If port 8080 is already in use:

```bash
# Find the process using port 8080
lsof -i :8080  # macOS/Linux
netstat -ano | findstr :8080  # Windows

# Use a different port
PORT=9000 ./vesper
```

### Database Locked Error

If you get a "database is locked" error:

```bash
# Stop all running instances
pkill vesper

# Remove database lock files
rm -f ./data/tasks.db-shm ./data/tasks.db-wal

# Restart the application
make run
```

### Build Errors

If you encounter build errors:

```bash
# Clean and rebuild
make clean
go clean -cache
make install
make build
```

### Migration Failures

If migrations fail:

```bash
# Remove the database and start fresh
rm -f ./data/tasks.db

# Run migrations again
make migrate
```

### Permission Denied

If you get permission errors:

```bash
# Make the binary executable
chmod +x vesper

# Ensure data directory has proper permissions
chmod 755 data
```

## Platform-Specific Notes

### Windows

- Use PowerShell or Git Bash for running commands
- The `.air.toml` file has Windows-specific paths (already configured)
- Use `./vesper.exe` instead of `./vesper`

### macOS

- Ensure Go is added to your PATH
- You may need to allow the binary in Security & Privacy settings

### Linux

- Works out of the box on most distributions
- For systemd service setup, see below

## Running as a System Service (Linux)

Create a systemd service file `/etc/systemd/system/vesper.service`:

```ini
[Unit]
Description=Vesper Time Block Planner
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/vesper
ExecStart=/opt/vesper/vesper
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
```

Enable and start:

```bash
sudo systemctl daemon-reload
sudo systemctl enable vesper
sudo systemctl start vesper
sudo systemctl status vesper
```

## Next Steps

After installation:

1. **Read the [API Documentation](API.md)** to learn about available endpoints
2. **Try the example workflow** in the API docs
3. **Check out the [Contributing Guide](CONTRIBUTING.md)** if you want to contribute
4. **Star the repository** if you find it useful!

## Getting Help

If you encounter issues not covered here:

- Check existing [GitHub Issues](https://github.com/Adjanour/vesper/issues)
- Open a new issue with details about your problem
- Include your OS, Go version, and error messages

---

**Installation complete! Start building your time blocks! 🎯**
