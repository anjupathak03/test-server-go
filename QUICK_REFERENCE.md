# Quick Reference Guide

## Project Overview
A REST API Todo application with MySQL database and Keploy mocking integration.

## Quick Start

### Windows PowerShell
```powershell
# Install dependencies
.\run.ps1 deps

# Run unit tests
.\run.ps1 test-unit

# Start the server
.\run.ps1 run
```

### Cross-Platform
```bash
# Install dependencies
go mod download

# Run tests
go test ./handlers/... ./repository/... -v

# Start server
go run main.go
```

## API Quick Reference

**Base URL:** `http://localhost:8080`

| Endpoint | Method | Description | Body Example |
|----------|--------|-------------|--------------|
| `/api/todos` | GET | Get all todos | - |
| `/api/todos` | POST | Create todo | `{"title":"Task","description":"Details"}` |
| `/api/todos/:id` | GET | Get todo by ID | - |
| `/api/todos/:id` | PUT | Update todo | `{"completed":true}` |
| `/api/todos/:id` | DELETE | Delete todo | - |

## Testing Commands

### Unit Tests (No Database)
```bash
go test ./handlers/... ./repository/... -v
```

### Integration Tests (Requires Database)
```bash
go test ./... -v
```

### Keploy Tests (No Database)
```bash
# 1. Record mocks (needs database)
keploy record -c "go run main.go"
# Make API calls...

# 2. Test with mocks (no database)
keploy test -c "go test ./... -v"
```

## Environment Variables

```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_NAME=todo_db
PORT=8080
```

## Common Tasks

### Create Todo
```bash
curl -X POST http://localhost:8080/api/todos \
  -H "Content-Type: application/json" \
  -d '{"title":"My Task","description":"Task details"}'
```

### List All Todos
```bash
curl http://localhost:8080/api/todos
```

### Update Todo
```bash
curl -X PUT http://localhost:8080/api/todos/1 \
  -H "Content-Type: application/json" \
  -d '{"completed":true}'
```

### Delete Todo
```bash
curl -X DELETE http://localhost:8080/api/todos/1
```

## File Structure

```
test-server/
├── main.go              # Application entry point
├── database/            # Database connection
├── models/              # Data structures
├── repository/          # Database operations
├── handlers/            # HTTP handlers
├── routes/              # Route definitions
├── *_test.go            # Test files
├── README.md            # Full documentation
├── KEPLOY_GUIDE.md      # Keploy usage guide
└── QUICK_REFERENCE.md   # This file
```

## Troubleshooting

**Can't connect to database?**
- Check MySQL is running
- Verify credentials in .env file

**Tests failing?**
- Run `go mod tidy` to ensure dependencies are installed
- Check database connection for integration tests

**Keploy not working?**
- Ensure Keploy CLI is installed
- Record mocks before testing with them

## Documentation

- **Full Documentation:** [README.md](README.md)
- **Keploy Guide:** [KEPLOY_GUIDE.md](KEPLOY_GUIDE.md)
- **API Examples:** See README.md

## Key Commands Summary

```powershell
# Windows PowerShell
.\run.ps1 help          # Show all commands
.\run.ps1 deps          # Install dependencies  
.\run.ps1 test-unit     # Run unit tests
.\run.ps1 keploy-record # Record mocks
.\run.ps1 keploy-test   # Test with mocks
```

```bash
# Linux/Mac
make help               # Show all commands
make deps               # Install dependencies
make test-unit          # Run unit tests
make keploy-record      # Record mocks
make keploy-test        # Test with mocks
```
