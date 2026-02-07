# Contributing to Vesper

Thank you for your interest in contributing to Vesper! This document provides guidelines and instructions for contributing to the project.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [How to Contribute](#how-to-contribute)
- [Development Workflow](#development-workflow)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Pull Request Process](#pull-request-process)
- [Reporting Bugs](#reporting-bugs)
- [Suggesting Features](#suggesting-features)

## Code of Conduct

We expect all contributors to be respectful and professional. Please:

- Be welcoming and inclusive
- Be respectful of differing viewpoints and experiences
- Accept constructive criticism gracefully
- Focus on what is best for the community

## Getting Started

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/vesper.git
   cd vesper
   ```
3. **Add the upstream repository**:
   ```bash
   git remote add upstream https://github.com/Adjanour/vesper.git
   ```

## Development Setup

### Prerequisites

- Go 1.24 or higher
- SQLite (embedded via modernc.org/sqlite)
- Git
- Optional: [Air](https://github.com/cosmtrek/air) for hot reload during development

### Initial Setup

1. **Install dependencies**:
   ```bash
   make install
   ```

2. **Run database migrations**:
   ```bash
   make migrate
   ```

3. **Build the project**:
   ```bash
   make build
   ```

4. **Run the server**:
   ```bash
   make run
   ```

The server will start on `http://localhost:8080`.

### Development Mode (Hot Reload)

For development with automatic reloading:

```bash
# Install Air (if not already installed)
go install github.com/cosmtrek/air@latest

# Run in dev mode
make dev
```

## How to Contribute

There are many ways to contribute to Vesper:

- **Bug fixes**: Fix existing issues
- **New features**: Implement features from the roadmap
- **Documentation**: Improve or add documentation
- **Tests**: Add test coverage
- **Code quality**: Refactor and improve code quality
- **Examples**: Add example code or tutorials

## Development Workflow

1. **Create a feature branch**:
   ```bash
   git checkout -b feat/your-feature-name
   # or
   git checkout -b fix/your-bug-fix
   ```

   Use prefixes:
   - `feat/` for new features
   - `fix/` for bug fixes
   - `docs/` for documentation
   - `refactor/` for code refactoring
   - `test/` for adding tests

2. **Make your changes**:
   - Write clear, concise code
   - Follow the coding standards (see below)
   - Add tests for new functionality
   - Update documentation as needed

3. **Test your changes**:
   ```bash
   make test
   make lint
   ```

4. **Commit your changes**:
   ```bash
   git add .
   git commit -m "feat: add your feature description"
   ```

   Follow [Conventional Commits](https://www.conventionalcommits.org/):
   - `feat:` for new features
   - `fix:` for bug fixes
   - `docs:` for documentation changes
   - `refactor:` for code refactoring
   - `test:` for test additions

5. **Keep your branch up to date**:
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

6. **Push to your fork**:
   ```bash
   git push origin feat/your-feature-name
   ```

7. **Open a Pull Request** on GitHub

## Coding Standards

### Go Code Style

- Follow standard Go conventions and idioms
- Use `gofmt` for code formatting (automatically done by `make fmt`)
- Use meaningful variable and function names
- Keep functions small and focused
- Add comments for exported functions and complex logic

### Code Organization

- Place new models in `internal/models/`
- Place API handlers in `internal/api/`
- Place database operations in `internal/database/`
- Keep business logic separate from HTTP handlers

### Error Handling

- Use domain-specific errors defined in `internal/database/db.go`
- Return appropriate HTTP status codes
- Provide clear error messages

### Example

```go
// Good
func (q *Queries) GetTask(ctx context.Context, id string) (*models.Task, error) {
	row := q.db.QueryRowContext(ctx, getTaskSQL, id)
	
	var t models.Task
	err := row.Scan(&t.ID, &t.Title, &t.Start, &t.End, &t.Status, &t.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &t, nil
}
```

## Testing

While the project is in early stages, we encourage adding tests:

```bash
# Run all tests
make test

# Run tests with coverage
go test -cover ./...

# Run tests for a specific package
go test -v ./internal/database/
```

## Pull Request Process

1. **Ensure all tests pass** and code is formatted
2. **Update the README.md** if you've added new features or changed behavior
3. **Update API.md** if you've modified API endpoints
4. **Add or update tests** for your changes
5. **Write a clear PR description**:
   - Explain what changes you made
   - Why you made them
   - How to test them
   - Link to any related issues

6. **Request a review** from maintainers
7. **Address review feedback** promptly
8. Once approved, a maintainer will merge your PR

### PR Title Format

Use conventional commit format in PR titles:

- `feat: add Google Calendar sync`
- `fix: correct task overlap detection`
- `docs: improve API documentation`

## Reporting Bugs

Found a bug? Please create an issue with:

1. **Clear title** describing the bug
2. **Steps to reproduce** the issue
3. **Expected behavior**
4. **Actual behavior**
5. **Environment details** (OS, Go version, etc.)
6. **Any relevant logs or error messages**

## Suggesting Features

Have an idea for a new feature?

1. **Check existing issues** to avoid duplicates
2. **Open a new issue** with:
   - Clear description of the feature
   - Use cases and benefits
   - Possible implementation approach (optional)
3. **Wait for feedback** from maintainers before starting work

## Project Structure

```
vesper/
├── cmd/
│   └── server/          # Main application entry point
├── internal/
│   ├── api/            # HTTP handlers and routing
│   ├── database/       # Database operations and migrations
│   └── models/         # Data models
├── data/               # SQLite database storage (gitignored)
├── API.md              # API documentation
├── CONTRIBUTING.md     # This file
├── Makefile            # Build and development tasks
└── README.md           # Project overview
```

## Getting Help

- **Questions?** Open a discussion on GitHub
- **Stuck?** Comment on your PR or issue
- **Need clarification?** Reach out to the maintainers

## Recognition

All contributors will be recognized in the project. Thank you for helping make Vesper better!

---

**Happy Contributing! 🚀**
