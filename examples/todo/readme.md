# Todo Reference App

This is a Todo List reference application, a basic example of how to use TopSpin library to create a simple Todo list service.

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
[INF] 2023/05/04 21:10:02.378042 NATS client connecting to nats://localhost:4222
[INF] 2023/05/04 21:10:02.378205 Server rest-server initializing at port :8081
[INF] 2023/05/04 21:10:02.380262 NATS subscribed through: [::1]:4222
[INF] 2023/05/04 21:10:02.380424 Listening on 'command' subject
[INF] 2023/05/04 21:10:07.036192 NATS publishing through: [::1]:4222
```

## Call
```shell
curl --location --request POST 'http://localhost:8081/api/v1/cmd/create-list' \
--header 'Content-Type: application/json' \
--data-raw '{
  "userID": "e014aa9d-0e21-42a0-953c-46fa3704826a",
  "name": "Todo",
  "description": "Buy apples"
}'
```

## Acknowledging reception
```shell
[INF] 2023/05/04 21:10:07.036731 Received a command event with ID: 398aa599-9459-4135-b555-06befa2b8b2e
[INF] 2023/05/04 21:10:07.17.483514 Processing 'create-list' command with this data: {UserID:e014aa9d-0e21-42a0-953c-46fa3704826a Name:Todo Description:Buy apples}
```
