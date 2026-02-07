# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Complete project setup and documentation
- `.gitignore` file to exclude build artifacts and dependencies
- `Makefile` with common development tasks
- Database migration runner (`internal/database/migrate/migrate.go`)
- Comprehensive API documentation (`API.md`)
- Contributing guide (`CONTRIBUTING.md`)
- Installation guide (`INSTALL.md`)
- Example environment configuration (`.env.example`)
- This changelog file

### Changed
- Updated README.md with references to new documentation files

## [0.1.0] - 2025-11-02

### Added
- Initial release of Vesper backend
- HTTP API server with chi router
- SQLite database integration
- Task creation with overlap detection
- Task retrieval by ID
- Task deletion
- Health check endpoint
- CORS support for browser-based clients
- Database migrations for tasks and users tables
- Basic error handling with domain-specific errors
- Docker support with multi-stage builds
- Air configuration for hot reload during development

### Core Features
- Time block (task) management
- Overlap prevention for scheduled tasks
- Task status tracking (scheduled, deleted, replaced)
- User association for tasks

[Unreleased]: https://github.com/Adjanour/vesper/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/Adjanour/vesper/releases/tag/v0.1.0
