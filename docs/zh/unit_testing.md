## 文档
* [使用指南](https://github.com/go-nunu/nunu/blob/main/docs/zh/guide.md)
* [分层架构](https://github.com/go-nunu/nunu/blob/main/docs/zh/architecture.md)
* [上手教程](https://github.com/go-nunu/nunu/blob/main/docs/zh/tutorial.md)
* [高效编写单元测试](https://github.com/go-nunu/nunu/blob/main/docs/zh/unit_testing.md)


# 单元测试

## 介绍

在项目中进行单元测试是一种重要的开发实践。然而，当被测代码依赖其他模块或组件时，编写单元测试变得复杂且不稳定。为了解决这个问题，我们可以使用mock来模拟被测代码的依赖。通过使用mock对象，我们可以控制外部模块的行为，使得被测代码在测试过程中不会真正依赖和调用外部模块，从而实现对被测代码的隔离。在Go语言中，使用golang/mock库来生成mock代码，并使用sqlmock和redismock来模拟数据库和缓存的行为。通过使用mock，我们可以提高单元测试的可靠性和效率。本文将介绍如何使用mock来编写简洁高效的单元测试。

## 导读

首先我们先来看下项目中的依赖注入文件`cmd/server/wire.go`：

> tip: 该文件由`google/wire`工具自动编译生成，禁止人为编辑

```
// Injectors from wire.go:

func newApp(viperViper *viper.Viper, logger *log.Logger) (*gin.Engine, func(), error) {
	jwt := middleware.NewJwt(viperViper)
	handlerHandler := handler.NewHandler(logger)
	sidSid := sid.NewSid()
	serviceService := service.NewService(logger, sidSid, jwt)
	db := repository.NewDB(viperViper)
	client := repository.NewRedis(viperViper)
	repositoryRepository := repository.NewRepository(db, client, logger)
	userRepository := repository.NewUserRepository(repositoryRepository)
	userService := service.NewUserService(serviceService, userRepository)
	userHandler := handler.NewUserHandler(handlerHandler, userService)
	engine := server.NewServerHTTP(logger, jwt, userHandler)
	return engine, func() {
	}, nil
}
```

从这段代码我们可以得知`repository`、`service`、`repository`之间的依赖关系，

`userHandler`依赖于`userService`，而`userService`又依赖于`userRepository。

比如`handler/user.go`下面的`GetProfile`代码如下：
```
func (h *userHandler) GetProfile(ctx *gin.Context) {
	userId := GetUserIdFromCtx(ctx)
	if userId == "" {
		resp.HandleError(ctx, http.StatusUnauthorized, 1, "unauthorized", nil)
		return
	}

	user, err := h.userService.GetProfile(ctx, userId)
	if err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}

	resp.HandleSuccess(ctx, user)
}
```
我们会发现在它的内部调用了`userService.GetProfile`。

因此在编写单元测试的时候，我们就不可避免的需要先初始化`userService`实例，而当我们去初始化`userService`的时候，我们又会发现它又依赖于`userRepository`。

明明我们只需要测试一个最底层的`handler`，却需要先初始化执行`service`、`repository`等代码。 这很明显违背了单元测试的（单一职责原则），每个单元测试只关注一个功能点或一个代码单元。

有什么比较好的办法解决该问题呢，我们的最终答案就是`mock`。


### Mock（依赖隔离好帮手）

在进行单元测试时，我们希望测试的是被测代码单元的逻辑，而不希望依赖其他外部模块或组件的状态或行为。这样做可以更好地隔离被测代码，使得测试更加可靠和可重复。

Mock是一种测试模式，用于模拟或替代被测代码所依赖的外部模块或组件。通过使用Mock对象，我们可以控制外部模块的行为，使得被测代码在测试过程中不会真正依赖和调用外部模块，从而实现对被测代码的隔离。

Mock对象可以模拟外部模块的返回值、异常、超时等，使得测试可以更加可控和可预测。它解决了以下问题：

1. 依赖其他模块：某些代码单元可能依赖其他模块，例如数据库、网络请求等。通过使用Mock对象，我们可以模拟这些依赖，使得测试不需要真正依赖这些模块，从而避免测试的不稳定性和复杂性。

2. 隔离外部环境：某些代码单元可能受到外部环境的影响，例如当前时间、系统状态等。通过使用Mock对象，我们可以控制这些外部环境的状态，使得测试可以在不同环境下运行，从而增加测试的覆盖范围和准确性。

3. 提高测试效率：某些外部模块可能执行耗时操作，例如网络请求、文件读写等。通过使用Mock对象，我们可以避免真实执行这些操作，从而提高测试的执行速度和效率。


在nunu项目中，我们采用以下mock库来帮助我们编写单元测试

* github.com/golang/mock            // google开源的mock库
* github.com/go-redis/redismock/v9  // 提供redis查询的模拟测试，兼容github.com/redis/go-redis/v9
* github.com/DATA-DOG/go-sqlmock    // sqlmock是一个实现sql/driver 的模拟库

## 面向接口编程

使用`golang/mock`有个前提，我们需要遵循"面向接口编程"的方式来编写我们的`repository`和`service`。

可能有的同学不了解"面向接口编程"是什么意思，我们这儿以一段代码举例：

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

func (r *UserRepository) FirstById(id int64) (*model.User, error) {
	var user model.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

```

上面的代码中，我们先定义一个`UserRepository interface`,然后通过`userRepository struct`去实现它的所有方法。
```
type UserRepository interface {
	FirstById(id int64) (*model.User, error)
}
type userRepository struct {
	*Repository
}
func (r *UserRepository) FirstById(id int64) (*model.User, error) {
    // ...
}

```
而不是直接写成
```
type UserRepository struct {
	*Repository
}

func (r *UserRepository) FirstById(id int64) (*model.User, error) {
    // ...
}
```

这就是所谓的**面向接口编程**，它可以提高代码的灵活性、可扩展性、可测试性和可维护性，是Go语言非常推崇的一种编程风格。


## go-mock快速上手

`golang/mock`的使用其实简单，我们首先安装一下它：

```
go install github.com/golang/mock/mockgen@v1.6.0
```

`mockgen`是`go-mock`的一个命令行工具，可以解析我们代码中的`interface`定义，自动生成正确的mock代码


示例：
```
mockgen -source=internal/service/user.go -destination mocks/service/user.go
```

上面的命令指定了两个参数，interface源文件以及最终生成mock代码的目标文件，我们将目标文件放置在`mocks/service`目录下面。

生成了`UserService`的`mock`代码，我们就可以去编写`UserHandler`的单元测试了。

最终的单测代码如下：

```

func TestUserHandler_GetProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_service.NewMockUserService(ctrl)
	
	// 关键代码，定义mockUserService.GetProfile的返回值
	mockUserService.EXPECT().GetProfile(gomock.Any(), userId).Return(&model.User{
		Id:       1,
		UserId:   userId,
		Username: "xxxxx",
		Nickname: "xxxxx",
		Password: "xxxxx",
		Email:    "xxxxx@gmail.com",
	}, nil)

	router := setupRouter(mockUserService)
	req, _ := http.NewRequest("GET", "/user", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, resp.Code, http.StatusOK)
	// Add assertions for the response body if needed
}

```


完整的源码位于： https://github.com/go-nunu/nunu-layout-advanced/blob/main/test/server/handler/user_test.go

## sqlmock与redismock

相对于`handler`和`service`的单元测试，`repository`的稍微有些不一样，因为它依赖的不再是我们自己的业务模块，而是依赖于rpc、redis、MySQL这些外部数据源。

这种情况下，为了避免连接真实的数据库和缓存，减少测试的不确定性，我们同样进行mock。

代码如下
```
package repository

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-nunu/nunu-layout-advanced/internal/model"
	"github.com/go-nunu/nunu-layout-advanced/internal/repository"
	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupRepository(t *testing.T) (repository.UserRepository, sqlmock.Sqlmock) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      mockDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm connection: %v", err)
	}

	rdb, _ := redismock.NewClientMock()

	repo := repository.NewRepository(db, rdb, nil)
	userRepo := repository.NewUserRepository(repo)

	return userRepo, mock
}


func TestUserRepository_GetByUsername(t *testing.T) {
	userRepo, mock := setupRepository(t)

	ctx := context.Background()
	username := "test"

    // 模拟查询测试数据
	rows := sqlmock.NewRows([]string{"id", "user_id", "username", "nickname", "password", "email", "created_at", "updated_at"}).
		AddRow(1, "123", "test", "Test", "password", "test@example.com", time.Now(), time.Now())
	mock.ExpectQuery("SELECT \\* FROM `users`").WillReturnRows(rows)

	user, err := userRepo.GetByUsername(ctx, username)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "test", user.Username)

	assert.NoError(t, mock.ExpectationsWereMet())
}

```

完整代码位于：https://github.com/go-nunu/nunu-layout-advanced/blob/main/test/server/repository/user_test.go


## 测试覆盖率
Golang官方原生支持生成测试覆盖率报告。

```
go test -coverpkg=./internal/handler,./internal/service,./internal/repository -coverprofile=./coverage.out ./test/server/...

go tool cover -html=./coverage.out -o coverage.html
```

上面的2条命令将会生成一个网页可视化的覆盖率报告文件`coverage.html`，我们可以直接使用浏览器打开它。

效果如下：

![coverage](https://github.com/go-nunu/nunu/blob/main/.github/assets/coverage.png)

## 总结


单元测试在项目中是一种重要的开发实践，可以确保代码的正确性并提供自动化验证功能。在进行单元测试时，我们需要面向接口编程，使用mock对象来隔离被测代码的依赖关系。在Go语言中，我们可以使用golang/mock库来生成mock代码。对于依赖外部数据源的repository，我们可以使用sqlmock和redismock来模拟数据库和缓存的行为。通过使用mock对象，我们可以控制外部模块的行为，使得被测代码在测试过程中不会真正依赖和调用外部模块，从而实现对被测代码的隔离。这样可以提高测试的可靠性、可重复性和效率。