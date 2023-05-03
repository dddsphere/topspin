format:
	gofmt -s -w .

build:
	go build ./...

test:
	go clean -testcache
	go test ./... -v

run:
	go run cmd/todo.go

.PHONY: openapihttp
openapihttp:
	oapi-codegen -generate types -o examples/todo/internal/app/ports/openapi/todotypes.go -package openapi api/openapi/todo.yml
	oapi-codegen -generate chi-server -o examples/todo/internal/app/ports/openapi/todoapi.go -package openapi api/openapi/todo.yml
	oapi-codegen -generate types -o examples/todo/internal/client/ports/openapi/todotypes.go -package openapi api/openapi/todo.yml
	oapi-codegen -generate client -o examples/todo/internal/client/ports/openapi/todoapi.go -package openapi api/openapi/todo.yml

.PHONY: genuml
genuml:
	plantuml ./docs/diagrams/sequence-draft.txt
