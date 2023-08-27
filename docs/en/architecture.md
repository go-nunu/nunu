## Documentation
* [User Guide](https://github.com/go-nunu/nunu/blob/main/docs/en/guide.md)
* [Architecture](https://github.com/go-nunu/nunu/blob/main/docs/en/architecture.md)
* [Getting Started Tutorial](https://github.com/go-nunu/nunu/blob/main/docs/en/tutorial.md)
* [Unit Testing](https://github.com/go-nunu/nunu/blob/main/docs/en/unit_testing.md)


[切换简体中文](https://github.com/go-nunu/nunu/blob/main/docs/zh/architecture.md)

# Exploring the Nunu Architecture

Nunu adopts a classic layered architecture. Additionally, to achieve better modularity and decoupling, it utilizes the dependency injection framework `Wire`.

![Nunu Layout](https://github.com/go-nunu/nunu/blob/main/.github/assets/layout.png)

## Directory Structure
```
.
├── cmd
│   └── server
│       ├── wire
│       │   ├── wire.go
│       │   └── wire_gen.go
│       └── main.go
├── config
│   ├── local.yml
│   └── prod.yml
├── deploy
├── internal
│   ├── handler
│   │   ├── handler.go
│   │   └── user.go
│   ├── job
│   │   └── job.go
│   ├── model
│   │   └── user.go
│   ├── pkg
│   ├── repository
│   │   ├── repository.go
│   │   └── user.go
│   ├── server
│   │   ├── http.go
│   │   └── server.go
│   └── service
│       ├── service.go
│       └── user.go
├── pkg
├── scripts
├── storage
├── test
├── web
├── Makefile
├── go.mod
└── go.sum
```


- `cmd`: Entry point of the application, containing different subcommands.
- `config`: Configuration files.
- `deploy`: Files related to deployment, such as Dockerfile and docker-compose.yml.
- `internal`: Main code of the application, organized according to the layered architecture.
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