# Todo App PowerShell Helper Script
# Easy commands to manage the Todo application

param(
    [Parameter(Mandatory=$true)]
    [ValidateSet('help', 'deps', 'build', 'run', 'test', 'test-unit', 'keploy-record', 'keploy-test', 'keploy-list', 'clean')]
    [string]$Command
)

function Show-Help {
    Write-Host "Todo App Helper Script" -ForegroundColor Cyan
    Write-Host "======================" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "Usage: .\run.ps1 <command>" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "Available commands:" -ForegroundColor Green
    Write-Host "  help           - Show this help message"
    Write-Host "  deps           - Install Go dependencies"
    Write-Host "  build          - Build the application"
    Write-Host "  run            - Run the application"
    Write-Host "  test           - Run all tests"
    Write-Host "  test-unit      - Run unit tests only"
    Write-Host "  keploy-record  - Record mocks with Keploy"
    Write-Host "  keploy-test    - Run tests with Keploy mocks"
    Write-Host "  keploy-list    - List available mock sets"
    Write-Host "  clean          - Clean build artifacts"
    Write-Host ""
    Write-Host "Examples:" -ForegroundColor Magenta
    Write-Host "  .\run.ps1 deps"
    Write-Host "  .\run.ps1 test-unit"
    Write-Host "  .\run.ps1 keploy-record"
}

function Install-Deps {
    Write-Host "Installing dependencies..." -ForegroundColor Green
    go mod download
    go mod tidy
    Write-Host "Dependencies installed!" -ForegroundColor Green
}

function Build-App {
    Write-Host "Building application..." -ForegroundColor Green
    New-Item -ItemType Directory -Force -Path "bin" | Out-Null
    go build -o bin/todo-server.exe main.go
    Write-Host "Build complete! Binary: bin/todo-server.exe" -ForegroundColor Green
}

function Run-App {
    Write-Host "Starting Todo Server..." -ForegroundColor Green
    Write-Host "Press Ctrl+C to stop" -ForegroundColor Yellow
    go run main.go
}

function Run-Tests {
    Write-Host "Running all tests..." -ForegroundColor Green
    go test ./... -v
}

function Run-UnitTests {
    Write-Host "Running unit tests..." -ForegroundColor Green
    go test ./handlers/... ./repository/... -v
}

function Record-Keploy {
    Write-Host "Starting Keploy recording..." -ForegroundColor Green
    Write-Host "The application will start. Make API calls to record mocks." -ForegroundColor Yellow
    Write-Host "Press Ctrl+C when done recording." -ForegroundColor Yellow
    Write-Host ""
    keploy record -c "go run main.go"
}

function Test-Keploy {
    Write-Host "Running tests with Keploy mocks..." -ForegroundColor Green
    keploy test -c "go test ./... -v"
}

function List-Keploy {
    Write-Host "Available Keploy mock sets:" -ForegroundColor Green
    if (Test-Path "keploy") {
        Get-ChildItem -Path "keploy" -Directory | ForEach-Object {
            Write-Host "  - $($_.Name)" -ForegroundColor Cyan
        }
    } else {
        Write-Host "No mocks found. Run 'keploy-record' first." -ForegroundColor Yellow
    }
}

function Clean-All {
    Write-Host "Cleaning build artifacts..." -ForegroundColor Green
    if (Test-Path "bin") {
        Remove-Item -Recurse -Force "bin"
    }
    go clean -testcache
    Write-Host "Clean complete!" -ForegroundColor Green
}

# Execute command
switch ($Command) {
    'help' { Show-Help }
    'deps' { Install-Deps }
    'build' { Build-App }
    'run' { Run-App }
    'test' { Run-Tests }
    'test-unit' { Run-UnitTests }
    'keploy-record' { Record-Keploy }
    'keploy-test' { Test-Keploy }
    'keploy-list' { List-Keploy }
    'clean' { Clean-All }
}
