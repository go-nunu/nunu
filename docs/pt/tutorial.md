## Documentação
* [Guia do Usuário](https://github.com/go-nunu/nunu/blob/main/docs/pt/guide.md)
* [Arquitetura](https://github.com/go-nunu/nunu/blob/main/docs/pt/architecture.md)
* [Tutorial de Início Rápido](https://github.com/go-nunu/nunu/blob/main/docs/pt/tutorial.md)
* [Teste de Unidade](https://github.com/go-nunu/nunu/blob/main/docs/pt/unit_testing.md)

- [Chinês](https://github.com/go-nunu/nunu/blob/main/docs/zh/tutorial.md)
- [Português](https://github.com/go-nunu/nunu/blob/main/docs/pt/tutorial.md)

# Guia do Usuário do Framework Nunu

Nunu é um framework web baseado na linguagem de programação Go. Ele oferece uma estrutura de projeto elegante e operações de comando que permitem aos desenvolvedores desenvolver aplicações web de forma eficiente.

## Requisitos
Para usar o Nunu com o Layout Avançado, você precisa ter o seguinte software instalado no seu sistema:

* Golang 1.19 ou superior
* Git
* MySQL 5.7 ou superior
* Redis

## Instalação

Antes de começar a usar o Nunu, você precisa instalá-lo. Você pode fazer isso executando o seguinte comando:

```bash
go install github.com/go-nunu/nunu@latest
```

Para usuários na China, você pode usar `GOPROXY` para acelerar o `go install`.
```bash
$ go env -w GO111MODULE=on
$ go env -w GOPROXY=https://goproxy.cn,direct
```

> Dicas: Se `go install` for bem-sucedido, mas você receber um erro dizendo "comando nunu não encontrado", significa que a variável de ambiente não está configurada. Você pode adicionar o diretório GOBIN à variável de ambiente.

## Criando um Projeto

Criar um novo projeto com o Nunu é muito simples. Basta executar o seguinte comando na linha de comando:

```bash
nunu new nomeDoProjeto
```

Substitua `nomeDoProjeto` pelo nome do seu projeto. Aqui, escolheremos o Layout Avançado.

**Usando um Repositório Acelerado na China:**

Por padrão, `nunu new` puxa do repositório do GitHub, mas você também pode usar um repositório acelerado na China.

```bash
# Usando o template avançado (recomendado)
nunu new nomeDoProjeto -r https://gitee.com/go-nunu/nunu-layout-advanced.git

# Usando o template básico
nunu new nomeDoProjeto -r https://gitee.com/go-nunu/nunu-layout-basic.git
```

Após executar o comando acima, o Nunu criará automaticamente um projeto Go bem estruturado com alguns arquivos e diretórios comumente usados.


## Criando Componentes

No Nunu, você pode usar o seguinte comando para criar componentes Handler, Service, Repository e Model em lotes:

```bash
nunu create all order
```

Aqui, `order` é o nome do componente que você deseja criar.

Após executar o comando acima, o Nunu criará automaticamente os componentes nos diretórios correspondentes e escreverá as estruturas correspondentes e alguns métodos comumente usados.
```bash
# Informação de log
Created new handler: internal/handler/order.go
Created new service: internal/service/order.go
Created new repository: internal/repository/order.go
Created new model: internal/model/order.go
```

## Registrando Rotas
Edite  `internal/server/http.go`.

Adicione `handler.OrderHandler` como um parâmetro para `NewServerHTTP`, que configura a dependência para `OrderHandler`.

Em seguida, registre uma nova rota: `noAuthRouter.GET("/order", orderHandler.GetOrderById)`.
```go
func NewServerHTTP(
	// ...
	orderHandler *handler.OrderHandler,     // novo
) *gin.Engine {
    // ...

	// No authentication routes
	noAuthRouter := r.Group("/").Use(middleware.RequestLogMiddleware(logger))
	{
		noAuthRouter.GET("/order", orderHandler.GetOrderById)   // novo
```

## Escrevendo Provedores Wire
Edite `cmd/server/wire.go` e adicione as funções de fábrica geradas a partir dos arquivos ao `providerSet`, conforme mostrado abaixo:
```go
//go:build wireinject
// +build wireinject

package main

// ...

var HandlerSet = wire.NewSet(
	handler.NewHandler,
	handler.NewUserHandler,

	handler.NewOrderHandler, // novo
)

var ServiceSet = wire.NewSet(
	service.NewService,
	service.NewUserService,

	service.NewOrderService, // novo
)

var RepositorySet = wire.NewSet(
	repository.NewDB,
	repository.NewRedis,
	repository.NewRepository,
	repository.NewUserRepository,

	repository.NewOrderRepository, // novo
)

func newApp(*viper.Viper, *log.Logger) (*gin.Engine, func(), error) {
	panic(wire.Build(
		ServerSet,
		RepositorySet,
		ServiceSet,
		HandlerSet,
		SidSet,
		JwtSet,
	))
}

```
## Compilando o Wire

No Nunu, você pode usar o seguinte comando para compilar o Wire:

```bash
nunu wire all
```

Após executar o comando acima, você verá que o arquivo `wire_gen.go` é gerado a partir do arquivo `cmd/server/wire.go`.

Abra o arquivo `wire_gen.go`, e você verá que o código de dependência para `orderRepository`, `orderService` e `orderHandler` foi gerado automaticamente.

```go
func NewApp(viperViper *viper.Viper, logger *log.Logger) (*gin.Engine, func(), error) {
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
	
	
	orderRepository := repository.NewOrderRepository(repositoryRepository)
	orderService := service.NewOrderService(serviceService, orderRepository)
	orderHandler := handler.NewOrderHandler(handlerHandler, orderService)
	
	
	engine := server.NewServerHTTP(logger, jwt, userHandler, orderHandler)
	return engine, func() {
	}, nil
}

```

Neste ponto, completamos o processo central do projeto Nunu.

Em seguida, você precisa modificar as informações de configuração do MySQL e Redis em `config/local.yml` e escrever seu código lógico nos arquivos relevantes.
```bash
internal/handler/order.go            // Manipular parâmetros de solicitação e respostas
internal/service/order.go            // Implementar lógica de negócios
internal/repository/order.go         // Interagir com bancos de dados e Redis
internal/model/order.go              // Entidade da tabela do banco de dados, modelo GORM
```

## Iniciando o Projeto
Finalmente, no Nunu, você pode usar o seguinte comando para iniciar o projeto:

```bash
# Por favor, modifique as informações de configuração do MySQL e Redis em config/local.yml antes de iniciar o servidor
# Antes de iniciar o servidor pela primeira vez, execute a seguinte migração de banco de dados
nunu run ./cmd/migration  

# Iniciar o servidor
nunu run ./cmd/server    

# Ou
nunu run

# Ou
nunu run ./cmd/server  --excludeDir=".git,.idea,tmp,vendor" --includeExt="go,yml,vue"  -- --conf=./config/local.yml
```

Após executar o comando acima, o Nunu iniciará automaticamente o projeto e monitorará atualizações de arquivos, suportando recarregamento em tempo real.



## Geração Automática de Documentação Swagger

Primeiro, precisamos instalar a ferramenta de linha de comando swag em nossa máquina local. Você pode fazer isso executando o seguinte comando:
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

[swaggo](https://github.com/swaggo/swag) nos permite gerar automaticamente documentação OpenAPI com base em nossos comentários de código. Tudo o que precisamos fazer é escrever os comentários antes de nossas funções de manipulação. Por exemplo:
```go
// GetProfile godoc
// @Summary get user info.
// @Schemes
// @Description
// @Tags 用户模块
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} response.Response
// @Router /user [get]
func (h *userHandler) GetProfile(ctx *gin.Context) {
    // ...
}
```

Execute o comando `swag init` para gerar os arquivos de documentação:
```bash
swag init -g cmd/server/main.go -o ./docs --parseDependency

# ou
make swag
```

Abra a página de documentação no seu navegador:
```bash
http://127.0.0.1:8000/swagger/index.html
```


## Conclusão

O framework Nunu oferece uma estrutura de projeto elegante e operações de comando que permitem aos desenvolvedores desenvolver aplicações web de forma eficiente. Neste tutorial, você aprendeu como criar um projeto, criar Handlers, criar Services, criar Repositories, compilar o Wire e iniciar o projeto usando o Nunu. Esperamos que este conteúdo ajude você a fazer melhor uso do framework Nunu.
