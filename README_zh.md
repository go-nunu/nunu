# Nunu — A CLI tool for building go aplication.


Nunu是一个基于Golang的应用脚手架，它的名字来自于英雄联盟中的游戏角色，一个骑在雪怪肩膀上的小男孩。和努努一样，该项目也是站在巨人的肩膀上，它是由Golang生态中各种非常流行的库整合而成的，它们的组合可以帮助你快速构建一个高效、可靠的应用程序。

[英文介绍](https://github.com/go-nunu/nunu/blob/main/README.md)

![Nunu](https://github.com/go-nunu/nunu/blob/main/.github/assets/banner.png)


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
- **gocron**:  https://github.com/go-co-op/gocron
- More...
## 特性
* **超低学习成本和定制**：Nunu封装了Gopher最熟悉的一些流行库。您可以轻松定制应用程序以满足特定需求。
* **高性能和可扩展性**：Nunu旨在具有高性能和可扩展性。它使用最新的技术和最佳实践，确保您的应用程序可以处理高流量和大量数据。
* **安全可靠**：Nunu使用了稳定可靠的第三方库，确保您的应用程序安全可靠。
* **模块化和可扩展**：Nunu旨在具有模块化和可扩展性。您可以通过使用第三方库或编写自己的模块轻松添加新功能和功能。
* **文档完善和测试完备**：Nunu文档完善，测试完备。它提供了全面的文档和示例，帮助您快速入门。它还包括一套测试套件，确保您的应用程序按预期工作。

## 简洁分层架构
Nunu采用了经典的分层架构。同时，为了更好地实现模块化和解耦，采用了依赖注入框架`Wire`。

![Nunu Layout](https://github.com/go-nunu/nunu/blob/main/.github/assets/layout.jpg)

## Nunu CLI

![Nunu](https://github.com/go-nunu/nunu/blob/main/.github/assets/screenshot.jpg)


## 文档
* [使用指南](https://github.com/go-nunu/nunu/blob/main/docs/zh/guide.md)
* [分层架构](https://github.com/go-nunu/nunu/blob/main/docs/zh/architecture.md)
* [上手教程](https://github.com/go-nunu/nunu/blob/main/docs/zh/tutorial.md)
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

这是一个经典的Golang 项目的目录结构，包含以下目录：

- `cmd`：存放命令行应用的代码，例如 `main.go`。
- `config`：存放配置文件，例如 `config.yaml`。
- `internal`：存放项目内部的代码，不对外暴露。
    - `dao`：存放数据访问对象（Data Access Object）的代码。
    - `handler`：存放 HTTP 请求处理器的代码。
    - `middleware`：存放 HTTP 中间件的代码。
    - `model`：存放数据模型的代码。
    - `provider`：存放依赖注入的代码。
    - `server`：存放 HTTP 服务器的代码。
    - `service`：存放业务逻辑的代码。
- `pkg`：存放可重用的代码，对外暴露。
    - `config`：存放读取配置文件的代码。
    - `helper`：存放辅助函数的代码。
    - `http`：存放 HTTP 相关的代码。
    - `log`：存放日志相关的代码。

## 要求
要使用Nunu，您需要在系统上安装以下软件：

* Golang 1.16或更高版本
* Git
* MySQL5.7或更高版本(可选)
* Redis（可选）



### 安装

您可以通过以下命令安装Nunu：

```bash
go install github.com/go-nunu/nunu@latest
```

> tips: 如果`go install`成功，却提示找不到nunu命令，这是因为环境变量没有配置，可以把 GOBIN 目录配置到环境变量中即可

### 创建新项目

您可以使用以下命令创建一个新的Golang项目：

```bash
nunu new projectName
```
默认拉取github源，你也可以使用国内加速仓库
```
// 使用基础模板
nunu new projectName -r https://gitee.com/go-nunu/nunu-layout-basic.git
// 使用高级模板
nunu new projectName -r https://gitee.com/go-nunu/nunu-layout-advanced.git
```

此命令将创建一个名为`projectName`的目录，并在其中生成一个优雅的Golang项目结构。

### 创建组件

您可以使用以下命令为项目创建handler、service和dao等组件：

```bash
nunu create handler user
nunu create service user
nunu create dao user
nunu create model user
```
或
```
nunu create all user
```
这些命令将分别创建一个名为`UserHandler`、`UserService`、`UserDao`和`UserModel`的组件，并将它们放置在正确的目录中。

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