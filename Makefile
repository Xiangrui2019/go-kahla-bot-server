default: deps
	go build

deps:
	go mod tidy
	go mod vendor

.PHONY: go test