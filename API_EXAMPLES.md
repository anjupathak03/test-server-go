# API Testing Examples

This file contains example curl commands for testing the Todo API.

## Setup

Start the server:
```bash
go run main.go
```

Or if recording with Keploy:
```bash
keploy record -c "go run main.go"
```

## Create Todos

### Create Todo 1
```bash
curl -X POST http://localhost:8080/api/todos \
  -H "Content-Type: application/json" \
  -d '{"title":"Buy groceries","description":"Milk, eggs, bread, cheese"}'
```

### Create Todo 2
```bash
curl -X POST http://localhost:8080/api/todos \
  -H "Content-Type: application/json" \
  -d '{"title":"Finish project","description":"Complete the Go todo app"}'
```

### Create Todo 3
```bash
curl -X POST http://localhost:8080/api/todos \
  -H "Content-Type: application/json" \
  -d '{"title":"Exercise","description":"Go for a 30 minute run"}'
```

### Create Todo 4
```bash
curl -X POST http://localhost:8080/api/todos \
  -H "Content-Type: application/json" \
  -d '{"title":"Read documentation","description":"Read Keploy docs"}'
```

## Read Todos

### Get All Todos
```bash
curl http://localhost:8080/api/todos
```

### Get Todo by ID (replace {id} with actual ID)
```bash
curl http://localhost:8080/api/todos/1
curl http://localhost:8080/api/todos/2
curl http://localhost:8080/api/todos/3
```

## Update Todos

### Mark Todo as Completed
```bash
curl -X PUT http://localhost:8080/api/todos/1 \
  -H "Content-Type: application/json" \
  -d '{"completed":true}'
```

### Update Todo Title
```bash
curl -X PUT http://localhost:8080/api/todos/2 \
  -H "Content-Type: application/json" \
  -d '{"title":"Finish project by Friday"}'
```

### Update Multiple Fields
```bash
curl -X PUT http://localhost:8080/api/todos/3 \
  -H "Content-Type: application/json" \
  -d '{"title":"Morning Exercise","description":"Run 5km","completed":true}'
```

## Delete Todos

### Delete Todo by ID
```bash
curl -X DELETE http://localhost:8080/api/todos/4
```

### Delete Another Todo
```bash
curl -X DELETE http://localhost:8080/api/todos/2
```

## Complete Workflow Example

This sequence demonstrates a complete CRUD workflow:

```bash
# 1. Create a todo
curl -X POST http://localhost:8080/api/todos \
  -H "Content-Type: application/json" \
  -d '{"title":"Test Workflow","description":"Testing CRUD operations"}'

# Note the ID from the response, let's say it's 5

# 2. Get the created todo
curl http://localhost:8080/api/todos/5

# 3. Update the todo
curl -X PUT http://localhost:8080/api/todos/5 \
  -H "Content-Type: application/json" \
  -d '{"completed":true}'

# 4. Get all todos to verify
curl http://localhost:8080/api/todos

# 5. Delete the todo
curl -X DELETE http://localhost:8080/api/todos/5

# 6. Verify deletion
curl http://localhost:8080/api/todos/5
# Should return 404 Not Found
```

## PowerShell Examples

If you're using Windows PowerShell, use these commands instead:

### Create Todo
```powershell
$body = @{
    title = "Buy groceries"
    description = "Milk, eggs, bread"
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:8080/api/todos" -Method Post -Body $body -ContentType "application/json"
```

### Get All Todos
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/api/todos" -Method Get
```

### Get Todo by ID
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/api/todos/1" -Method Get
```

### Update Todo
```powershell
$body = @{
    completed = $true
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:8080/api/todos/1" -Method Put -Body $body -ContentType "application/json"
```

### Delete Todo
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/api/todos/1" -Method Delete
```

## Testing Error Cases

### Invalid Request - Missing Required Field
```bash
curl -X POST http://localhost:8080/api/todos \
  -H "Content-Type: application/json" \
  -d '{"description":"Missing title"}'
# Should return 400 Bad Request
```

### Invalid ID Format
```bash
curl http://localhost:8080/api/todos/invalid
# Should return 400 Bad Request
```

### Non-existent Todo
```bash
curl http://localhost:8080/api/todos/9999
# Should return 404 Not Found
```

### Invalid JSON
```bash
curl -X POST http://localhost:8080/api/todos \
  -H "Content-Type: application/json" \
  -d '{invalid json}'
# Should return 400 Bad Request
```

## Bulk Operations Script

Save this as a bash script to quickly populate test data:

```bash
#!/bin/bash
# populate_todos.sh

BASE_URL="http://localhost:8080/api/todos"

echo "Creating test todos..."

curl -X POST $BASE_URL -H "Content-Type: application/json" -d '{"title":"Learn Go","description":"Study Go programming language"}'
curl -X POST $BASE_URL -H "Content-Type: application/json" -d '{"title":"Build API","description":"Create REST API with Go"}'
curl -X POST $BASE_URL -H "Content-Type: application/json" -d '{"title":"Write tests","description":"Add unit and integration tests"}'
curl -X POST $BASE_URL -H "Content-Type: application/json" -d '{"title":"Deploy app","description":"Deploy to production"}'
curl -X POST $BASE_URL -H "Content-Type: application/json" -d '{"title":"Monitor","description":"Set up monitoring and alerts"}'

echo ""
echo "Getting all todos..."
curl $BASE_URL

echo ""
echo "Done!"
```

Make it executable and run:
```bash
chmod +x populate_todos.sh
./populate_todos.sh
```

## PowerShell Bulk Operations Script

Save this as populate_todos.ps1:

```powershell
# populate_todos.ps1

$baseUrl = "http://localhost:8080/api/todos"

$todos = @(
    @{ title = "Learn Go"; description = "Study Go programming language" },
    @{ title = "Build API"; description = "Create REST API with Go" },
    @{ title = "Write tests"; description = "Add unit and integration tests" },
    @{ title = "Deploy app"; description = "Deploy to production" },
    @{ title = "Monitor"; description = "Set up monitoring and alerts" }
)

Write-Host "Creating test todos..." -ForegroundColor Green

foreach ($todo in $todos) {
    $body = $todo | ConvertTo-Json
    Invoke-RestMethod -Uri $baseUrl -Method Post -Body $body -ContentType "application/json"
    Write-Host "Created: $($todo.title)" -ForegroundColor Cyan
}

Write-Host "`nGetting all todos..." -ForegroundColor Green
$allTodos = Invoke-RestMethod -Uri $baseUrl -Method Get
$allTodos | ConvertTo-Json

Write-Host "`nDone!" -ForegroundColor Green
```

Run it:
```powershell
.\populate_todos.ps1
```

## Tips for Keploy Recording

When recording with Keploy, execute a comprehensive set of operations:

1. Create several todos
2. Get all todos
3. Get individual todos by ID
4. Update todos (both completed status and content)
5. Delete some todos
6. Try error cases (invalid IDs, missing fields)

This ensures your mocks cover all scenarios!
