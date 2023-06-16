## Documentation
* [User Guide](https://github.com/go-nunu/nunu/blob/main/docs/en/guide.md)
* [Architecture](https://github.com/go-nunu/nunu/blob/main/docs/en/architecture.md)
* [Getting Started Tutorial](https://github.com/go-nunu/nunu/blob/main/docs/en/tutorial.md)


[切换简体中文](https://github.com/go-nunu/nunu/blob/main/docs/zh/architecture.md)

# Exploring the Nunu Architecture

Nunu adopts a classic layered architecture. Additionally, to achieve better modularity and decoupling, it utilizes the dependency injection framework `Wire`.

![Nunu Layout](https://github.com/go-nunu/nunu/blob/main/.github/assets/layout.png)

## Directory Structure
```
.
├── cmd
│   ├── job
│   │   ├── main.go
│   │   ├── wire.go
│   │   └── wire_gen.go
│   ├── migration
│   │   ├── main.go
│   │   ├── wire.go
│   │   └── wire_gen.go
│   └── server
│       ├── main.go
│       ├── wire.go
│       └── wire_gen.go
├── config
│   ├── local.yml
│   └── prod.yml
├── deploy
├── internal
│   ├── handler
│   │   ├── handler.go
│   │   └── user.go
│   ├── job
│   │   └── job.go
│   ├── middleware
│   ├── migration
│   │   └── migration.go
│   ├── model
│   │   └── user.go
│   ├── repository
│   │   ├── repository.go
│   │   └── user.go
│   ├── server
│   │   └── http.go
│   └── service
│       ├── service.go
│       └── user.go
├── mocks
│   ├── repository
│   │   └── user.go
│   └── service
│       └── user.go
├── pkg
├── scripts
├── storage
├── test
│   └── server
│       ├── handler
│       │   └── user_test.go
│       ├── repository
│       │   └── user_test.go
│       └── service
│           └── user_test.go
├── web
│   └── index.html
├── LICENSE
├── Makefile
├── README.md
├── README_zh.md
├── coverage.html
├── go.mod
└── go.sum

```


- `cmd`: Entry point of the application, containing different subcommands.
- `config`: Configuration files.
- `deploy`: Files related to deployment, such as Dockerfile and docker-compose.yml.
- `internal`: Main code of the application, organized according to the layered architecture.
- `mocks`: Mock code for testing.
- `pkg`: Common code, including configuration, logging, and HTTP.
- `scripts`: Script files for deployment and other automation tasks.
- `storage`: Storage files, such as log files.
- `test`: Test code.
- `web`: Front-end code.

## internal

- `internal/handler` (or `controller`): Handles HTTP requests, calls services in the business logic layer, and returns HTTP responses.
- `internal/server` (or `router`): HTTP server that starts the HTTP service, listens to ports, and handles HTTP requests.
- `internal/service` (or `logic`): Services that implement specific business logic and call the data access layer (repository).
- `internal/model` (or `entity`): Data models that define the data structures needed by the business logic layer.
- `internal/repository` (or `dao`): Data access objects that encapsulate database operations and provide CRUD operations on the data.
- `internal/middleware`: Middleware used for handling requests and responses, such as logging, CORS, and signing.

## Dependency Injection

This project utilizes the dependency injection framework `Wire` to achieve modularity and decoupling. `Wire` generates dependency injection code `wire_gen.go` by precompiling `wire.go`, simplifying the process of dependency injection.

- `cmd/job/wire.go`: `Wire` configuration file that defines the dependencies required by the `job` subcommand.
- `cmd/migration/wire.go`: `Wire` configuration file that defines the dependencies required by the `migration` subcommand.
- `cmd/server/wire.go`: `Wire` configuration file that defines the dependencies required by the `server` subcommand.

Wire official documentation: https://github.com/google/wire/blob/main/docs/guide.md

Note: The `wire_gen.go` file is automatically generated during compilation and should not be manually modified.

## Common Code

To achieve code reuse and centralized management, this project adopts a common code approach, where some common code is placed under the `pkg` directory.

- `pkg/config`: Handles reading and parsing configuration files.
- `pkg/helper`: Contains various utility functions, such as MD5 encryption and UUID generation.
- `pkg/http`: Contains HTTP-related code, such as HTTP clients and HTTP servers.
- `pkg/log`: Contains logging-related code, such as log initialization and writing.
- `more...`: Of course, you can freely add and expand more packages as needed.