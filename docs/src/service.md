# Service

在开发 web 项目时，service 通常指一组可被多个处理器或其他服务调用的功能模块。其主要作用是封装业务逻辑，提供可重用的功能，以便在不同上下文中使用。

在 Nunu 项目中，我们设有一个服务目录，用于存放处理特定业务逻辑的服务函数。这些服务函数可以被处理器调用，以执行数据处理操作并返回结果。

## 核心用途

* **业务逻辑封装**: 将复杂的业务逻辑封装在服务中，便于多个处理器或其他服务重用，减少代码重复。
* **数据访问**: 处理与数据库 repository 或外部 API 的交互，提供数据的创建、更新、查询和删除等功能。
* **事务管理**: 在需要时，管理多个操作的事务，确保数据的一致性和完整性。
* **错误处理**: 处理服务内部可能发生的错误，并返回适当的错误信息，使错误在 handler 层被捕获并记录，便于调用者进行相应处理和开发者后续的错误追踪。

## 设计

在 Nunu 中，我们提供了基础的 Service，内置 logger、sid、jwt 和数据库事务等组件，简化全局唯一ID 和事务的管理，方便日志记录与追踪。

```go
type Service struct {
 logger *log.Logger
 sid    *sid.Sid
 jwt    *jwt.JWT
 tm     repository.Transaction
}
```

在创建对应业务的 Service 时，确保在多个业务的 Service 中使用统一的 logger、sid、jwt 和数据库事务，例如 Nunu 中提供的 UserService。

```go
type userService struct {
 *Service
 userRepo repository.UserRepository
}
```

在 UserService 中，通过注入 Service 的方式，统一注入全局唯一ID、数据库事务等组件，简化管理，确保其单例性。
