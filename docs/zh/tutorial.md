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

## 创建项目

使用Nunu创建一个新项目非常简单，只需要在命令行中输入以下命令：

```bash
nunu new projectName
```
其中`projectName`是你想要创建的项目名称,这里我们选择高级Layout

**国内加速源：**

`nunu new`默认拉取github源，你也可以使用国内加速仓库
```
// 使用基础模板
nunu new projectName -r https://gitee.com/go-nunu/nunu-layout-basic.git
// 使用高级模板(推荐)
nunu new projectName -r https://gitee.com/go-nunu/nunu-layout-advanced.git
```


执行完上述命令后，Nunu会自动创建一个目录结构优雅的Go项目，包含了一些常用的文件和目录，如下所示：

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

## 创建组件

在Nunu中，可以使用以下命令批量创建Handler、Service、Dao、Model组件：

```bash
nunu create all order
```

其中，`order`是你想要创建的组件名称。

执行完上述命令后，Nunu会自动在对应目录创建组件，并写入对应的结构体和一些常用的方法。
```
// 日志信息
Created new handler: internal/handler/order.go
Created new service: internal/service/order.go
Created new dao: internal/dao/order.go
Created new model: internal/model/order.go
```

## 注册路由
编辑 `internal/server/http.go`

将` *handler.OrderHandler`添加为`NewServerHTTP`的参数，这样就写好了`OrderHandler`的依赖关系。

紧接着我们我们再注册一个路由，`noAuthRouter.GET("/order", orderHandler.GetOrderById)`
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

## 编写Wire Provider
编辑 `internal/provider/provider.go`，将刚刚生成的文件中的工厂函数添加到`providerSet`中，如下所示：
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
## 编译Wire

在Nunu中，可以使用以下命令编译Wire：

```bash
nunu wire
```

执行完上述命令后，我们选择`cmd/server/wire/wire.go`文件，生成对应的`wire_gen.go`文件。


打开`cmd/server/wire/wire_gen.go`文件，我们可以看到`orderDao`、`orderService`、`orderHandler`的依赖关系代码自动生成了。

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

至此，我们已经走完了Nunu项目中的核心流程，

接下来，你需要修改`config/local.yml`中的Mysql和Redis配置信息，

并在相关的文件中编写你的逻辑代码即可。
```
internal/handler/order.go   // 处理请求参数和响应
internal/service/order.go   // 实现业务逻辑
internal/dao/order.go       // 与数据库和Redis等交互
internal/model/order.go     // 数据表实体，GORM model
```

## 启动项目
最后，在Nunu中，可以使用以下命令启动项目：

```bash
nunu run
```

执行完上述命令后，Nunu会自动启动项目，并监听文件更新，支持热重启。

## 总结

Nunu框架提供了一套优雅的项目结构和命令操作，使得开发者可以更加高效地开发Web应用程序。通过本教程，你已经学会了如何使用Nunu创建项目、创建Handler、创建Service、创建Dao、编译Wire和启动项目。希望这些内容能够帮助你更好地使用Nunu框架。