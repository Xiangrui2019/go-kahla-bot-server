default: deps
	go build

deps:
	go mod tidy
	go mod vendor

clean:
	rm -rf vendor
	rm -rf go.sum

.PHONY: go test