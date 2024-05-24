# Get all library
go mod tidy

# Build the executable
go build -o kubota-gasoline-api ./cmd

# Run the executable
./kubota-gasoline-api