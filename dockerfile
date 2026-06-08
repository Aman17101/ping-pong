# Step 1: Build the Go binary
FROM golang:1.21-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o pingpong .

# Step 2: Create a small runtime image
FROM alpine:latest

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/pingpong .

# IMPORTANT: Copy the static folder so the server can find the HTML/JS
COPY --from=builder /app/static ./static

# Expose the port the app runs on
EXPOSE 8080

# Run the app
CMD ["./pingpong"]