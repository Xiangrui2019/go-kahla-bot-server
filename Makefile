default: deps
	go install
	cp ./config.toml "${GOPATH}/bin/config.toml"

deps:
	go mod tidy
	go mod vendor

run: default
	cd ${GOPATH}/bin;./go-kahla-bot-server

docker: deps
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
	chmod +x ./main
	sudo docker build -t go-kahla-bot-server .

clean:
	rm -rf vendor/
	rm -rf go.sum

.PHONY: go test