.PHONY: build
build:
	go build -o bin/rpc-aggregator cmd/rpc-aggregator/main.go

.PHONY: run
run: build
	./bin/rpc-aggregator --config ./config.yaml

.PHONY: runp
runp: build
	./bin/rpc-aggregator --config ./_private/solana.yaml

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
	docker build . -t rpc-aggregator

.PHONY: run-docker
run-docker: image
	docker run -p 8080:8080 -it rpc-aggregator

.PHONY: clean
clean:
	rm -rf bin/