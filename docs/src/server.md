# Server基础概念
在Nunu中，我们将`HTTP`、`GRPC`、`WebSocket`、`Task`、`Job`等服务都抽象为`Server`。

```go
type Server interface {
	Start(context.Context) error
	Stop(context.Context) error
}
```

每个`Server`都必须实现`Server`接口中的方法，也就是`Start(ctx)`和`Stop(ctx)`

## HTTP
HTTP服务，我们使用`gin`作为HTTP框架，`gin`的`Engine`实现了`Server`接口，因此，我们只需要将`Engine`作为`Server`即可。

## Task
Task服务，我们使用`cron`作为Task框架，`cron`的`Cron`实现了`Server`接口，因此，我们只需要将`Cron`作为`Server`即可。

## Job
Job服务，我们使用`cron`作为Job框架，`cron`的`Cron`实现了`Server`接口，因此，我们只需要将`Cron`作为`Server`即可。

## Migration
Migration服务，我们使用`migrate`作为Migration框架，`migrate`的`Migrate`实现了`Server`接口，因此，我们只需要将`Migrate`作为`Server`即可。

