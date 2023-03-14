# Q&A Service

A service that exposes a REST API which allows to create, update, delete and retrieve answers as key-value pairs. 
The API also supports a `get history` operation for a given key.

Design Patterns used:
- [Event Sourcing](https://www.eventstore.com/event-sourcing)
- [CQRS](https://www.eventstore.com/cqrs-pattern)

Databases used:
- [EventStoreDB](https://www.eventstore.com/) as event store
- [MongoDB](https://www.mongodb.com/) as read repository


## Run it with in-memory db

```
$ godotenv go run . -memdb
```


## Run it with databases

```
$ docker-compose up -d
```

Event Store Dashboard 
http://localhost:2113/


## Usage examples

```
- create answer
$ curl http://localhost:8080/answers \
  -X POST -d '{"key":"name","value":"john"}' | jq .
{
  "ok": "AnswerCreatedEvent"
}

- get answer
$ curl http://localhost:8080/answers/name | jq .
{
  "key": "name",
  "value": "john"
}

# error on conflict
> curl http://localhost:8080/answers \
  -X POST -d '{"key":"name","value":"john"}' | jq .
{
  "error": "Key exists"
}

# get history for given key
$ curl http://localhost:8080/answers/name/history | jq .
[
  {
    "type": "AnswerCreatedEvent",
    "data": {
      "key": "name",
      "value": "john"
    }
  }
]

# update answer
$ curl http://localhost:8080/answers/name \
  -X PATCH -d '{"key":"name","value":"jack"}'  | jq .
{
  "ok": "AnswerUpdatedEvent"
}

# fetch updated
$ curl http://localhost:8080/answers/name | jq .
{
  "key": "name",
  "value": "jack"
}
```


## Development

```
$ docker-compose up -d mongo eventstore.db
$ godotenv go run .
```


### Test

Run all tests (requires databases to be up):
```
$ docker-compose up -d mongo eventstore.db
$ godotenv go test ./...
```

Skip tests flagged as "short":
```
$ godotenv go test -short ./...
```
