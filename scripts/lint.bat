@echo off
setlocal

REM Check if golangci-lint is installed
where golangci-lint >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo golangci-lint not found. Installing...
    powershell -Command "Invoke-WebRequest -Uri 'https://github.com/golangci/golangci-lint/releases/download/v1.60.3/golangci-lint-1.60.3-windows-amd64.zip' -OutFile 'golangci-lint.zip'"
    powershell -Command "Expand-Archive -Path 'golangci-lint.zip' -DestinationPath '.'"
    copy "golangci-lint-1.60.3-windows-amd64\golangci-lint.exe" "%GOPATH%\bin\"
    del "golangci-lint.zip"
    rmdir /s /q "golangci-lint-1.60.3-windows-amd64"
)

echo Running golangci-lint...
golangci-lint run ./... --timeout=5m

echo Linting completed successfully!