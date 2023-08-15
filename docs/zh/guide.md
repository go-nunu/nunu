## 文档
* [使用指南](https://github.com/go-nunu/nunu/blob/main/docs/zh/guide.md)
* [分层架构](https://github.com/go-nunu/nunu/blob/main/docs/zh/architecture.md)
* [上手教程](https://github.com/go-nunu/nunu/blob/main/docs/zh/tutorial.md)
* [高效编写单元测试](https://github.com/go-nunu/nunu/blob/main/docs/zh/unit_testing.md)


[进入英文版](https://github.com/go-nunu/nunu/blob/main/docs/en/guide.md)

# Nunu 使用指南

Nunu 是一个基于 Golang 的应用脚手架，它可以帮助您快速构建高效、可靠的应用程序。本指南将介绍如何使用 Nunu 创建、开发您的应用程序。

## 安装

您可以通过以下命令安装 Nunu：

```bash
go install github.com/go-nunu/nunu@latest
```

国内用户可以使用`GOPROXY`加速`go install`

```
$ go env -w GO111MODULE=on
$ go env -w GOPROXY=https://goproxy.cn,direct
```

> tips: 如果`go install`成功，却提示找不到nunu命令，这是因为环境变量没有配置，可以把 GOBIN 目录配置到环境变量中即可

## 创建新项目

您可以使用以下命令创建一个新的 Golang 项目：

```bash
nunu new projectName

// 推荐新用户选择Advanced Layout
```

此命令将创建一个名为 `projectName` 的目录，并在其中生成一个优雅的 Golang 项目结构。

**国内加速源：**

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
## 创建组件

您可以使用以下命令为项目创建 handler、service 、 repository和model 等组件：

```bash
nunu create handler user
nunu create service user
nunu create repository user
nunu create model user
```

这些命令将分别创建一个名为 `UserHandler`、`UserService` 、 `UserRepository` 和 `UserModel` 的组件，并将它们放置在正确的目录中。

如果你想在自定义的目录创建相应组件则可以这么做：
```bash
nunu create handler internal/handler/user/center
nunu create service internal/service/user/center
nunu create repository internal/repository/user/center
nunu create model internal/model/user/center
```


你还可以使用以下命令一次性创建 handler、service、repository 和 model 等组件：

```bash
nunu create all user
```

## 启动项目

您可以使用以下命令快速启动项目：

```bash
// 等价于  go run ./cmd/server

nunu run
```

此命令将启动您的 Golang 项目，并支持文件更新热重启。

## 编译 wire.go

您可以使用以下命令快速编译 `wire.go`：

```bash
// 等价于 cd cmd/server  && wire
nunu wire
```

此命令会自动寻找项目中的`wire.go`文件，并快速编译生成所需的依赖项。

## 配置文件

### 指定配置文件启动
Nunu 使用 Viper 库来管理配置文件。

默认会加载`config/local.yml`，你可以使用环境变量或参数来指定配置文件路径

```
// Linux or MacOS
APP_CONF=config/prod.yml nunu run

// Windows

set APP_CONF=config\prod.yml && nunu run

```
或者使用传参的方式:`go run ./cmd/server -conf=config/prod.yml`

### 读取配置项

您可以在 `config` 目录下创建一个名为 `local.yaml` 的文件来存储您的配置信息。例如：


```yaml
data:
  mysql:
    user: root:123456@tcp(127.0.0.1:3380)/user?charset=utf8mb4&parseTime=True&loc=Local
  redis:
    addr: 127.0.0.1:6350
    password: ""
    db: 0
    read_timeout: 0.2s
    write_timeout: 0.2s
```

您可以在代码中使用依赖注入`conf *viper.Viper`来读取配置信息：

```go
package repository

import (
	"context"
	"fmt"
	"github.com/go-nunu/nunu-layout-advanced/pkg/log"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type Repository struct {
	db     *gorm.DB
	rdb    *redis.Client
	logger *log.Logger
}

func NewRepository(db *gorm.DB, rdb *redis.Client, logger *log.Logger) *Repository {
	return &Repository{
		db:     db,
		rdb:    rdb,
		logger: logger,
	}
}

func NewDB(conf *viper.Viper) *gorm.DB {
	db, err := gorm.Open(mysql.Open(conf.GetString("data.mysql.user")), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}
func NewRedis(conf *viper.Viper) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.GetString("data.redis.addr"),
		Password: conf.GetString("data.redis.password"),
		DB:       conf.GetInt("data.redis.db"),
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("redis error: %s", err.Error()))
	}

	return rdb
}


```
tips：通过参数进行依赖注入之后，别忘记执行`nunu wire`命令生成依赖文件。

## 日志

Nunu 使用 Zap 库来管理日志。您可以在 `config` 中配置日志。例如：

```yaml
log:
  log_level: info
  encoding: json           # json or console
  log_file_name: "./storage/logs/server.log"
  max_backups: 30              # 日志文件最多保存多少个备份
  max_age: 7                   #  文件最多保存多少天
  max_size: 1024               #  每个日志文件保存的最大尺寸 单位：M
  compress: true               # 是否压缩
```

您可以在代码中使用以下方式来记录日志：

```go
package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-nunu/nunu-layout-basic/internal/service"
	"github.com/go-nunu/nunu-layout-basic/pkg/helper/resp"
	"go.uber.org/zap"
	"net/http"
)

// ...

func (h *userHandler) GetUserById(ctx *gin.Context) {
	h.logger.Info("GetUserByID", zap.Any("user", user))
	// ...
}

// ...

```

## 数据库

Nunu 使用 GORM 库来管理数据库。您可以在 `config` 目录下配置数据库。例如：

```yaml
data:
  mysql:
    user: root:123456@tcp(127.0.0.1:3380)/user?charset=utf8mb4&parseTime=True&loc=Local
  redis:
    addr: 127.0.0.1:6350
    password: ""
    db: 0
    read_timeout: 0.2s
    write_timeout: 0.2s
```

您可以在代码中使用以下方式来连接数据库：

```go
package repository

import (
	"github.com/go-nunu/nunu-layout-advanced/internal/model"
)


type UserRepository interface {
	FirstById(id int64) (*model.User, error)
}
type userRepository struct {
	*Repository
}

func NewUserRepository(repository *Repository) *UserRepository {
	return &UserRepository{
		Repository: repository,
	}
}

func (r *userRepository) FirstById(id int64) (*model.User, error) {
	var user model.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

```

需要注意的是Nunu中的`xxxRepository`、`xxxService`、`xxxHandler`都是基于`interface`实现，

这就是所谓的**面向接口编程**，它可以提高代码的灵活性、可扩展性、可测试性和可维护性，是Go语言非常推崇的一种编程风格。



比如上面的代码我们写成了
```
type UserRepository interface {
	FirstById(id int64) (*model.User, error)
}
type userRepository struct {
	*Repository
}

```
而不是直接写成
```
type UserRepository struct {
	*Repository
}
```
> tips: Nunu高级Layout中的单元测试就是基于`interface`特性进行mock操作的。


## 测试

Nunu 使用 testify、redismock、gomock、go-sqlmock等 库来编写测试。

具体的测试用例可以查看[Nunu advanced layout](https://github.com/go-nunu/nunu-layout-advanced/tree/main/test/server)

您可以使用以下命令运行测试：

```bash
go test -coverpkg=./internal/handler,./internal/service,./internal/repository -coverprofile=./.nunu/coverage.out ./test/server/...
go tool cover -html=./.nunu/coverage.out -o coverage.html

```

上面的命令将会生成一个html文件`coverage.html`，我们可以直接使用浏览器打开它，然后我们就会看到详细的单元测试覆盖率。

## 结论

Nunu 是一个非常实用的 Golang 应用脚手架，它可以帮助您快速构建高效、可靠的应用程序。希望本指南能够帮助您更好地使用 Nunu。
