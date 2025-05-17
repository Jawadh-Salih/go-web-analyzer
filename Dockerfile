# Start from the official Golang base image
FROM golang:1.24-alpine as server

# Set environment variables
ENV GO111MODULE=on

# Set working directory
WORKDIR /app

# Copy go mod and source
COPY go.mod go.sum ./
COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY web/ ./web/

RUN ls -R /app

# Build the binary
RUN go mod tidy 
RUN go build -o webanalyzer ./cmd/main.go

FROM scratch

COPY --from=server app/webanalyzer /webanalyzer
# Expose port
EXPOSE 8080

# Run the binary
CMD ["./webanalyzer"]
