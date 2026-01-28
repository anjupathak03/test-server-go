# Todo App with Keploy Mocking

A simple REST API Todo application built with Go and MySQL, featuring automated testing with Keploy for mocking database calls.

## Features

- RESTful API for Todo CRUD operations
- MySQL database integration
- Unit tests with standard Go testing
- Integration tests with Keploy mocking
- Clean architecture with separated layers (handlers, repository, models)

## Prerequisites

- Go 1.21 or higher
- MySQL 5.7 or higher
- Keploy CLI (for mocking)

## Project Structure

```
test-server/
├── database/          # Database connection and initialization
├── handlers/          # HTTP request handlers
├── models/            # Data models and DTOs
├── repository/        # Database operations
├── routes/            # Route definitions
├── main.go            # Application entry point
├── main_test.go       # Integration tests
└── *_test.go          # Unit tests
```

## Setup

1. **Install Dependencies**
   ```bash
   go mod download
   ```

2. **Setup MySQL Database**
   ```sql
   CREATE DATABASE todo_db;
   CREATE DATABASE todo_test_db;  -- For testing
   ```

3. **Configure Environment Variables**
   Create a `.env` file or set environment variables:
   ```
   DB_HOST=localhost
   DB_PORT=3306
   DB_USER=root
   DB_PASSWORD=password
   DB_NAME=todo_db
   PORT=8080
   ```

## Running the Application

```bash
go run main.go
```

The server will start on `http://localhost:8080`

## API Endpoints

| Method | Endpoint          | Description        |
|--------|-------------------|--------------------|
| POST   | /api/todos        | Create a new todo  |
| GET    | /api/todos        | Get all todos      |
| GET    | /api/todos/:id    | Get todo by ID     |
| PUT    | /api/todos/:id    | Update a todo      |
| DELETE | /api/todos/:id    | Delete a todo      |

### Example Requests

**Create Todo:**
```bash
curl -X POST http://localhost:8080/api/todos \
  -H "Content-Type: application/json" \
  -d '{"title":"Buy groceries","description":"Milk, eggs, bread"}'
```

**Get All Todos:**
```bash
curl http://localhost:8080/api/todos
```

**Update Todo:**
```bash
curl -X PUT http://localhost:8080/api/todos/1 \
  -H "Content-Type: application/json" \
  -d '{"completed":true}'
```

**Delete Todo:**
```bash
curl -X DELETE http://localhost:8080/api/todos/1
```

## Testing

### Unit Tests (with mocks)

Run unit tests without database:
```bash
go test ./handlers/... -v
go test ./repository/... -v
```

### Integration Tests (with real database)

Run integration tests with real database:
```bash
go test ./... -v
```

## Using Keploy for Mocking

Keploy records all database calls during execution and can replay them during testing, eliminating the need for a real database in tests.

### Step 1: Record Database Calls

First, start the application with Keploy in record mode to capture all database interactions:

```bash
# For Windows PowerShell (recommended approach):
keploy record -c "go run main.go"

# Or if using keploy directly:
keploy record --path ./keploy --command "go run main.go"
```

While recording, make API calls to the application:
```bash
# Create some todos
curl -X POST http://localhost:8080/api/todos -H "Content-Type: application/json" -d '{"title":"Test Todo 1","description":"First test"}'
curl -X POST http://localhost:8080/api/todos -H "Content-Type: application/json" -d '{"title":"Test Todo 2","description":"Second test"}'

# Get all todos
curl http://localhost:8080/api/todos

# Get specific todo
curl http://localhost:8080/api/todos/1

# Update todo
curl -X PUT http://localhost:8080/api/todos/1 -H "Content-Type: application/json" -d '{"completed":true}'

# Delete todo
curl -X DELETE http://localhost:8080/api/todos/2
```

Stop the application (Ctrl+C) after recording.

### Step 2: Run Tests with Keploy Mocks

Now run your tests using the recorded mocks:

```bash
# Run tests with Keploy mocking
keploy test -c "go test ./... -v"

# Or with specific test path:
keploy test --path ./keploy --command "go test -v"
```

The tests will use the recorded database responses instead of connecting to a real database.

### Using Keploy MCP Tools

If you're using Keploy through Model Context Protocol (MCP):

1. **List available mocks:**
   Use `keploy_list_mocks` to see all recorded mock sets

2. **Record new mocks:**
   Use `keploy_mock_record` with command: `go run main.go`

3. **Run tests with mocks:**
   Use `keploy_mock_test` with command: `go test ./... -v`

### Benefits of Keploy Mocking

- **No Database Required**: Tests run without needing MySQL
- **Fast Tests**: No network or database overhead
- **Consistent Results**: Same responses every time
- **Easy CI/CD**: No database setup in CI pipelines
- **Realistic Data**: Uses actual database responses

## Dependencies

- [github.com/gorilla/mux](https://github.com/gorilla/mux) - HTTP router
- [github.com/go-sql-driver/mysql](https://github.com/go-sql-driver/mysql) - MySQL driver
- [github.com/DATA-DOG/go-sqlmock](https://github.com/DATA-DOG/go-sqlmock) - SQL mocking for unit tests
- [github.com/stretchr/testify](https://github.com/stretchr/testify) - Testing toolkit

## Architecture

The application follows a layered architecture:

1. **Handlers Layer** (`handlers/`): HTTP request handling and response formatting
2. **Repository Layer** (`repository/`): Database operations and data access
3. **Models Layer** (`models/`): Data structures and DTOs
4. **Database Layer** (`database/`): Database connection and initialization
5. **Routes Layer** (`routes/`): API route definitions

This separation allows for:
- Easy testing with mocks
- Clean dependency injection
- Better maintainability
- Clear separation of concerns

## Troubleshooting

**Database Connection Issues:**
- Verify MySQL is running
- Check credentials in environment variables
- Ensure databases exist

**Keploy Recording Issues:**
- Make sure Keploy CLI is installed
- Check that the application starts correctly
- Verify network calls are being made

**Test Failures:**
- Ensure test database exists and is clean
- Check that all dependencies are installed
- Verify database credentials are correct

## License

MIT
