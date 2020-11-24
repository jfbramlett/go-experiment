clean:
	rm -rf bin
	rm -rf vendor

vendor:
	go mod vendor

build: clean vendor
	go build -o bin/server ./cmd/...

lint:
	docker run --rm -t --entrypoint=linter -v `pwd`:$(GO_PROJECT_PATH) -w $(GO_PROJECT_PATH) golang:dev-latest

test:
	go test -cover ./pkg/...