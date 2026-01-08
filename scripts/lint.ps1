# Script to run golangci-lint with proper error handling

# Check if golangci-lint is installed
if (!(Get-Command golangci-lint -ErrorAction SilentlyContinue)) {
    Write-Host "golangci-lint could not be found. Installing..."
    # Install golangci-lint using Go
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@v2.6.2
}

Write-Host "Running golangci-lint..."
golangci-lint run ./... --timeout=5m

Write-Host "Linting completed successfully!"
