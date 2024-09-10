# Use the official Golang image as the base image
FROM golang:1.20-alpine

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files first to cache dependencies
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the rest of the application code
COPY . .

# Ensure SQLite3 is installed (Alpine package)
RUN apk add --no-cache sqlite-libs

# Build the Go app
RUN go build -o froum ./main.go

# Expose port (assuming the app runs on port 8080)
EXPOSE 8080

# Command to run the executable
CMD ["./froum"]
