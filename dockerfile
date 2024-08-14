# # Use an official Golang runtime as a parent image
# FROM golang:1.22

# # Install SQLite3
# RUN apt-get update && apt-get install -y sqlite3 libsqlite3-dev

# # Set the working directory in the container
# WORKDIR /app

# # Copy the entire project
# COPY . .

# # Download all dependencies
# RUN go mod download

# # Build the Go app
# RUN go build -o main ./cmd/server

# # Set the DB_PATH environment variable
# ENV DB_PATH=/app/data/forum.db

# # Expose port 8080 to the outside world
# EXPOSE 8080

# # Command to run the executable
# CMD ["./main"]

# FOR LIVE SERVER: 
FROM golang:1.22

WORKDIR /app

# Install Air
RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Ensure the tmp directory exists and has correct permissions
RUN mkdir -p /.cache && chmod 777 /.cache
RUN go build -o /app/tmp/main ./cmd/server
# Set the correct permissions for the entire /app directory
RUN chmod -R 755 /app

# Command to run Air
CMD ["air", "-c", ".air.toml"]