# Get all library
go mod tidy

# Set environment variables for Windows 64-bit
go build -o kubota-gasoline-api ./cmd

# Build the Windows executable
./kubota-gasoline-api