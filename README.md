# Q&A Service
A service that exposes a REST API which allows to create, update, delete and retrieve answers as key-value pairs.

Design Patterns used:
- [Event Sourcing](https://www.eventstore.com/event-sourcing)
- [CQRS](https://www.eventstore.com/cqrs-pattern)

Databases used:
- [EventStoreDB](https://www.eventstore.com/) as event store
- [MongoDB](https://www.mongodb.com/) as read repository

## Run it
```
$ docker-compose up -d
```

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
$ godotenv go run .
```

### Test
```
$ go clean -testcache && godotenv go test ./...
```

## TODO
- Handle errors