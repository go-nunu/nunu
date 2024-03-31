# Server基础概念
在Nunu中，我们将`HTTP`、`GRPC`、`WebSocket`、`Task`、`Job`等服务都抽象为`Server`。

```go
type Server interface {
	Start(context.Context) error
	Stop(context.Context) error
}
```

每个`Server`都必须实现`Server`接口中的方法，也就是`Start(ctx)`和`Stop(ctx)`




## 服务依赖注册

在Nunu中，如果你想给你的某个服务进程注册相应的启动服务，

你只需要关心`cmd/[服务器名称]/wire/wire.go`文件即可。

```go
var serverSet = wire.NewSet(
	server.NewHTTPServer,
	server.NewJob,
	// 你想注册的其它服务
)

// build App
func newApp(
	httpServer *http.Server, 
	job *server.Job,
    // 你想注册的其它服务,需要从参数传入
	) *app.App {
	return app.NewApp(
		app.WithServer(httpServer, job),
		app.WithName("demo-server"),
	)
}

func NewWire(*viper.Viper, *log.Logger) (*app.App, func(), error) {

	panic(wire.Build(
		repositorySet,
		serviceSet,
		handlerSet,
		serverSet,  // 集中注册服务
		sid.NewSid,
		jwt.NewJwt,
		newApp,
	))
}


```

