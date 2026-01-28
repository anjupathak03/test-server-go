# Test API Script for Keploy Recording
# This script exercises all API endpoints to generate comprehensive mocks

$baseUrl = "http://localhost:8080/api/todos"

Write-Host "Testing Todo API..." -ForegroundColor Green
Start-Sleep -Seconds 2

# Create Todo 1
Write-Host "`nCreating Todo 1..." -ForegroundColor Cyan
$todo1 = @{
    title = "Buy groceries"
    description = "Milk, eggs, bread, cheese"
} | ConvertTo-Json

try {
    $response1 = Invoke-RestMethod -Uri $baseUrl -Method Post -Body $todo1 -ContentType "application/json"
    Write-Host "Created Todo ID: $($response1.id)" -ForegroundColor Green
} catch {
    Write-Host "Error creating todo: $_" -ForegroundColor Red
}

Start-Sleep -Seconds 1

# Create Todo 2
Write-Host "`nCreating Todo 2..." -ForegroundColor Cyan
$todo2 = @{
    title = "Finish project"
    description = "Complete the Go todo app"
} | ConvertTo-Json

try {
    $response2 = Invoke-RestMethod -Uri $baseUrl -Method Post -Body $todo2 -ContentType "application/json"
    Write-Host "Created Todo ID: $($response2.id)" -ForegroundColor Green
} catch {
    Write-Host "Error creating todo: $_" -ForegroundColor Red
}

Start-Sleep -Seconds 1

# Create Todo 3
Write-Host "`nCreating Todo 3..." -ForegroundColor Cyan
$todo3 = @{
    title = "Exercise"
    description = "Go for a 30 minute run"
} | ConvertTo-Json

try {
    $response3 = Invoke-RestMethod -Uri $baseUrl -Method Post -Body $todo3 -ContentType "application/json"
    Write-Host "Created Todo ID: $($response3.id)" -ForegroundColor Green
} catch {
    Write-Host "Error creating todo: $_" -ForegroundColor Red
}

Start-Sleep -Seconds 1

# Get All Todos
Write-Host "`nGetting all todos..." -ForegroundColor Cyan
try {
    $allTodos = Invoke-RestMethod -Uri $baseUrl -Method Get
    Write-Host "Retrieved $($allTodos.Count) todos" -ForegroundColor Green
} catch {
    Write-Host "Error getting todos: $_" -ForegroundColor Red
}

Start-Sleep -Seconds 1

# Get Todo by ID
Write-Host "`nGetting Todo by ID (1)..." -ForegroundColor Cyan
try {
    $todo = Invoke-RestMethod -Uri "$baseUrl/1" -Method Get
    Write-Host "Retrieved: $($todo.title)" -ForegroundColor Green
} catch {
    Write-Host "Error getting todo: $_" -ForegroundColor Red
}

Start-Sleep -Seconds 1

# Update Todo
Write-Host "`nUpdating Todo 1 (mark as completed)..." -ForegroundColor Cyan
$update = @{
    completed = $true
} | ConvertTo-Json

try {
    $updated = Invoke-RestMethod -Uri "$baseUrl/1" -Method Put -Body $update -ContentType "application/json"
    Write-Host "Updated - Completed: $($updated.completed)" -ForegroundColor Green
} catch {
    Write-Host "Error updating todo: $_" -ForegroundColor Red
}

Start-Sleep -Seconds 1

# Update Todo Title
Write-Host "`nUpdating Todo 2 (change title)..." -ForegroundColor Cyan
$update2 = @{
    title = "Finish project by Friday"
    description = "Complete the Go todo app and deploy"
} | ConvertTo-Json

try {
    $updated2 = Invoke-RestMethod -Uri "$baseUrl/2" -Method Put -Body $update2 -ContentType "application/json"
    Write-Host "Updated: $($updated2.title)" -ForegroundColor Green
} catch {
    Write-Host "Error updating todo: $_" -ForegroundColor Red
}

Start-Sleep -Seconds 1

# Delete Todo
Write-Host "`nDeleting Todo 3..." -ForegroundColor Cyan
try {
    $deleted = Invoke-RestMethod -Uri "$baseUrl/3" -Method Delete
    Write-Host "Deleted successfully" -ForegroundColor Green
} catch {
    Write-Host "Error deleting todo: $_" -ForegroundColor Red
}

Start-Sleep -Seconds 1

# Get All Todos Again
Write-Host "`nGetting all todos after operations..." -ForegroundColor Cyan
try {
    $finalTodos = Invoke-RestMethod -Uri $baseUrl -Method Get
    Write-Host "Final count: $($finalTodos.Count) todos" -ForegroundColor Green
    $finalTodos | ConvertTo-Json | Write-Host
} catch {
    Write-Host "Error getting todos: $_" -ForegroundColor Red
}

Write-Host "`nAPI Testing Complete!" -ForegroundColor Green
Write-Host "Press Ctrl+C in the server terminal to stop recording" -ForegroundColor Yellow
