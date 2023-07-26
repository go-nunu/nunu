## 文档
* [使用指南](https://github.com/go-nunu/nunu/blob/main/docs/zh/guide.md)
* [分层架构](https://github.com/go-nunu/nunu/blob/main/docs/zh/architecture.md)
* [上手教程](https://github.com/go-nunu/nunu/blob/main/docs/zh/tutorial.md)
* [高效编写单元测试](https://github.com/go-nunu/nunu/blob/main/docs/zh/unit_testing.md)


[进入英文版](https://github.com/go-nunu/nunu/blob/main/docs/en/tutorial.md)

# Nunu框架使用教程

Nunu是一个基于Go语言的Web框架，它提供了一套优雅的项目结构和命令操作，使得开发者可以更加高效地开发Web应用程序。


## 要求
要使用Nunu 高级Layout，您需要在系统上安装以下软件：

* Golang 1.16或更高版本
* Git
* MySQL5.7或更高版本
* Redis

## 安装

在开始使用Nunu之前，需要先安装它。可以通过以下命令进行安装：



```bash

go install github.com/go-nunu/nunu@latest
```

国内用户可以使用`GOPROXY`加速`go install`
```
$ go env -w GO111MODULE=on
$ go env -w GOPROXY=https://goproxy.cn,direct
```

> tips: 如果`go install`成功，却提示找不到nunu命令，这是因为环境变量没有配置，可以把 GOBIN 目录配置到环境变量中即可



## 创建项目

使用Nunu创建一个新项目非常简单，只需要在命令行中输入以下命令：

```bash
nunu new projectName
```
其中`projectName`是你想要创建的项目名称,**这里我们选择Advanced Layout**

**国内加速源：**

`nunu new`默认拉取github源，你也可以使用国内加速仓库
```
// 使用高级模板(推荐)
nunu new projectName -r https://gitee.com/go-nunu/nunu-layout-advanced.git

// 使用基础模板
nunu new projectName -r https://gitee.com/go-nunu/nunu-layout-basic.git
```


执行完上述命令后，Nunu会自动创建一个目录结构优雅的Go项目，包含了一些常用的文件和目录，如下所示：

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

## 创建组件

在Nunu中，可以使用以下命令批量创建Handler、Service、Repository、Model组件：

```bash
nunu create all order
```

其中，`order`是你想要创建的组件名称。

执行完上述命令后，Nunu会自动在对应目录创建组件，并写入对应的结构体和一些常用的方法。
```
// 日志信息
Created new handler: internal/handler/order.go
Created new service: internal/service/order.go
Created new repository: internal/repository/order.go
Created new model: internal/model/order.go
```

## 注册路由
编辑 `internal/server/http.go`

将` *handler.OrderHandler`添加为`NewServerHTTP`的参数，这样就写好了`OrderHandler`的依赖关系。

紧接着我们我们再注册一个路由，`noAuthRouter.GET("/order", orderHandler.GetOrderById)`
```
func NewServerHTTP(
	// ...
	orderHandler handler.OrderHandler,     // new
) *gin.Engine {
    // ...

	// 无权限路由
	noAuthRouter := r.Group("/").Use(middleware.RequestLogMiddleware(logger))
	{
		noAuthRouter.GET("/order", orderHandler.GetOrderById)   // new
```

## 编写Wire Provider
编辑 `cmd/server/wire.go`，将刚刚生成文件中的工厂函数添加到`providerSet`中，如下所示：
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
## 编译Wire

在Nunu中，可以使用以下命令编译Wire：

```bash
nunu wire
```

执行完上述命令后，我们选择`cmd/server/wire.go`文件，生成对应的`wire_gen.go`文件。


打开`cmd/server/wire_gen.go`文件，我们可以看到`orderRepository`、`orderService`、`orderHandler`的依赖关系代码自动生成了。

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

至此，我们已经走完了Nunu项目中的核心流程，

接下来，你需要修改`config/local.yml`中的Mysql和Redis配置信息，

并在相关的文件中编写你的逻辑代码即可。
```
internal/handler/order.go            // 处理请求参数和响应
internal/service/order.go            // 实现业务逻辑
internal/repository/order.go         // 与数据库和Redis等交互
internal/model/order.go              // 数据表实体，GORM model
```

## 启动项目
最后，在Nunu中，可以使用以下命令启动项目：

```bash
// 请先修改config/local.yml中的 MySQL 和 Redis 配置信息

// 初次启动server之前，请先执行以下数据库迁移
nunu run ./cmd/migration  

 // 启动server
nunu run ./cmd/server    

// 或

nunu run
```

执行完上述命令后，Nunu会自动启动项目，并监听文件更新，支持热重启。

## 总结

Nunu框架提供了一套优雅的项目结构和命令操作，使得开发者可以更加高效地开发Web应用程序。通过本教程，你已经学会了如何使用Nunu创建项目、创建Handler、创建Service、创建Repository、编译Wire和启动项目。希望这些内容能够帮助你更好地使用Nunu框架。
