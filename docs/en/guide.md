## Documentation
* [Guide](https://github.com/go-nunu/nunu/blob/main/docs/en/guide.md)
* [Architecture](https://github.com/go-nunu/nunu/blob/main/docs/en/architecture.md)
* [Tutorial](https://github.com/go-nunu/nunu/blob/main/docs/en/tutorial.md)

[进入简体中文版](https://github.com/go-nunu/nunu/blob/main/docs/zh/guide.md)

# Nunu User Guide

Nunu is a Golang-based application scaffold that helps you quickly build efficient and reliable applications. This guide will show you how to use Nunu to create and develop your applications.

## Installation

You can install Nunu using the following command:

```bash
go install github.com/go-nunu/nunu@latest
```

## Creating a New Project

You can use the following command to create a new Golang project:

```bash
nunu new projectName
```

This command will create a directory named `projectName` and generate an elegant Golang project structure within it.

> Nunu comes with two types of layouts:

* **Basic Layout**

The Basic Layout contains a very minimal architecture directory structure, suitable for developers who are very familiar with Nunu projects.

* **Advanced Layout**

**Recommendation: We recommend that beginners choose the Advanced Layout first.**

The Advanced Layout contains many examples of Nunu's usage (db, redis, jwt, cron, migration, etc.), which is suitable for developers to quickly learn and understand Nunu's architectural ideas.

## Creating Components

You can use the following commands to create handler, service, and dao components for your project:

```bash
nunu create handler user
nunu create service user
nunu create dao user
```

These commands will create components named `UserHandler`, `UserService`, and `UserDao`, respectively, and place them in the correct directory.

If you want to create the corresponding components in a custom directory, you can do so like this:

```bash
nunu create handler internale/handler/user/center
nunu create service internale/service/user/center
nunu create dao internale/dao/user/center
nunu create model internale/model/user/center
```

You can also use the following command to create handler, service, dao and model components at once:

```bash
nunu create all user
```

## Starting the Project

You can use the following command to quickly start your project:

```bash
// equivalent to go run cmd/server/main.go

nunu run
```

This command will start your Golang project and support file update hot restart.

## Compiling wire.go

You can use the following command to quickly compile `wire.go`:

```bash
// equivalent to cd cmd/server/wire && wire
nunu wire
```

This command will automatically find the `wire.go` file in your project and quickly compile the required dependencies.

## Configuration File

### Starting with a Specified Configuration File

Nunu uses the Viper library to manage configuration files.

By default, `config/local.yml` will be loaded, and you can use environment variables or parameters to specify the configuration file path.

```
// Linux or MacOS
APP_CONF=config/prod.yml nunu run

// Windows

set APP_CONF=config\prod.yml && nunu run

```

Or use the parameter method: `go run cmd/server/main.go -conf=config/prod.yml`

### Reading Configuration Items

You can create a file named `local.yaml` in the `config` directory to store your configuration information. For example:

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

You can use dependency injection `conf *viper.Viper` to read the configuration information in your code:

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
Tips: After dependency injection through parameters, don't forget to execute the `nunu wire` command to generate the dependency file.

## Logging

Nunu uses the Zap library to manage logging. You can configure logging in the `config` directory. For example:

```yaml
log:
  log_level: info
  encoding: json           # json or console
  log_file_name: "./storage/logs/server.log"
  max_backups: 30              # maximum number of backup log files to keep
  max_age: 7                   # maximum number of days to keep a log file
  max_size: 1024               # maximum size in megabytes of each log file
  compress: true               # whether to compress log files
```

You can use the following method in your code to record logs:

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

## Database

Nunu uses the GORM library to manage databases. You can configure the database in the `config` directory. For example:

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

You can use the following method in your code to connect to the database:

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


## Testing

Nunu uses the Testify library to write tests. You can create a file named `main_test.go` in the `test` directory to write tests. For example:

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

You can use the following command to run tests:

```bash
go test ./test/...
```

## Conclusion

Nunu is a very practical Golang application scaffold that helps you quickly build efficient and reliable applications. We hope this guide will help you better use N