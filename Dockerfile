# Stage 1: Build the Go binary using a Debian-based image with proper build tools
FROM golang:1.20-buster AS builder

# Install dependencies for CGO and SQLite3
RUN apt-get update && apt-get install -y \
    gcc \
    libc6-dev \
    sqlite3 \
    libsqlite3-dev

# Set environment variables to enable CGO
# Use default GCC for the current system architecture
ENV CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

# Set the working directory inside the container
WORKDIR /app

# Copy Go module files
COPY go.mod go.sum ./

# Download Go modules
RUN go mod download

# Copy the entire project
COPY . .

# Build the Go project targeting Linux AMD64 architecture
RUN go build -o froum main.go

# Stage 2: Create a smaller runtime-only container using Alpine
FROM alpine:3.18

# Install SQLite runtime to support the application
RUN apk add --no-cache sqlite-libs

# Set working directory inside the container
WORKDIR /app

# Copy the compiled Go binary from the builder stage
COPY --from=builder /app/froum /app/froum

# Copy static assets, templates, and other required files
COPY css/ /app/css/
COPY handlers/ /app/handlers/
COPY models/ /app/models/
COPY storage/ /app/storage/
COPY templates/ /app/templates/

# Expose the port your app will run on
EXPOSE 8080

# Run the binary
CMD ["./froum"]