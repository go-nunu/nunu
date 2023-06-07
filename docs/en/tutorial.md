## Documentation
* [Guide](https://github.com/go-nunu/nunu/blob/main/docs/en/guide.md)
* [Architecture](https://github.com/go-nunu/nunu/blob/main/docs/en/architecture.md)
* [Tutorial](https://github.com/go-nunu/nunu/blob/main/docs/en/tutorial.md)

[进入简体中文版](https://github.com/go-nunu/nunu/blob/main/docs/zh/tutorial.md)

# Nunu Framework Tutorial

Nunu is a web framework based on the Go language. It provides an elegant project structure and command operations, making it easier for developers to develop web applications more efficiently.

## Requirements
To use Nunu's advanced layout, you need to install the following software on your system:

* Golang 1.16 or higher
* Git
* MySQL 5.7 or higher
* Redis

## Installation

Before using Nunu, you need to install it. You can install it using the following command:

```bash
go install github.com/go-nunu/nunu@latest
```

## Creating a Project

Creating a new project with Nunu is very simple. Just enter the following command in the command line:

```bash
nunu new projectName
```

Here, `projectName` is the name of the project you want to create. Here we choose the advanced layout.

After executing the above command, Nunu will automatically create a directory structure for the Go project, including some commonly used files and directories, as shown below:

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
│   ├── http
│   │   └── http.go
│   └── log
│       └── log.go
├── script
│   └── deploy.sh
├── storage
├── test
├── web
├── LICENSE
├── README.md
├── README_zh.md
├── go.mod
└── go.sum
```

## Creating Components

In Nunu, you can use the following command to create Handler, Service, Dao, and Model components in batches:

```bash
nunu create all order
```

Here, `order` is the name of the component you want to create.

After executing the above command, Nunu will automatically create the components in the corresponding directory and write the corresponding structures and some commonly used methods.

```
// Log information
Created new handler: internal/handler/order.go
Created new service: internal/service/order.go
Created new dao: internal/dao/order.go
Created new model: internal/model/order.go
```

## Registering Routes

Edit `internal/server/http.go` and add `*handler.OrderHandler` as a parameter to `NewServerHTTP`, so that the dependency relationship of `OrderHandler` is written.

Then we register a route, `noAuthRouter.GET("/order", orderHandler.GetOrderById)`.
```
func NewServerHTTP(
	// ...
	orderHandler *handler.OrderHandler,
) *gin.Engine {
    // ...

	// 无权限路由
	noAuthRouter := r.Group("/").Use(middleware.RequestLogMiddleware(logger))
	{
		noAuthRouter.GET("/order", orderHandler.GetOrderById)
```

## Writing Wire Provider

Edit `internal/provider/provider.go` and add the factory function generated in the file you just generated to `providerSet`, as shown below:

```

var DaoSet = wire.NewSet(
	dao.NewDB,
	dao.NewRedis,
	dao.NewDao,
	dao.NewUserDao,

	dao.NewOrderDao, // new
)

var ServiceSet = wire.NewSet(
	service.NewService,
	service.NewUserService,

	service.NewOrderService, // new
)

var HandlerSet = wire.NewSet(
	handler.NewHandler,
	handler.NewUserHandler,

	handler.NewOrderHandler, // new
)
```
## Compiling Wire

In Nunu, you can use the following command to compile Wire:

```bash
nunu wire
```

After executing the above command, we select the `cmd/server/wire/wire.go` file, and the corresponding `wire_gen.go` file is generated.

Open the `cmd/server/wire/wire_gen.go` file, and you can see that the dependency relationship code of `orderDao`, `orderService`, and `orderHandler` is automatically generated.

```
func NewApp(viperViper *viper.Viper, logger *log.Logger) (*gin.Engine, func(), error) {
	jwt := middleware.NewJwt(viperViper)
	sonyflakeSonyflake := sonyflake.NewSonyflake()
	handlerHandler := handler.NewHandler(logger, sonyflakeSonyflake)
	serviceService := service.NewService(logger)
	db := dao.NewDB(viperViper)
	client := dao.NewRedis(viperViper)
	daoDao := dao.NewDao(db, client, logger)
	userDao := dao.NewUserDao(daoDao)
	userService := service.NewUserService(serviceService, userDao)
	userHandler := handler.NewUserHandler(handlerHandler, userService)
	
	
	orderDao := dao.NewOrderDao(daoDao)
	orderService := service.NewOrderService(serviceService, orderDao)
	orderHandler := handler.NewOrderHandler(handlerHandler, orderService)
	engine := server.NewServerHTTP(logger, jwt, userHandler, orderHandler)
	
	
	return engine, func() {
	}, nil
}

```

At this point, we have completed the core process of the Nunu project.

Next, you need to modify the Mysql and Redis configuration information in `config/local.yml` and write your logic code in the relevant files.

```
internal/handler/order.go   // handle request parameters and responses
internal/service/order.go   // implement business logic
internal/dao/order.go       // interact with databases and Redis
internal/model/order.go     // data table entity, GORM model
```

## Starting the Project
Finally, in Nunu, you can use the following command to start the project:

```bash
nunu run
```

After executing the above command, Nunu will automatically start the project and listen for file updates, supporting hot restart.

## Conclusion

The Nunu framework provides an elegant project structure and command operations, making it easier for developers to develop web applications more efficiently. Through this tutorial, you have learned how to use Nunu to create projects, create Handlers, create Services, create Dao, compile Wire, and start projects. I hope these contents can help you better use the Nunu framework.