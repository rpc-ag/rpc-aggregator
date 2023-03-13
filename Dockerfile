FROM golang:latest as builder

# Set the working directory
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /app/rpc-aggregator ./cmd/rpc-aggregator/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/rpc-aggregator /usr/local/bin/rpc-aggregator
COPY ./_private/solana.yaml /app/config.yaml
COPY ./_private/auth.yaml /app/auth.yaml
EXPOSE 8080
#CMD ["/bin/sh"]
CMD ["/usr/local/bin/rpc-aggregator","--config","/app/config.yaml","--auth","/app/auth.yaml"]
