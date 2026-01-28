# Keploy Integration Guide

This guide explains how to use Keploy to record and replay database calls for testing the Todo application.

## What is Keploy?

Keploy is a tool that records all outgoing calls (database queries, API calls, etc.) made by your application and replays them during tests. This eliminates the need for a real database in your test environment.

## Prerequisites

1. Keploy CLI installed
2. MySQL database running (only needed for recording)
3. Application configured with database credentials

## Step-by-Step Guide

### 1. Recording Mocks

First, we need to record the application's database interactions. Make sure your MySQL database is running and accessible.

#### Option A: Using Keploy MCP Tools (Recommended)

If you're using VS Code with Keploy MCP integration:

1. **Check available mocks:**
   ```
   Use: keploy_list_mocks
   ```

2. **Record mocks:**
   ```
   Use: keploy_mock_record
   Command: go run main.go
   ```

#### Option B: Using Keploy CLI

```powershell
# Start recording
keploy record -c "go run main.go"
```

### 2. Generate Test Data

While the application is running in record mode, make API calls to generate traffic:

```powershell
# Open a new terminal and run these commands

# Create todos
curl -X POST http://localhost:8080/api/todos -H "Content-Type: application/json" -d '{\"title\":\"Buy groceries\",\"description\":\"Milk, eggs, bread\"}'

curl -X POST http://localhost:8080/api/todos -H "Content-Type: application/json" -d '{\"title\":\"Write code\",\"description\":\"Finish the project\"}'

# Get all todos
curl http://localhost:8080/api/todos

# Get specific todo
curl http://localhost:8080/api/todos/1

# Update todo
curl -X PUT http://localhost:8080/api/todos/1 -H "Content-Type: application/json" -d '{\"completed\":true}'

# Delete todo
curl -X DELETE http://localhost:8080/api/todos/2
```

### 3. Stop Recording

Press `Ctrl+C` in the terminal where Keploy is running. The mocks will be saved in the `./keploy` directory.

### 4. Running Tests with Mocks

Now you can run your tests using the recorded mocks. **No database required!**

#### Option A: Using Keploy MCP Tools

```
Use: keploy_list_mocks  # See available mock sets
Use: keploy_mock_test
Command: go test ./... -v
```

#### Option B: Using Keploy CLI

```powershell
keploy test -c "go test ./... -v"
```

## Directory Structure After Recording

```
keploy/
├── test-run-0/           # First recorded session
│   ├── mocks.yaml        # Recorded database responses
│   └── config.yaml       # Recording configuration
├── test-run-1/           # Second recorded session
│   └── ...
└── ...
```

## Benefits

1. **No Database Setup**: Tests run without MySQL
2. **Fast Execution**: No network or database overhead
3. **Consistent Results**: Same responses every time
4. **CI/CD Friendly**: No infrastructure requirements
5. **Realistic Data**: Uses actual database responses

## Troubleshooting

### Issue: "No mocks found"
**Solution**: You need to record mocks first using `keploy_mock_record` or `keploy record`

### Issue: "Connection refused during recording"
**Solution**: Ensure MySQL is running and credentials in .env are correct

### Issue: "Mock not found during test"
**Solution**: The test made a call that wasn't recorded. Re-record with more comprehensive test scenarios

### Issue: "Tests still connecting to database"
**Solution**: Ensure you're running tests with `keploy test` command, not plain `go test`

## Example Workflow

### Initial Setup
```powershell
# 1. Setup environment
cp .env.example .env
# Edit .env with your MySQL credentials

# 2. Install dependencies
go mod download

# 3. Start MySQL
# (Use your preferred method)
```

### Recording Session
```powershell
# 4. Record mocks
keploy record -c "go run main.go"

# 5. In another terminal, generate traffic
curl -X POST http://localhost:8080/api/todos -H "Content-Type: application/json" -d '{\"title\":\"Test\",\"description\":\"Test desc\"}'
# ... more API calls ...

# 6. Stop recording (Ctrl+C)
```

### Testing Session
```powershell
# 7. Run tests with mocks (no database needed!)
keploy test -c "go test ./... -v"

# Or run specific tests
keploy test -c "go test ./handlers/... -v"
```

## Advanced Usage

### Recording Multiple Scenarios

You can create multiple mock sets for different test scenarios:

```powershell
# Record scenario 1: Normal operations
keploy record -c "go run main.go"
# ... make API calls for normal flow ...

# Record scenario 2: Error cases
keploy record -c "go run main.go"
# ... make API calls that trigger errors ...
```

### Using Specific Mock Sets

```powershell
# List available mocks
keploy test --list

# Use specific mock set
keploy test --mock-name test-run-0 -c "go test ./... -v"
```

## Integration with CI/CD

Add to your CI pipeline (e.g., GitHub Actions):

```yaml
- name: Run tests with Keploy
  run: keploy test -c "go test ./... -v"
```

No database setup required in CI!

## Comparison with Traditional Mocking

| Aspect | Traditional Mocks | Keploy |
|--------|------------------|--------|
| Setup | Manual mock setup | Automatic recording |
| Maintenance | Update mocks with code | Re-record when needed |
| Realism | Mock data | Real database responses |
| Coverage | Manual scenarios | All executed calls |
| Database | Not needed | Only for recording |

## Next Steps

1. Record comprehensive scenarios covering all API endpoints
2. Set up CI/CD to use Keploy tests
3. Update mocks when database schema changes
4. Use different mock sets for different test suites

## Additional Resources

- [Keploy Documentation](https://docs.keploy.io/)
- [Keploy GitHub](https://github.com/keploy/keploy)
- Project README.md for general setup instructions
