# Use an official Golang runtime as a parent image
FROM golang:1.22

# Install SQLite3
RUN apt-get update && apt-get install -y sqlite3 libsqlite3-dev

# Set the working directory in the container
WORKDIR /app

# Copy the entire project
COPY . .

# Download all dependencies
RUN go mod download

# Build the Go app
RUN go build -o main ./cmd/server

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]