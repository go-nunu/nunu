# 依赖注入
## 为什么需要依赖注入？{#why-do-we-need-dependency-injection}

依赖注入是一种编程模式，用于管理代码中各个组件之间的依赖关系。在没有依赖注入的情况下，通常会使用全局变量或硬编码的方式来获取所需的依赖，这种方式会带来一些问题。

## 全局变量{#global-variable-issues}
让我们通过一个简单的 Go 语言示例来说明：

假设我们有一个简单的服务，它需要一个数据库连接。在没有依赖注入的情况下，我们可能会这样实现：

```go
package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	// 初始化数据库连接
	db, _ = sql.Open("mysql", "user:password@tcp(localhost:3306)/dbname")
}

func main() {
	// 使用全局变量 db
	// 这里可能会有很多其他的逻辑
}

```
在上面的代码中，我们使用了全局变量 `db` 来保存数据库连接。这样做可能会引发以下问题：

**并发安全问题**：多个并发执行的程序可能会同时访问和修改全局变量，这可能导致竞态条件和数据竞争，从而引发难以调试的错误。在并发环境中，全局变量的状态不受控制，很容易引发意外的行为。

**难以测试**： 全局变量使得单元测试变得困难，因为在测试过程中很难控制全局变量的状态。这可能导致测试覆盖不全或者需要编写更多的集成测试来覆盖全局变量的不同状态。

**不可控的副作用**： 全局变量的值可以在程序的任何地方被修改，这可能导致不可预测的副作用。例如，一个函数可能会依赖于某个全局变量的值，但是在调用该函数之前，另一个部分的代码可能会修改该全局变量，从而改变函数的行为。

**可维护性差**： 全局变量增加了代码的耦合性，使得代码难以重构和维护。当代码规模增大时，全局变量会导致代码难以理解和修改，降低了代码的可维护性。

## 手动依赖注入 {#manual-dependency-injection}
依赖注入可以帮助解决这些问题。在依赖注入中，依赖关系不再硬编码到组件内部，而是通过构造函数、方法参数或者其他方式传递进来。

下面是使用依赖注入改进上面示例的代码：

```go
package main

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

type Service struct {
    db *sql.DB
}

func NewService(db *sql.DB) *Service {
    return &Service{db}
}

func (s *Service) DoSomething() {
    // 使用 s.db 来执行数据库操作
}

func main() {
    // 初始化数据库连接
    db, _ := sql.Open("mysql", "user:password@tcp(localhost:3306)/dbname")

    // 创建服务实例，并将数据库连接传递进去
    service := NewService(db)

    // 使用服务进行操作
    service.DoSomething()
}

```

在上面的代码中，我们通过将 `db` 作为参数传递给 `NewService` 函数，实现了依赖注入。这样做带来了以下好处：

**解耦：** `Service` 和数据库连接之间的关系解耦，使得 `Service` 更易于重用和测试。

**易于测试：** 现在我们可以轻松地通过传递模拟的数据库连接来测试 `Service` 的行为。

**可维护性：** 通过明确地传递依赖关系，代码变得更易于理解和维护。

## 自动依赖注入{#automatic-dependency-injection}
**当项目变得庞大时，手动进行依赖注入可能会变得繁琐。**

如下是一个真实项目的部分代码，它使用了依赖注入来管理各个组件之间的依赖关系。

```go 
jwtJWT := jwt.NewJwt(viperViper)
handlerHandler := handler.NewHandler(logger)
sidSid := sid.NewSid()
db := repository.NewDB(viperViper, logger)
client := repository.NewRedis(viperViper)
officialAccount, err := repository.NewWechatOfficial(viperViper)
if err != nil {
    return nil, nil, err
}
clientV3, err := repository.NewWechatPay(viperViper)
if err != nil {
    return nil, nil, err
}
cosClient := repository.NewCosClient(viperViper)
repositoryRepository := repository.NewRepository(db, client, logger, officialAccount, clientV3, viperViper, sidSid, cosClient)
serviceService := service.NewService(logger, viperViper, sidSid, jwtJWT, repositoryRepository)
userRepository := repository.NewUserRepository(repositoryRepository)
walletRepository := repository.NewWalletRepository(repositoryRepository)
cosRepository := repository.NewCosRepository(repositoryRepository)
wechatRepository := repository.NewWechatRepository(repositoryRepository)
captchaRepository := repository.NewCaptchaRepository(repositoryRepository)
userService := service.NewUserService(serviceService, userRepository, walletRepository, cosRepository, wechatRepository, captchaRepository)
userHandler := handler.NewUserHandler(handlerHandler, userService)
categoryRepository := repository.NewCategoryRepository(repositoryRepository)
categoryService := service.NewCategoryService(serviceService, categoryRepository)
categoryHandler := handler.NewCategoryHandler(handlerHandler, categoryService)
fileRepository := repository.NewFileRepository(repositoryRepository)
fileService := service.NewFileService(serviceService, fileRepository, cosRepository)
uploadHandler := handler.NewUploadHandler(handlerHandler, viperViper, fileService)
promptRepository := repository.NewPromptRepository(repositoryRepository)
promptService := service.NewPromptService(serviceService, promptRepository, fileRepository, walletRepository, categoryRepository, userRepository, cosRepository, viperViper)
promptHandler := handler.NewPromptHandler(handlerHandler, promptService)
wechatService := service.NewWechatService(serviceService, wechatRepository, userRepository, walletRepository)
wechatHandler := handler.NewWechatHandler(handlerHandler, wechatService, userService)
orderRepository := repository.NewOrderRepository(repositoryRepository)
walletService := service.NewWalletService(serviceService, viperViper, walletRepository, orderRepository, promptRepository, clientV3)
walletHandler := handler.NewWalletHandler(handlerHandler, walletService)
translateService := service.NewTranslateService(serviceService, viperViper)
translateHandler := handler.NewTranslateHandler(handlerHandler, translateService)
captchaService := service.NewCaptchaService(serviceService, captchaRepository)
authHandler := handler.NewAuthHandler(handlerHandler, userService, captchaService)
httpServer := server.NewHTTPServer(logger, viperViper, jwtJWT, userHandler, categoryHandler, uploadHandler, promptHandler, wechatHandler, walletHandler, translateHandler, authHandler)
job := server.NewJob(logger)
appApp := newApp(httpServer, job)
```

看到这样的代码，相信你一定会感觉依赖注入使我们代码使变得复杂，但其实这段代码并不是人工手写的，而是通过`wire`自动生成的。

我们需要做的仅仅是需要声明我们的依赖关系，执行`wire`命令后，就会自动帮我们生成代码。

**声明依赖关系的代码非常简单，代码大概像下面这样**
```
var handlerSet = wire.NewSet(
	handler.NewHandler,
	handler.NewUserHandler,
	handler.NewCategoryHandler,
)

var serviceSet = wire.NewSet(
	service.NewService,
	service.NewUserService,
)

var repositorySet = wire.NewSet(
	repository.NewDB,
	repository.NewRedis,
	repository.NewWechatOfficial,
	repository.NewWechatPay,
	repository.NewRepository,
	repository.NewUserRepository,
)
var serverSet = wire.NewSet(
	server.NewHTTPServer,
	server.NewJob,
	server.NewTask,
)

// build App
func newApp(httpServer *http.Server, job *server.Job) *app.App {
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
		serverSet,
		sid.NewSid,
		jwt.NewJwt,
		newApp,
	))
}
```

<div class="tip custom-block" style="padding-top: 8px">

想学习更多关于Wire知识？前往到[Wire官方文档](https://github.com/google/wire/blob/main/docs/guide.md)。

</div>
