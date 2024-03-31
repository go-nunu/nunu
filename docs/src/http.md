# HTTP
**用途**：处理基于HTTP协议的请求和响应。

**工作原理**：监听指定的HTTP端口，并根据接收到的请求执行相应的操作，然后返回相应的HTTP响应。

**示例应用场景**：Web应用程序、API服务等。

## 路由定义
`HTTP`的路由定义非常简单，大家可以直接参考`internal/server/http.go`中的`NewHTTPServer`方法。

在`NewHTTPServer`方法中，我们首先创建了一个`gin.Engine`对象，然后定义了路由规则，包括`GET`、`POST`、`PUT`、`DELETE`等方法。

需要注意的是在`高级Layout`示例中，我们为大家定义了三个路由组，`noAuthRouter`、`noStrictAuthRouter`和`strictAuthRouter`，他们的具体用法如下：
1. `noAuthRouter`：无需认证即可访问，用于一些无需认证的接口，例如登录、注册等。
2. `noStrictAuthRouter`：无需严格认证即可访问，用于一些无需严格认证的接口，例如获取用户信息等。
3. `strictAuthRouter`：需要严格认证即可访问，用于一些需要严格认证的接口，例如修改用户信息等。

三个路由组是基于不同的中间件实现的，具体中间件的实现可以参考`internal/middleware`目录下的代码。

## 依赖注入Handler
`HTTP`模块的依赖注入非常简单，只需要在`NewHTTPServer`方法中传入的用到的`Handler`结构即可。
```go
func NewHTTPServer(
	logger *log.Logger,
	conf *viper.Viper,
	jwt *jwt.JWT,
	userHandler *handler.UserHandler,
	// 更多handler
) *http.Server {
	
}
```