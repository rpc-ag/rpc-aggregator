FROM golang:latest as builder

# Set the working directory
WORKDIR /app

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./

# Download the dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the binary
RUN go build -o rpc-proxy cmd/rpc-proxy/main.go

# Use a lightweight image as the final image
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy the binary from the builder image
COPY --from=builder /app/rpc-proxy .

# Copy the config.yaml file
COPY config.yaml .

# Expose the port the server will listen on
EXPOSE 8080

# Start the rpc-proxy server
CMD ["./rpc-proxy"]
