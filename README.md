# Nunu — A CLI tool for building go aplication.

Nunu is an application scaffold based on Golang, named after a game character in League of Legends, a little boy riding on the shoulder of a snow monster. Like Nunu, this project also stands on the shoulders of giants, integrating various popular libraries in the Golang ecosystem. Their combination can help you quickly build an efficient and reliable application.

[中文介绍](https://github.com/go-nunu/nunu/blob/main/README_zh.md)


![Nunu](https://github.com/go-nunu/nunu/blob/main/.github/assets/banner.png)







## Features

- **Gin**: https://github.com/gin-gonic/gin
- **Gorm**: https://github.com/go-gorm/gorm
- **Wire**: https://github.com/google/wire
- **Viper**: https://github.com/spf13/viper
- **Zap**: https://github.com/uber-go/zap
- **Golang-jwt**: https://github.com/golang-jwt/jwt
- **Go-redis**: https://github.com/go-redis/redis
- **Testify**: https://github.com/stretchr/testify
- **Sonyflake**: https://github.com/sony/sonyflake
- **robfig-cron**: https://github.com/robfig/cron
- More...

## Features
* **Low learning cost and customization**: Nunu encapsulates some popular libraries that Gopher is most familiar with. You can easily customize your application to meet specific needs.
* **High performance and scalability**: Nunu aims to have high performance and scalability. It uses the latest technology and best practices to ensure that your application can handle high traffic and large amounts of data.
* **Secure and reliable**: Nunu uses stable and reliable third-party libraries to ensure the security and reliability of your application.
* **Modular and extensible**: Nunu is designed to be modular and extensible. You can easily add new features and functionality by using third-party libraries or writing your own modules.
* **Complete documentation and testing**: Nunu has complete documentation and testing. It provides comprehensive documentation and examples to help you get started quickly. It also includes a set of test suites to ensure that your application works as expected.

## Nunu CLI

![Nunu](https://github.com/go-nunu/nunu/blob/main/.github/assets/screenshot.jpg)


## Documentation
* [Guide](https://github.com/go-nunu/nunu/blob/main/docs/en/guide.md)
* [Architecture](https://github.com/go-nunu/nunu/blob/main/docs/en/architecture.md)
* [Tutorial](https://github.com/go-nunu/nunu/blob/main/docs/en/tutorial.md)


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
├── internal
│   ├── dao
│   │   ├── dao.go
│   │   └── user.go
│   ├── handler
│   │   ├── handler.go
│   │   └── user.go
│   ├── middleware
│   │   └── cors.go
│   ├── model
│   │   └── user.go
│   ├── provider
│   │   └── provider.go
│   ├── server
│   │   └── http.go
│   └── service
│       ├── service.go
│       └── user.go
├── pkg
│   ├── config
│   │   └── config.go
│   ├── helper
│   │   ├── md5
│   │   │   └── md5.go
│   │   ├── resp
│   │   │   └── resp.go
│   │   ├── sonyflake
│   │   │   └── sonyflake.go
│   │   └── uuid
│   │       └── uuid.go
│   ├── http
│   │   └── http.go
│   └── log
│       └── log.go
├── LICENSE
├── README.md
├── README_zh.md
├── go.mod
└── go.sum
```


This is the directory structure of a classic Golang project, which includes the following directories:

- `cmd`: Contains the code for command-line applications, such as `main.go`.
- `config`: Contains configuration files, such as `config.yaml`.
- `internal`: Contains internal code that is not exposed externally.
    - `dao`: Contains the code for Data Access Objects (DAOs).
    - `handler`: Contains the code for HTTP request handlers.
    - `middleware`: Contains the code for HTTP middleware.
    - `model`: Contains the code for data models.
    - `provider`: Contains the code for dependency injection.
    - `server`: Contains the code for HTTP servers.
    - `service`: Contains the code for business logic.
- `pkg`: Contains reusable code that is exposed externally.
    - `config`: Contains the code for reading configuration files.
    - `helper`: Contains the code for helper functions.
    - `http`: Contains HTTP-related code.
    - `log`: Contains code related to logging.

## Requirements
To use Nunu, you need to install the following software on your system:

* Golang 1.16 or higher
* Git
* MySQL 5.7 or higher (optional)
* Redis (optional)

### Installation

You can install Nunu using the following command:

```bash
go install github.com/go-nunu/nunu@latest
```


### Creating a New Project

You can create a new Golang project using the following command:

```bash
nunu new projectName

// or

nunu new projectName -r https://github.com/go-nunu/nunu-layout-advanced.git
```

This command will create a directory named `projectName` and generate an elegant Golang project structure within it.

### Creating Components

You can create handlers, services, and daos for your project using the following commands:

```bash
nunu create handler user
nunu create service user
nunu create dao user
nunu create model user
```
or
```
nunu create all user
```

These commands will create components named `UserHandler`, `UserService`, `UserDao` and `UserModel`, respectively, and place them in the correct directories.

### Starting the Project

You can quickly start your project using the following command:

```bash
nunu run
```

This command will start your Golang project and support file update hot reload.

### Compiling wire.go

You can quickly compile your `wire.go` file using the following command:

```bash
nunu wire
```

This command will compile your `wire.go` file and generate the required dependencies.

## Contributing

If you find any issues or have any improvement suggestions, please feel free to raise an issue or submit a pull request. We welcome your contributions!

## License

Nunu is released under the MIT license. See [LICENSE](LICENSE) for more information.