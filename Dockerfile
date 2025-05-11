# Start from the official Golang base image
FROM golang:1.21-alpine

# Set environment variables
ENV GO111MODULE=on

# Set working directory
WORKDIR /app

# Copy go mod and source
COPY go.mod ./
COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY web/ ./web/

# Build the binary
RUN go build -o webanalyzer ./cmd/main.go

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./webanalyzer"]
