.DEFAULT_GOAL := build

fmt:
	go fmt ./...
.PHONY:fmt

lint: fmt
	golint ./...
.PHONY:lint

vet: fmt
	go vet ./...
.PHONY:vet

run:
	go run ./cmd/donationService/main.go
test:
	go test ./tests/...
build:
	go build ./cmd/donationService