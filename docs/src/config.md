# 配置




## 指定配置文件启动
Nunu 使用 [Viper](https://github.com/spf13/viper) 库来管理配置文件。

预设了两套配置文件，分别是`local.yml`和`prod.yml`，分别用于本地开发环境与生产环境，默认会加载`config/local.yml`。

你可以使用环境变量或参数来指定配置文件路径

::: code-group
```bash [Linux/MacOS]
APP_CONF=config/prod.yml nunu run
```
```bash [Windows]
set APP_CONF=config\prod.yml && nunu run
```
:::

或者使用传参的方式:`go run ./cmd/server -conf=config/prod.yml`  

## 读取配置

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
:::tip
通过参数进行依赖注入之后，别忘记执行`nunu wire`命令生成依赖文件。
:::