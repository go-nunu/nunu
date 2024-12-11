# HTTP中间件

在一个 Web 应用程序中，Middleware（中间件）用于在处理请求和生成响应的过程中执行一些通用的逻辑。通常，这些逻辑包括对请求进行预处理、验证、日志记录、错误处理等操作。Middleware 是一种类似于过滤器或拦截器的机制，可以在请求到达处理器之前、处理器执行过程中、处理器执行完成之后等不同阶段插入自定义的逻辑。

Middleware 通常被组织到一个名为 middleware 的目录中。这个目录用于存放各种中间件的代码文件，以便于组织、管理和维护。

## 核心用途

以下是一些常见的 Middleware 的应用场景：

* **身份验证和授权**：检查用户是否已登录，以及用户是否有权限访问某些资源。如果用户未经身份验证或者没有足够的权限，可以重定向到登录页面或返回相应的错误信息。

* **日志记录**：记录请求和响应的相关信息，如请求方法、请求路径、响应状态码、处理时间等，以便后续的监控和分析。

* **错误处理**：捕获处理器执行过程中发生的错误，并返回适当的错误响应给客户端。这样可以避免将错误暴露给客户端，并提供更友好的错误信息。

* **请求预处理**：对请求进行预处理，例如解析请求体、验证请求参数、设置请求上下文等。这些操作可以减少处理器中的重复代码，并提高代码的可复用性和可维护性。

* **性能监控**：监控请求处理的性能指标，如处理时间、内存占用、数据库查询次数等。这些指标可以帮助开发人员识别潜在的性能瓶颈，并进行优化。

## 预置中间件

### 限流器

### 日志Trace

我们通过`random.UUIdV4()`生成一个 UUID 并将其转换为 MD5 字符串作为请求的唯一标识，并记录该请求的 HTTP 方法、请求头、请求 URL 与请求参数

在对于响应的记录中我们采用了自定义响应写入类，并重写了`Write()`方法实现了对响应的记录，允许我们在写入响应时同时捕获响应体的内容。

```go
type bodyLogWriter struct {
 gin.ResponseWriter
 // body 指向 bytes.Buffer 的指针，勇于存储响应题的内容
 body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
 w.body.Write(b)
 return w.ResponseWriter.Write(b)
}
```
我们通过`bodyLogWriter.Write(b []byte)`方法调用`w.ResponseWriter.Write(b)`将响应内容写入客户端从而能够使响应内容被 logger 进行记录

### 签名验证

在签名校验中间件中，我们默认采用 MD5 的加密方式，需要提供四个必要的请求头分别为`Timestamp`、`Nonce`、`Sign`与`App-Version`

#### 校验参数 

| 请求头 | Timestamp | Nonce | Sign | App-Version |
|---|---|---|---|---|
| 含义 | 时间戳 | 请求唯一标识 | 请求签名 | App 版本号 |

#### 校验步骤

1. **获取请求头**  
首先我们回遍历请求头，检查每个必要的请求头都存在且不为空。如果缺少任意一个必要的请求头都将返回 400 错误并终止请求处理

2. **构建签名数据**

    ```go
    data := map[string]string{
        "AppKey":     conf.GetString("security.api_sign.app_key"),
        "Timestamp":  ctx.Request.Header.Get("Timestamp"),
        "Nonce":      ctx.Request.Header.Get("Nonce"),
        "AppVersion": ctx.Request.Header.Get("App-Version"),
    }
    ```

    创建一个包含应用密钥、时间戳、请求唯一标识与 App 版本号的 map 用于于存储签名数据

3. **排序与构建签名字符串**  
我们按照键进行陪许，并构建用于生成签名的字符串，并附加加密密钥

4. **签名校验**
我们通过计算生成的字符串的 MD5 值是否与请求中的 Sign 字段是否匹配。如果不匹配则返回 400 错误并终止请求处理，如果匹配则放行。

### 跨域
