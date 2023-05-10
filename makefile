format:
	gofmt -s -w .

build:
	go build ./...

test:
	go clean -testcache
	go test ./... -v

run:
	go run examples/todo/cmd/todo.go --config-file=examples/todo/configs/config.yml

create-list-command:
	curl --location 'http://localhost:8081/api/v1/cmd/create-list' \
	--header 'Content-Type: application/json' \
	--data '{\
  	"userUUID": "e014aa9d-0e21-42a0-953c-46fa3704826a",\
  	"name": "Todo",\
  	"description": "Buy apples"\
	}'

.PHONY: openapihttp
openapihttp:
	oapi-codegen -generate types -o examples/todo/internal/ports/openapi/todotypes.go -package openapi examples/todo/api/openapi/todo.yml
	oapi-codegen -generate chi-server -o examples/todo/internal/ports/openapi/todoapi.go -package openapi examples/todo/api/openapi/todo.yml
	oapi-codegen -generate types -o examples/todo/internal/client/ports/openapi/todotypes.go -package openapi examples/todo/api/openapi/todo.yml
	oapi-codegen -generate client -o examples/todo/internal/client/ports/openapi/todoapi.go -package openapi examples/todo/api/openapi/todo.yml


.PHONY: genuml
genuml:
	plantuml ./docs/diagrams/sequence-draft.txt
