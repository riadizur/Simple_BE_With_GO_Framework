# Get all library
go mod tidy

# Set environment variables for Windows 64-bit
export GOOS=windows
export GOARCH=amd64

# Build the Windows executable
go build -o kubota-gasoline-api.exe ./cmd