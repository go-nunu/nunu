## Documentation
* [User Guide](https://github.com/go-nunu/nunu/blob/main/docs/en/guide.md)
* [Architecture](https://github.com/go-nunu/nunu/blob/main/docs/en/architecture.md)
* [Getting Started Tutorial](https://github.com/go-nunu/nunu/blob/main/docs/en/tutorial.md)

[Go to Chinese version](https://github.com/go-nunu/nunu/blob/main/docs/zh/tutorial.md)

# Nunu Framework User Guide

Nunu is a web framework based on the Go programming language. It provides an elegant project structure and command operations that allow developers to efficiently develop web applications.

## Requirements
To use Nunu with Advanced Layout, you need to have the following software installed on your system:

* Golang 1.16 or higher
* Git
* MySQL 5.7 or higher
* Redis

## Installation

Before you can start using Nunu, you need to install it. You can do so by running the following command:

```bash
go install github.com/go-nunu/nunu@latest
```

For users in China, you can use `GOPROXY` to speed up `go install`.
```
$ go env -w GO111MODULE=on
$ go env -w GOPROXY=https://goproxy.cn,direct
```

> tips: If `go install` is successful but you get an error saying "nunu command not found," it means that the environment variable is not configured. You can add the GOBIN directory to the environment variable.

## Creating a Project

Creating a new project with Nunu is very simple. Just run the following command in the command line:

```bash
nunu new projectName
```

Replace `projectName` with the name of your project. Here, we will choose the Advanced Layout.

**Using an Accelerated Repository in China:**

By default, `nunu new` pulls from the GitHub repository, but you can also use an accelerated repository in China.

```
// Using the advanced template (recommended)
nunu new projectName -r https://gitee.com/go-nunu/nunu-layout-advanced.git

// Using the basic template
nunu new projectName -r https://gitee.com/go-nunu/nunu-layout-basic.git
```

After running the above command, Nunu will automatically create a well-structured Go project with some commonly used files and directories, as shown below:

```
.
├── cmd
│   ├── job
│   │   ├── main.go
│   │   ├── wire.go
│   │   └── wire_gen.go
│   ├── migration
│   │   ├── main.go
│   │   ├── wire.go
│   │   └── wire_gen.go
│   └── server
│       ├── main.go
│       ├── wire.go
│       └── wire_gen.go
├── config
│   ├── local.yml
│   └── prod.yml
├── deploy
│   ├── build
│   │   └── Dockerfile
│   ├── docker-compose
│   │   └── docker-compose.yml
│   └── docker-composer
│       └── conf
│           ├── mysql
│           │   └── conf.d
│           └── redis
│               └── cache
│                   └── redis.conf
├── internal
│   ├── handler
│   │   ├── handler.go
│   │   └── user.go
│   ├── job
│   │   └── job.go
│   ├── middleware
│   │   ├── cors.go
│   │   ├── jwt.go
│   │   ├── log.go
│   │   └── sign.go
│   ├── migration
│   │   └── migration.go
│   ├── model
│   │   └── user.go
│   ├── repository
│   │   ├── repository.go
│   │   └── user.go
│   ├── server
│   │   └── http.go
│   └── service
│       ├── service.go
│       └── user.go
├── mocks
│   ├── repository
│   │   └── user.go
│   └── service
│       └── user.go
├── pkg
│   ├── config
│   │   └── config.go
│   ├── helper
│   │   ├── convert
│   │   │   └── convert.go
│   │   ├── md5
│   │   │   └── md5.go
│   │   ├── resp
│   │   │   └── resp.go
│   │   ├── sid
│   │   │   └── sid.go
│   │   └── uuid
│   │       └── uuid.go
│   ├── http
│   │   └── http.go
│   └── log
│       └── log.go
├── scripts
│   └── deploy.sh
├── storage
│   └── logs
│       └── server.log
├── test
│   └── server
│       ├── handler
│       │   ├── storage
│       │   │   └── logs
│       │   │       └── server.log
│       │   └── user_test.go
│       ├── repository
│       │   ├── repository_test.go
│       │   └── user_test.go
│       └── service
│           └── user_test.go
├── web
│   └── index.html
├── LICENSE
├── Makefile
├── README.md
├── README_zh.md
├── coverage.html
├── go.mod
└── go.sum

```

## Creating Components

In Nunu, you can use the following command to create Handler, Service, Repository, and Model components in batches:

```bash
nunu create all order
```

Here, `order` is the name of the component you want to create.

After running the above command, Nunu will automatically create the components in the corresponding directories and write the corresponding structures and some commonly used methods.
```
// Log information
Created new handler: internal/handler/order.go
Created new service: internal/service/order.go
Created new repository: internal/repository/order.go
Created new model: internal/model/order.go
```

## Registering Routes
Edit `internal/server/http.go`.

Add `*handler.OrderHandler` as a parameter to `NewServerHTTP`, which sets up the dependency for `OrderHandler`.

Next, register a new route: `noAuthRouter.GET("/order", orderHandler.GetOrderById)`.
```
func NewServerHTTP(
	// ...
	orderHandler *handler.OrderHandler,     // new
) *gin.Engine {
    // ...

	// No authentication routes
	noAuthRouter := r.Group("/").Use(middleware.RequestLogMiddleware(logger))
	{
		noAuthRouter.GET("/order", orderHandler.GetOrderById)   // new
```

## Writing Wire Providers
Edit `cmd/server/wire.go` and add the factory functions generated from the files to `providerSet`, as shown below:
```
//go:build wireinject
// +build wireinject

package main

// ...

var HandlerSet = wire.NewSet(
	handler.NewHandler,
	handler.NewUserHandler,

	handler.NewOrderHandler, // new
)

var ServiceSet = wire.NewSet(
	service.NewService,
	service.NewUserService,

	service.NewOrderService, // new
)

var RepositorySet = wire.NewSet(
	repository.NewDB,
	repository.NewRedis,
	repository.NewRepository,
	repository.NewUserRepository,

	repository.NewOrderRepository, // new
)

func newApp(*viper.Viper, *log.Logger) (*gin.Engine, func(), error) {
	panic(wire.Build(
		ServerSet,
		RepositorySet,
		ServiceSet,
		HandlerSet,
		SidSet,
		JwtSet,
	))
}

```
## Compiling Wire

In Nunu, you can use the following command to compile Wire:

```bash
nunu wire
```

After running the above command, you will find that the `wire_gen.go` file is generated from the `cmd/server/wire.go` file.

Open the `wire_gen.go` file, and you will see that the dependency code for `orderRepository`, `orderService`, and `orderHandler` has been automatically generated.

```
func NewApp(viperViper *viper.Viper, logger *log.Logger) (*gin.Engine, func(), error) {
	jwt := middleware.NewJwt(viperViper)
	handlerHandler := handler.NewHandler(logger)
	sidSid := sid.NewSid()
	serviceService := service.NewService(logger, sidSid, jwt)
	db := repository.NewDB(viperViper)
	client := repository.NewRedis(viperViper)
	repositoryRepository := repository.NewRepository(db, client, logger)
	userRepository := repository.NewUserRepository(repositoryRepository)
	userService := service.NewUserService(serviceService, userRepository)
	userHandler := handler.NewUserHandler(handlerHandler, userService)
	
	
	orderRepository := repository.NewOrderRepository(repositoryRepository)
	orderService := service.NewOrderService(serviceService, orderRepository)
	orderHandler := handler.NewOrderHandler(handlerHandler, orderService)
	
	
	engine := server.NewServerHTTP(logger, jwt, userHandler, orderHandler)
	return engine, func() {
	}, nil
}

```

At this point, we have completed the core process of the Nunu project.

Next, you need to modify the MySQL and Redis configuration information in `config/local.yml` and write your logic code in the relevant files.
```
internal/handler/order.go            // Handle request parameters and responses
internal/service/order.go            // Implement business logic
internal/repository/order.go         // Interact with databases and Redis
internal/model/order.go              // Database table entity, GORM model
```

## Starting the Project
Finally, in Nunu, you can use the following command to start the project:

```bash
// Please modify the MySQL and Redis configuration information in config/local.yml before starting the server

// Before starting the server for the first time, run the following database migration
nunu run ./cmd/migration  

// Start the server
nunu run ./cmd/server    

// Or

nunu run
```

After running the above command, Nunu will automatically start the project and monitor file updates, supporting hot-reloading.

## Conclusion

The Nunu framework provides an elegant project structure and command operations that allow developers to efficiently develop web applications. In this tutorial, you have learned how to create a project, create Handlers, create Services, create Repositories, compile Wire, and start the project using Nunu. We hope that this content will help you make better use of the Nunu framework.