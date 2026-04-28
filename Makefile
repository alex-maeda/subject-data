.PHONY: build test run lint fmt swagger

build: swagger
	CGO_ENABLED=0 go build -o bin/run_service ./cmd/run_service

test:
	go test -race ./...

run:
	go run ./cmd/run_service

lint:
	golangci-lint run ./...

fmt:
	gofmt -w .

swagger:
	swag init -d ./cmd/run_service,./internal/api -g doc.go --output internal/swaggerdocs --packageName swaggerdocs --parseDependency --parseInternal
