# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

MTG Tracker is a Magic: The Gathering game tracking system with:
- Go backend API server with PostgreSQL database
- Dart HTTP client library
- Firebase authentication
- Docker containerization
- S3 storage integration for file uploads

## Architecture

### Backend (Go)
- **Entry point**: `cmd/server/main.go` - HTTP server with Firebase auth, CORS, and JSON middleware
- **Service layer**: `internal/mtgtracker/` - Business logic and HTTP handlers
- **Repository layer**: `internal/repository/` - Database operations using GORM
- **Middleware**: `internal/middleware/` - Authentication, CORS, JSON handling
- **Models**: Defined in repository layer with GORM tags

### API Structure
The API follows RESTful patterns with versioned endpoints:
- Player endpoints: `/player/v1/*`
- Game endpoints: `/game/v1/*`

### Frontend/Client
- **Dart client**: `dart_client/` - Generated HTTP client with JSON serialization
- **Static files**: `static/` - HTML files served at root

## Common Development Commands

### Go Backend
```bash
# Run the server locally
export POSTGRES_DSN="host=localhost user=postgres password=postgres dbname=mtgtracker port=5432 sslmode=disable"
go run cmd/server/main.go

# Build the server
go build -o server cmd/server/main.go

# Run with Docker Compose (includes PostgreSQL)
docker-compose up

# Build Docker image
docker build -t mtgtracker-go:latest .
```

### Database Setup
```bash
# Start PostgreSQL with Docker
docker run -d --name postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_USER=postgres -e POSTGRES_DB=mtgtracker -p 5432:5432 postgres

# Populate test data
./populate.sh
./populate_events.sh
```

### Dart Client
```bash
# Setup and generate JSON serialization
cd dart_client
dart pub get
dart pub run build_runner build

# Run tests
dart test
```

## Key Configuration

### Environment Variables
- `POSTGRES_DSN`: PostgreSQL connection string (required)
- Firebase service account credentials (for authentication)
- AWS credentials (for S3 storage)

### Dependencies
- Go 1.23.2+
- PostgreSQL database
- Firebase project for authentication
- AWS S3 for file storage

## Testing and API Usage

Use the `populate.sh` script to create test players and games. The script uses HTTPie commands that demonstrate the API structure.

Server runs on port 8080 by default and serves static HTML files from the `static/` directory.