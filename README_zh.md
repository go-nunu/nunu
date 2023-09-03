# Nunu — A CLI tool for building go aplication.


Nunu是一个基于Golang的应用脚手架，它的名字来自于英雄联盟中的游戏角色，一个骑在雪怪肩膀上的小男孩。和努努一样，该项目也是站在巨人的肩膀上，它是由Golang生态中各种非常流行的库整合而成的，它们的组合可以帮助你快速构建一个高效、可靠的应用程序。

[英文介绍](https://github.com/go-nunu/nunu/blob/main/README.md)

![Nunu](https://github.com/go-nunu/nunu/blob/main/.github/assets/banner.png)

## 文档
* [使用指南](https://github.com/go-nunu/nunu/blob/main/docs/zh/guide.md)
* [分层架构](https://github.com/go-nunu/nunu/blob/main/docs/zh/architecture.md)
* [详细教程](https://github.com/go-nunu/nunu/blob/main/docs/zh/tutorial.md)
* [高效编写单元测试](https://github.com/go-nunu/nunu/blob/main/docs/zh/unit_testing.md)

## 功能
- **Gin**: https://github.com/gin-gonic/gin
- **Gorm**: https://github.com/go-gorm/gorm
- **Wire**: https://github.com/google/wire
- **Viper**: https://github.com/spf13/viper
- **Zap**: https://github.com/uber-go/zap
- **Golang-jwt**: https://github.com/golang-jwt/jwt
- **Go-redis**: https://github.com/go-redis/redis
- **Testify**: https://github.com/stretchr/testify
- **Sonyflake**: https://github.com/sony/sonyflake
- **Gocron**:  https://github.com/go-co-op/gocron
- **Go-sqlmock**:  https://github.com/DATA-DOG/go-sqlmock
- **Gomock**:  https://github.com/golang/mock
- **Swaggo**:  https://github.com/swaggo/swag
- More...
## 特性
* **超低学习成本和定制**：Nunu封装了Gopher最熟悉的一些流行库。您可以轻松定制应用程序以满足特定需求。
* **高性能和可扩展性**：Nunu旨在具有高性能和可扩展性。它使用最新的技术和最佳实践，确保您的应用程序可以处理高流量和大量数据。
* **安全可靠**：Nunu使用了稳定可靠的第三方库，确保您的应用程序安全可靠。
* **模块化和可扩展**：Nunu旨在具有模块化和可扩展性。您可以通过使用第三方库或编写自己的模块轻松添加新功能和功能。
* **文档完善和测试完备**：Nunu文档完善，测试完备。它提供了全面的文档和示例，帮助您快速入门。它还包括一套测试套件，确保您的应用程序按预期工作。

## 交流群组

微信入群，请备注Nunu

<p align="left"><img src="https://github.com/go-nunu/nunu/blob/main/.github/assets/qrcode.jpg" width="200"></p>

## 简洁分层架构
Nunu采用了经典的分层架构。同时，为了更好地实现模块化和解耦，采用了依赖注入框架`Wire`。

![Nunu Layout](https://github.com/go-nunu/nunu/blob/main/.github/assets/layout.png)

## Nunu CLI

![Nunu](https://github.com/go-nunu/nunu/blob/main/.github/assets/screenshot.jpg)



## 目录结构
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

该项目的架构采用了典型的分层架构，主要包括以下几个模块：

* `cmd`：该模块包含了应用的入口点，根据不同的命令进行不同的操作，例如启动服务器、执行数据库迁移等。每个子模块都有一个`main.go`文件作为入口文件，以及`wire.go`和`wire_gen.go`文件用于依赖注入。
* `config`：该模块包含了应用的配置文件，根据不同的环境（如开发环境和生产环境）提供不同的配置。
* `deploy`：该模块用于部署应用，包含了一些部署脚本和配置文件。
* `internal`：该模块是应用的核心模块，包含了各种业务逻辑的实现。

  - `handler`：该子模块包含了处理HTTP请求的处理器，负责接收请求并调用相应的服务进行处理。

  - `job`：该子模块包含了后台任务的逻辑实现。

  - `model`：该子模块包含了数据模型的定义。

  - `repository`：该子模块包含了数据访问层的实现，负责与数据库进行交互。

  - `server`：该子模块包含了HTTP服务器的实现。

  - `service`：该子模块包含了业务逻辑的实现，负责处理具体的业务操作。

* `pkg`：该模块包含了一些通用的功能和工具。

* `scripts`：该模块包含了一些脚本文件，用于项目的构建、测试和部署等操作。

* `storage`：该模块用于存储文件或其他静态资源。

* `test`：该模块包含了各个模块的单元测试，按照模块划分子目录。

* `web`：该模块包含了前端相关的文件，如HTML、CSS和JavaScript等。

此外，还包含了一些其他的文件和目录，如授权文件、构建文件、README等。整体上，该项目的架构清晰，各个模块之间的职责明确，便于理解和维护。

## 要求
要使用Nunu，您需要在系统上安装以下软件：

* Golang 1.19或更高版本
* Git
* Docker (可选)
* MySQL5.7或更高版本(可选)
* Redis（可选）



### 安装

您可以通过以下命令安装Nunu：

```bash
go install github.com/go-nunu/nunu@latest
```

国内用户可以使用`GOPROXY`加速`go install`

```
$ go env -w GO111MODULE=on
$ go env -w GOPROXY=https://goproxy.cn,direct
```

> tips: 如果`go install`成功，却提示找不到nunu命令，这是因为环境变量没有配置，可以把 GOBIN 目录配置到环境变量中即可


### 创建新项目

您可以使用以下命令创建一个新的Golang项目：

```bash
// 推荐新用户选择Advanced Layout
nunu new projectName
```

`nunu new`默认拉取github源，你也可以使用国内加速仓库
```
// 使用高级模板(推荐)
nunu new projectName -r https://gitee.com/go-nunu/nunu-layout-advanced.git

// 使用基础模板
nunu new projectName -r https://gitee.com/go-nunu/nunu-layout-basic.git

```


> Nunu内置了两种类型的Layout：

* **基础模板(Basic Layout)**

Basic Layout 包含一个非常精简的架构目录结构，适合非常熟悉Nunu项目的开发者使用。

* **高级模板(Advanced Layout)**

**建议：我们推荐新手优先选择使用Advanced Layout。**


Advanced Layout 包含了很多Nunu的用法示例（ db、redis、 jwt、 cron、 migration等），适合开发者快速学习了解Nunu的架构思想。

此命令将创建一个名为`projectName`的目录，并在其中生成一个优雅的Golang项目结构。

### 创建组件

您可以使用以下命令为项目创建handler、service、repository和model等组件：

```bash
nunu create handler user
nunu create service user
nunu create repository user
nunu create model user
```
或
```
nunu create all user
```
这些命令将分别创建一个名为`UserHandler`、`UserService`、`UserRepository`和`UserModel`的组件，并将它们放置在正确的目录中。

### 启动项目

您可以使用以下命令快速启动项目：

```bash
nunu run
```

此命令将启动您的Golang项目，并支持文件更新热重启。

### 编译wire.go

您可以使用以下命令快速编译`wire.go`：

```bash
nunu wire
```

此命令将编译您的`wire.go`文件，并生成所需的依赖项。

## 贡献

如果您发现任何问题或有任何改进意见，请随时提出问题或提交拉取请求。我们非常欢迎您的贡献！

## 许可证

Nunu是根据MIT许可证发布的。有关更多信息，请参见[LICENSE](LICENSE)文件。
