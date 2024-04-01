# Redis

Redis的使用非常简单，核心代码如下：
::: code-group


```go [repository/repository.go]
type Repository struct {
    db *gorm.DB
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
```go [wire.go]
var repositorySet = wire.NewSet(
	repository.NewRedis,
}
```

```go [repository/user.go]
func NewUserRepository(repository *Repository) CaptchaRepository {
	return &useraRepository{
		Repository: repository,
	}
}
const KeyCaptcha = "Token:%d"

func (r *captchaRepository) SetToken(ctx context.Context, account, token string) error {
	return r.rdb.Set(ctx, fmt.Sprintf(account,  token), code, 15*time.Minute).Err()
}
```

:::

我们首先在`repository/repository.go`定义redis的初始化函数，然后通过`wire.go`声明依赖注入，最后在`repository/user.go`中实现具体的操作即可。


::: tip
每次修改完`wire.go`，都需要执行`nunu wire`命令，重新编译wire哦。
:::