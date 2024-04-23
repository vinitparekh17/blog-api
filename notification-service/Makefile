.DEFAULT_GOAL := build

.PHONY: fmt vet build 
fmt:
	go fmt ./...
vet:
	go vet ./...
build:
	go build -o main


.PHONY: clean
clean:
	go clean

.PHONY: test
test:
	go test -v ./...

.PHONY: run
run: build
	./main
