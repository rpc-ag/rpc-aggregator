.PHONY: build
build:
	go build -o bin/rpc-proxy cmd/rpc-proxy/main.go

.PHONY: run
run: build
	./bin/rpc-proxy --config ./config.yaml

.PHONY: runp
runp: build
	./bin/rpc-proxy --config ./_private/solana.yaml

.PHONY: test
test:
	go test -v ./...

.PHONY: check
check:
	go mod verify
	go build -v ./...
	go vet ./...
	staticcheck ./...
	golint ./...

.PHONY: image
image: build
	docker build . -t rpc-proxy

.PHONY: run-docker
run-docker: image
	docker run -p 8080:8080 -it rpc-proxy

.PHONY: clean
clean:
	rm -rf bin/