# Mailing-List

* manage eMail-subscribers with mailing-list--microService
* create, read, upDate, & delete subscribers' eMails, confirmed-at--times, & opt-out--stati
* use JSON or gRPC API

## Environment-Variables

* `MAILING_LIST__DATABASE`: path to SQLite-dataBase (default: `dataBase.sqlite`)
* `MAILING_LIST__JSON_SERVER`: address for JSON-server (default: `localHost:3001`)
* `MAILING_LIST__GRPC_SERVER`: address for gRPC-server (default: `localHost:3002`)

## Command-Line--InterFace

### JSON

* `go run server/server.go`
```
curl -X POST http://localhost:3001/email/create -H "Content-Type: application/json" -d '{"value": "test@test.test"}' &&
curl -X PUT http://localhost:3001/email/update -H "Content-Type: application/json" -d '{"value": "test@test.test", "confirmedAt": 77777}' &&
curl -X DELETE http://localhost:3001/email/delete -H "Content-Type: application/json" -d '{"value": "test@test.test"}' &&
curl -X GET http://localhost:3001/email/get -H "Content-Type: application/json" -d '{"value": "test@test.test"}' &&
curl -X GET http://localhost:3001/email/get_batch -H "Content-Type: application/json" -d '{"count": 3, "page": 1}' &&
curl -X GET http://localhost:3001/email/get_batch -H "Content-Type: application/json" -d '{"count": 3, "page": 2}' &&
curl -X GET http://localhost:3001/email/get_batch -H "Content-Type: application/json" -d '{"count": 3, "page": 3}'
```

### GRPC (gRPC)

* `go run server/server.go`
* `go run client/client.go`

## Stack

![Go](https://img.shields.io/badge/-Go-79D4FD?style=flat-square&logo=go&logoColor=black)
![SQLite](https://img.shields.io/badge/-SQLite-044A64?style=flat-square&logo=sqlite)