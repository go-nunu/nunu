# nunu命令行工具 {#info}
nunu 命令行工具是一个为了协助快速开发 Golang 项目而创建的项目，通过 nunu 命令行工具， 您可以很容易的进行 Golang 项目的创建、热编译、开发、测试、和部署。

## 创建项目 {#new}
```
nunu new [projectName] 
```
## 启动项目 {#run}
```
nunu run
```
:::tip
通常情况下，nunu run 命令仅用于本地开发环境快速热编译运行使用。如果是生产环境，请使用 `go build`之后部署。
::: 

## 创建组件 {#create}
```
nunu create [type] [handler-name]
```
详细命令：

::: code-group
``` bash [创建Handler]
nunu create handler [handler-name]
```

``` bash [创建Service]
nunu create service [service-name]
```

``` bash [创建Repository]
nunu create  repository [repository-name]
```

``` bash [创建Model]
nunu create  model [model-name]
```
:::
如果你觉得每种组件单独创建太麻烦，你可以使用 `nunu create all` 创建所有组件。
```
nunu create all [name]
```

## 编译wire {#waire}
```
nunu wire all
```
:::tip
如果你的项目存在多个`wire.go`文件，而你只想编译指定的`wire.go`文件，你可以使用 `nunu wire`，然后自己选择对应的文件编译。
:::
## 版本升级 {#upgrade}
```
nunu upgrade
```