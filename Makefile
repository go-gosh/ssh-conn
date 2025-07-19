BINARY=ssh-conn

.PHONY: build
build:
	go build -o $(BINARY) .

.PHONY: test
test:
	go test ./cmd/... -v 