# Build stage
FROM golang:alpine AS builder

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum files (if exists)
COPY go.mod ./
COPY go.su[m] ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY main.go ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o frontend-server main.go

# Final stage
FROM alpine:latest

# Install necessary runtime dependencies
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/frontend-server ./

# Copy static files
COPY *.html *.css *.js ./

# Expose port
EXPOSE 3000

# Command to run the executable
CMD ["./frontend-server"]
