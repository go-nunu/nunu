## 文档
* [使用指南](https://github.com/go-nunu/nunu/blob/main/docs/zh/guide.md)
* [分层架构](https://github.com/go-nunu/nunu/blob/main/docs/zh/architecture.md)
* [上手教程](https://github.com/go-nunu/nunu/blob/main/docs/zh/tutorial.md)
* [高效编写单元测试](https://github.com/go-nunu/nunu/blob/main/docs/zh/unit_testing.md)


[进入英文版](https://github.com/go-nunu/nunu/blob/main/docs/en/architecture.md)

# Nunu架构详解

Nunu采用了经典的分层架构。同时，为了更好地实现模块化和解耦，采用了依赖注入框架`Wire`。

![Nunu Layout](https://github.com/go-nunu/nunu/blob/main/.github/assets/layout.png)

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


- `cmd`：应用程序的入口，包含了不同的子命令。
- `config`：配置文件。
- `deploy`：部署相关的文件，如 Dockerfile 、 docker-compose.yml等。
- `internal`：应用程序的主要代码，按照分层架构进行组织。
- `pkg`：公共的代码，包括配置、日志、HTTP 等。
- `script`：脚本文件，用于部署和其他自动化任务。
- `storage`：存储文件，如日志文件。
- `test`：测试代码。
- `web`：前端代码。

## internal

- `internal/handler`（ or `controller`）：处理 HTTP 请求，调用业务逻辑层的服务，返回 HTTP 响应。
- `internal/server`（or `router`）：HTTP 服务器，启动 HTTP 服务，监听端口，处理 HTTP 请求。
- `internal/service`（or `logic`）：服务，实现具体的业务逻辑，调用数据访问层repository。
- `internal/model`（or `entity`）：数据模型，定义了业务逻辑层需要的数据结构。
- `internal/repository`（or `dao`）：数据访问对象，封装了数据库操作，提供了对数据的增删改查。



## 依赖注入

本项目采用了依赖注入框架`Wire`，实现了模块化和解耦。`Wire`通过预编译`wire.go`，自动生成依赖注入的代码`wire_gen.go`，简化了依赖注入的过程。

- `cmd/job/wire.go`：`Wire`配置文件，定义了`job`子命令需要的依赖关系。
- `cmd/migration/wire.go`：`Wire`配置文件，定义了`migration`子命令需要的依赖关系。
- `cmd/server/wire.go`：`Wire`配置文件，定义了`server`子命令需要的依赖关系。

Wire官方文档：https://github.com/google/wire/blob/main/docs/guide.md

注意：`wire_gen.go`文件为自动编译生成，禁止手动修改

## 公共代码

为了实现代码的复用和统一管理，本项目采用了公共代码的方式，将一些通用的代码放在了`pkg`目录下。

- `pkg/config`：配置文件的读取和解析。
- `pkg/helper`：一些通用的辅助函数，如 MD5 加密、UUID 生成等。
- `pkg/http`：HTTP 相关的代码，如 HTTP 客户端、HTTP 服务器等。
- `pkg/log`：日志相关的代码，如日志的初始化、日志的写入等。
- `more...`：当然，你可以自由添加扩展更多的pkg。