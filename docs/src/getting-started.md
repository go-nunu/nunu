# 快速开始
## 安装Nunu {#install}

您可以通过以下命令安装 Nunu：

```bash
go install github.com/go-nunu/nunu@latest
```


:::tip
国内用户可以使用`GOPROXY`加速`go install`
:::

```bash
$ go env -w GO111MODULE=on
$ go env -w GOPROXY=https://goproxy.cn,direct
```

:::details 如果`go install`成功，却提示找不到nunu命令？

这是因为环境变量没有配置，可以把 GOBIN 目录配置到环境变量中即可，GOBIN的配置方法windows、macOS、Linux各不相同，请自行网络搜索。
:::


## 创建新项目{new}

::: tip 在创建项目之前，我们需要了解到Nunu内置了两种类型的Layout：


* **基础模板(Basic Layout)**

  `Basic Layout` 包含一个**非常精简**的架构目录结构，适合非常熟悉Nunu项目的开发者使用。

* **高级模板(Advanced Layout)**

  `Advanced Layout` 包含了很多`Nunu`的用法示例（ `db`、`redis`、 `jwt`、 `cron`、 `migration`等），适合开发者快速学习了解Nunu的架构思想。


:::

您可以使用以下命令创建一个新的 `Golang` 项目，推荐选择`Advanced Layout`：

```bash
nunu new projectName
```

此命令将创建一个名为 `projectName` 的目录，并在其中生成一个优雅的 Golang 项目结构。

**国内加速源：**

`nunu new`默认拉取github源，你也可以使用国内加速仓库
::: code-group

```sh [使用高级模板(推荐)]
$ nunu new projectName -r https://gitee.com/go-nunu/nunu-layout-advanced.git
```

```sh [使用基础模板]
$ nunu new projectName -r https://gitee.com/go-nunu/nunu-layout-basic.git
```
:::



## 启动服务{run}

创建好项目之后，我们进入项目执行`nunu run`命令即可启动服务。

```
nunu run cmd/server/main.go
```

随后打开浏览器访问`http://localhost:8000`即可看到欢迎页面。

API文档地址: `http://127.0.0.1:8000/swagger/index.html`

