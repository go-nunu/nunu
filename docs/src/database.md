# 数据库

## 支持多数据库{#support-multi-database}

由于`wire`不支持一个方法注入多个相同类型的依赖，所以如果我们想链接多个数据库，可以自己定义对应的新类型来进行注入。
```go
type OrderDB gorm.DB
type UserDB gorm.DB
```

具体实现代码如下:

::: code-group
```go [repository/repository.go]
type UserDB gorm.DB
type OrderDB gorm.DB
type Repository struct {
	userDB  *gorm.DB
	orderDB *gorm.DB
	logger  *log.Logger
}

func NewRepository(userDB *UserDB, orderDB *OrderDB, rdb *redis.Client, logger *log.Logger) *Repository {
	return &Repository{
		userDB:  (*gorm.DB)(userDB),
		orderDB: (*gorm.DB)(orderDB),
		rdb:     rdb,
		logger:  logger,
	}
}

func NewUserDB(conf *viper.Viper, l *log.Logger) *UserDB {
	db, err := gorm.Open(mysql.Open(conf.GetString("data.db.user")),)
	if err != nil {
		panic(err)
	}
	return (*UserDB)(db)
}
func NewOrderDB(conf *viper.Viper, l *log.Logger) *OrderDB {
	db, err := gorm.Open(mysql.Open(conf.GetString("data.db.order")),)
	if err != nil {
		panic(err)
	}
	return (*OrderDB)(db)
}
// ...

```

```go [repository/user.go]
func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	if err := r.userDB.Create(user).Error; err != nil {
		return err
	}
	return nil
}
```

:::
## 数据库事务{#transaction}

### 事务实现原理{#principle}
Nunu的事务基于`context.Context`实现，核心代码如下：

::: code-group
```go [repository/repository.go]
//这个接口定义了一个方法 Transaction(ctx context.Context, fn func(ctx context.Context) error) error。
//该方法接受一个上下文对象 ctx 和一个函数 fn，该函数用于执行事务内的操作。
//方法返回一个错误，可能是事务执行过程中出现的任何错误。
type Transaction interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

//这个函数是构造函数，用于创建一个新的事务对象。
//接受一个名为 r 的仓库对象指针参数。
//返回一个实现了 Transaction 接口的对象，即仓库对象。
func NewTransaction(r *Repository) Transaction {
	return r
}

//这个方法接受一个上下文对象 ctx。
//它首先尝试从上下文中获取事务对象，如果存在，则返回该事务对象的数据库连接。
//如果不存在事务对象，则创建一个新的数据库连接，传入上下文对象，并返回该连接。
func (r *Repository) DB(ctx context.Context) *gorm.DB {
	v := ctx.Value(ctxTxKey)
	if v != nil {
		if tx, ok := v.(*gorm.DB); ok {
			return tx
		}
	}
	return r.db.WithContext(ctx)
}


//它接受一个上下文对象 ctx 和一个函数 fn，该函数用于执行事务内的操作。
//在执行事务之前，它通过调用 WithContext 方法创建一个带有事务对象的新上下文，并将该上下文传递给事务函数 fn。
//执行事务函数 fn，如果函数执行成功，则提交事务；如果函数执行失败，则回滚事务。
func (r *Repository) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, ctxTxKey, tx)
		return fn(ctx)
	})
}
```
:::

这段代码利用了 `Go` 语言中的`上下文（context）`来传递事务对象，以确保在同一事务内的所有操作都共享同一个数据库连接。这种方式可以确保在一个事务内的所有数据库操作都可以进行回滚，以保持数据的一致性。

### 如何使用事务{#how-to-use}
在使用事务时，需要先创建一个事务对象，然后通过调用事务对象的 `DB` 方法获取数据库连接，并执行数据库操作。如果执行过程中出现错误，可以调用事务对象的 `Rollback` 方法回滚事务，否则调用 `Commit` 方法提交事务。
::: code-group
```go [service/user.go]
func (s *userService) Register(ctx context.Context, req *v1.RegisterRequest) error {
	// ...
    // Transaction demo
    err = s.tm.Transaction(ctx, func (ctx context.Context) error {
        // Create a user
        if err = s.userRepo.Create(ctx, user); err != nil {
            return err
        }
        // Create a user wallet
        if err = s.userWalletRepo.Create(ctx, user); err != nil {
            return err
        }
        // TODO: other repo
        return nil
    })
	// ...
}
```
:::