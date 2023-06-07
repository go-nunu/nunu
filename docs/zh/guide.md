# Nunu 使用指南

Nunu 是一个基于 Golang 的应用脚手架，它可以帮助您快速构建高效、可靠的应用程序。本指南将介绍如何使用 Nunu 创建、开发您的应用程序。

## 安装

您可以通过以下命令安装 Nunu：

```bash
go install github.com/go-nunu/nunu@latest
```

## 创建新项目

您可以使用以下命令创建一个新的 Golang 项目：

```bash
nunu new projectName
```

此命令将创建一个名为 `projectName` 的目录，并在其中生成一个优雅的 Golang 项目结构。

**国内加速源：**

`nunu new`默认拉取github源，你也可以使用国内加速仓库
```
// 使用基础模板
nunu new projectName -r https://gitee.com/go-nunu/nunu-layout-basic.git
// 使用高级模板(推荐)
nunu new projectName -r https://gitee.com/go-nunu/nunu-layout-advanced.git
```


> Nunu内置了两种类型的Layout：

* **基础模板(Basic Layout)**

Basic Layout 包含一个非常精简的架构目录结构，适合非常熟悉Nunu项目的开发者使用。

* **高级模板(Advanced Layout)**

**建议：我们推荐新手优先选择使用Advanced Layout。**


Advanced Layout 包含了很多Nunu的用法示例（ db、redis、 jwt、 cron、 migration等），适合开发者快速学习了解Nunu的架构思想。
## 创建组件

您可以使用以下命令为项目创建 handler、service 和 dao 等组件：

```bash
nunu create handler user
nunu create service user
nunu create dao user
```

这些命令将分别创建一个名为 `UserHandler`、`UserService` 和 `UserDao` 的组件，并将它们放置在正确的目录中。

如果你想在自定义的目录创建相应组件则可以这么做：
```bash
nunu create handler internale/handler/user/center
nunu create service internale/service/user/center
nunu create dao internale/dao/user/center
```


你还可以使用以下命令一次性创建 handler、service 和 dao 等组件：

```bash
nunu create hsd user
```

## 启动项目

您可以使用以下命令快速启动项目：

```bash
// 等价于  go run cmd/server/main.go 

nunu run
```

此命令将启动您的 Golang 项目，并支持文件更新热重启。

## 编译 wire.go

您可以使用以下命令快速编译 `wire.go`：

```bash
// 等价于 cd cmd/server/wire  && wire
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
或者使用传参的方式:`go run cmd/server/main.go -conf=config/prod.yml`

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
package dao

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

type Dao struct {
	db     *gorm.DB
	rdb    *redis.Client
	logger *log.Logger
}

func NewDao(db *gorm.DB, rdb *redis.Client, logger *log.Logger) *Dao {
	return &Dao{
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
import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

func main() {
    logger, err := zap.Config{
        Encoding:         viper.GetString("logger.encoding"),
        Level:            zap.NewAtomicLevelAt(zapcore.Level(viper.GetInt("logger.level"))),
        OutputPaths:      []string{viper.GetString("logger.outputPaths.0"), viper.GetString("logger.outputPaths.1")},
        ErrorOutputPaths: []string{viper.GetString("logger.errorOutputPaths.0"), viper.GetString("logger.errorOutputPaths.1")},
        InitialFields: map[string]interface{}{
            "app":     viper.GetString("logger.initialFields.app"),
            "version": viper.GetString("logger.initialFields.version"),
        },
    }.Build()
    if err != nil {
        panic(err)
    }

    logger.Info("Hello, Nunu!")
}
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
package dao

import (
	"github.com/go-nunu/nunu-layout-advanced/internal/model"
)

type UserDao struct {
	*Dao
}

func NewUserDao(dao *Dao) *UserDao {
	return &UserDao{
		Dao: dao,
	}
}

func (r *UserDao) FirstById(id int64) (*model.User, error) {
	var user model.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

```


## 测试

Nunu 使用 Testify 库来编写测试。您可以在 `test` 目录下创建一个名为 `main_test.go` 的文件来编写测试。例如：

```go
import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/go-nunu/nunu/pkg/handler"
    "github.com/stretchr/testify/assert"
)

func TestUserHandler_GetUser(t *testing.T) {
    r := httptest.NewRequest(http.MethodGet, "/users/1", nil)
    w := httptest.NewRecorder()

    h := handler.NewUserHandler()
    h.GetUser(w, r)

    assert.Equal(t, http.StatusOK, w.Code)
    assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
    assert.JSONEq(t, `{"id":1,"name":"Alice"}`, w.Body.String())
}
```

您可以使用以下命令运行测试：

```bash
go test ./test/...
```

## 结论

Nunu 是一个非常实用的 Golang 应用脚手架，它可以帮助您快速构建高效、可靠的应用程序。希望本指南能够帮助您更好地使用 Nunu。