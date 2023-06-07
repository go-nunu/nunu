## Documentation
* [Guide](https://github.com/go-nunu/nunu/blob/main/docs/en/guide.md)
* [Architecture](https://github.com/go-nunu/nunu/blob/main/docs/en/architecture.md)
* [Tutorial](https://github.com/go-nunu/nunu/blob/main/docs/en/tutorial.md)

[进入简体中文版](https://github.com/go-nunu/nunu/blob/main/docs/zh/architecture.md)


# Architecture Overview

Nunu adopts the classic layered architecture. At the same time, in order to better achieve modularity and decoupling, it uses the dependency injection framework `Wire`.

![Nunu Layout](https://github.com/go-nunu/nunu/blob/main/.github/assets/layout.jpg)

## Directory Structure

```
.
├── cmd
│   ├── job
│   │   ├── wire
│   │   │   ├── wire.go
│   │   │   └── wire_gen.go
│   │   └── main.go
│   ├── migration
│   │   ├── wire
│   │   │   ├── wire.go
│   │   │   └── wire_gen.go
│   │   └── main.go
│   └── server
│       ├── wire
│       │   ├── wire.go
│       │   └── wire_gen.go
│       └── main.go
├── config
│   ├── local.yml
│   └── prod.yml
├── deploy
│   ├── build
│   │   └── Dockerfile
│   └── docker-composer
│       └── docker-composer.yml
├── internal
│   ├── dao
│   │   ├── dao.go
│   │   └── user.go
│   ├── handler
│   │   ├── handler.go
│   │   └── user.go
│   ├── job
│   │   └── job.go
│   ├── middleware
│   │   ├── cors.go
│   │   ├── jwt.go
│   │   ├── log.go
│   │   └── sign.go
│   ├── migration
│   │   └── migration.go
│   ├── model
│   │   └── user.go
│   ├── provider
│   │   └── provider.go
│   ├── server
│   │   └── http.go
│   └── service
│       ├── service.go
│       └── user.go
├── pkg
│   ├── config
│   │   └── config.go
│   ├── helper
│   │   ├── md5
│   │   │   └── md5.go
│   │   ├── resp
│   │   │   └── resp.go
│   │   ├── sonyflake
│   │   │   └── sonyflake.go
│   │   └── uuid
│   │       └── uuid.go
│   ├── http
│   │   └── http.go
│   └── log
│       └── log.go
├── script
│   └── deploy.sh
├── storage
│   └── logs
├── test
│   └── server
│       ├── handler
│       │   └── user_test.go
│       └── service
│           └── user_test.go
├── web
├── LICENSE
├── README.md
├── README_zh.md
├── go.mod
└── go.sum
```

- `cmd`: the entry point of the application, containing different subcommands.
- `config`: configuration files.
- `deploy`: deployment-related files, such as Dockerfile, docker-compose.yml, etc.
- `internal`: the main code of the application, organized according to the layered architecture.
- `pkg`: public code, including configuration, logging, HTTP, etc.
- `script`: script files for deployment and other automation tasks.
- `storage`: storage files, such as log files.
- `test`: test code.
- `web`: front-end code.

## internal

- `internal/handler` (or `controller`): handles HTTP requests, calls services in the business logic layer, and returns HTTP responses.
- `internal/server` (or `router`): HTTP server, starts the HTTP service, listens on the port, and handles HTTP requests.
- `internal/service` (or `logic`): service, implements specific business logic, and calls the DAO in the data access layer.
- `internal/model` (or `entity`): data model, defines the data structure required by the business logic layer.
- `internal/dao` (or `repository`): data access object, encapsulates database operations and provides data manipulation functions such as CRUD.
- `internal/middleware`: middleware, used to handle requests and responses, such as logging, cross-domain, signature, etc.

## Dependency Injection

This project uses the dependency injection framework `Wire` to achieve modularity and decoupling. `Wire` generates dependency injection code `wire_gen.go` by pre-compiling `wire.go`, simplifying the dependency injection process.

- `cmd/job/wire/wire.go`: `Wire` configuration file, defines the dependency relationship required by the `job` subcommand.
- `cmd/migration/wire/wire.go`: `Wire` configuration file, defines the dependency relationship required by the `migration` subcommand.
- `cmd/server/wire/wire.go`: `Wire` configuration file, defines the dependency relationship required by the `server` subcommand.
- `internal/provider/wire.go`: `Wire` provider declaration.

Wire official documentation: https://github.com/google/wire/blob/main/docs/guide.md

Note: `wire_gen.go` is automatically compiled and generated, and manual modification is prohibited.

## Public Code

To achieve code reuse and unified management, this project uses public code to place some common code under the `pkg` directory.

- `pkg/config`: reads and parses configuration files.
- `pkg/helper`: some common auxiliary functions, such as MD5 encryption, UUID generation, etc.
- `pkg/http`: HTTP-related code, such as HTTP client, HTTP server, etc.
- `pkg/log`: logging-related code, such as log initialization, log writing, etc.
- `more...`: Of course, you can freely add and extend more `pkg`.