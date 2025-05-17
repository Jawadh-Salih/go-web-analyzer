# Start from the official Golang base image
FROM golang:1.24-alpine AS server

# Set environment variables
ENV GO111MODULE=on

# Set working directory
WORKDIR /app

# Copy go mod and source
COPY go.mod go.sum ./

# Build the binary
RUN go mod tidy 

COPY . .

RUN go build -o webanalyzer ./cmd/main.go

FROM scratch

COPY --from=server app/webanalyzer /webanalyzer
# Load the HTML folder into the image
COPY --from=server app/web /web  
# Load the CA certificates into the image
COPY --from=server /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/


# Expose port
EXPOSE 8080

# Run the binary
CMD ["./webanalyzer"]
