# Todo Reference App

This is a Todo List reference application, a basic example of how to use TopSpin library to create a simploptodo list microservice. This library is designed to help you build Go microservices that follow Domain Driven Design (DDD) practices, while remaining lightweight and efficient. CQRS and Event Sourcing are first-class citizens in microservice architecture, but we are not restricted to them. Throughout the development process, we will showcase examples of those DDD patterns but Additionally, we understand that there may be use cases where a plain RESTful approach is more suitable. 

## Prerequisites
Install NATS server
```shell
$ GO111MODULE=on go get github.com/nats-io/nats-server/v2
(...)
$ nats-server -m 8222
```

Alternatively
```shell
$ docker pull nats:latest
$ docker run -p 4222:4222 -ti nats:latest
````

## Run
```shell
$ make run
go run cmd/todo.go
[INF] 2023/05/03 20:53:34.805728 NATS client connecting to nats://localhost:4222
[INF] 2023/05/03 20:53:34.805858 Server rest-server initializing at port :8081
[INF] 2023/05/03 20:53:34.807780 NATS subscribed through: [::1]:4222
[INF] 2023/05/03 20:53:34.807921 Listening on 'commands' subject
```

## Call
```shell
curl --location --request POST 'http://localhost:8081/api/v1/cmd/create-list' \
--header 'Content-Type: application/json' \
--data-raw '{
  "userUUID": "e014aa9d-0e21-42a0-953c-46fa3704826a",
  "name": "Todo",
  "description": "Buy apples"
}'
```
