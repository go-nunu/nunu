## Documentação
* [Guia do Usuário](https://github.com/go-nunu/nunu/blob/main/docs/pt/guide.md)
* [Arquitetura](https://github.com/go-nunu/nunu/blob/main/docs/pt/architecture.md)
* [Tutorial de Início Rápido](https://github.com/go-nunu/nunu/blob/main/docs/pt/tutorial.md)
* [Teste de Unidade](https://github.com/go-nunu/nunu/blob/main/docs/pt/unit_testing.md)

- [Chinês](https://github.com/go-nunu/nunu/blob/main/docs/zh/unit_testing.md)
- [Português](https://github.com/go-nunu/nunu/blob/main/docs/pt/unit_testing.md)


# Testes Unitários

## Introdução

Testes unitários são uma prática importante em projetos de desenvolvimento. No entanto, escrever testes unitários se torna complexo e instável quando o código testado depende de outros módulos ou componentes. Este artigo apresentará como usar mocks para escrever testes unitários concisos e eficientes.

## Visão Geral

Primeiro, vamos dar uma olhada no arquivo de injeção de dependência no projeto `cmd/server/wire.go`:

> Dica: Este arquivo é compilado e gerado automaticamente pela ferramenta `google/wire` e não deve ser editado manualmente.

```go
// Injetores de wire.go:

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

Deste trecho de código, podemos ver as relações de dependência entre `handler`, `service` e `repository`.

`userHandler` depende de `userService`, e `userService` depende de `userRepository`.

Por exemplo, o código para `GetProfile` em `handler/user.go` é o seguinte:
```go
func (h *userHandler) GetProfile(ctx *gin.Context) {
	userId := GetUserIdFromCtx(ctx)
	if userId == "" {
		v1.HandleError(ctx, http.StatusUnauthorized, 1, "unauthorized", nil)
		return
	}

	user, err := h.userService.GetProfile(ctx, userId)
	if err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}

	resp.HandleSuccess(ctx, user)
}
```
Podemos ver que ele chama internamente `userService.GetProfile`.

Portanto, ao escrever testes unitários, inevitavelmente precisamos inicializar a instância de `userService` primeiro. No entanto, ao inicializarmos `userService`, descobrimos que ele depende de `userRepository`.

Embora só precisemos testar o `handler` de nível inferior, precisamos inicializar e executar `service`, `repository` e outros códigos. Isso obviamente viola o princípio do teste unitário (Princípio da Responsabilidade Única), onde cada teste unitário deve se concentrar em uma funcionalidade ou unidade de código específica.

Qual é a boa solução para este problema? Nossa resposta final é "mocking".


### Mocking (Um Bom Auxiliar para Isolamento de Dependências)

Ao realizar testes unitários, queremos testar a lógica da unidade de código testada sem depender do estado ou comportamento de outros módulos ou componentes externos. Esta abordagem pode isolar melhor o código testado e tornar os testes mais confiáveis e repetíveis.

Mocking é um padrão de teste usado para simular ou substituir módulos ou componentes externos dos quais o código testado depende. Ao usar objetos mock, podemos controlar o comportamento dos módulos externos, de modo que o código testado não precise realmente depender e chamar os módulos externos durante o teste, alcançando assim o isolamento do código testado.

Objetos mock podem simular valores de retorno, exceções, timeouts, etc., de módulos externos, tornando os testes mais controláveis e previsíveis. Ele resolve os seguintes problemas:

1. Dependência de outros módulos: Algumas unidades de código podem depender de outros módulos, como bancos de dados, solicitações de rede, etc. Ao usar objetos mock, podemos simular essas dependências, para que os testes não precisem realmente depender desses módulos, evitando assim a instabilidade e complexidade dos testes.

2. Isolamento do ambiente externo: Algumas unidades de código podem ser afetadas pelo ambiente externo, como o tempo atual, status do sistema, etc. Ao usar objetos mock, podemos controlar o estado desses ambientes externos, para que os testes possam ser executados em diferentes ambientes, aumentando assim a cobertura e precisão dos testes.

3. Melhoria da eficiência dos testes: Alguns módulos externos podem realizar operações demoradas, como solicitações de rede, operações de leitura/escrita de arquivos, etc. Ao usar objetos mock, podemos evitar executar essas operações na realidade, melhorando assim a velocidade e eficiência de execução dos testes.


No projeto nunu, usamos as seguintes bibliotecas de mocking para nos ajudar a escrever testes unitários:

* github.com/golang/mock            // Uma biblioteca de mocking open-source do Google
* github.com/go-redis/redismock/v9  // Fornece testes de mock para consultas Redis, compatível com github.com/redis/go-redis/v9
* github.com/DATA-DOG/go-sqlmock    // sqlmock é uma biblioteca de mocking que implementa sql/driver


## Programação Orientada a Interfaces

Usar `golang/mock` tem um pré-requisito. Precisamos seguir a abordagem de "programação orientada a interfaces" para escrever nosso `repository` e `service`.

Alguns podem não estar familiarizados com o que significa "programação orientada a interfaces". Vamos pegar um trecho de código como exemplo:

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

No código acima, primeiro definimos uma `UserRepository interface`, e então implementamos todos os seus métodos usando a `userRepository struct`.

```go
type UserRepository interface {
	FirstById(id int64) (*model.User, error)
}
type userRepository struct {
	*Repository
}
func (r *userRepository) FirstById(id int64) (*model.User, error) {
    // ...
}

```
Instead of directly writing it as:
```go
type UserRepository struct {
	*Repository
}

func (r *UserRepository) FirstById(id int64) (*model.User, error) {
    // ...
}
```

A **programação orientada a interfaces**, que pode melhorar a flexibilidade, escalabilidade, testabilidade e manutenção do código, é um estilo de programação altamente recomendado pela linguagem Go.


## Começando com go-mock

Usar o `golang/mock` é simples. Primeiro, vamos instalá-lo:

```bash
go install github.com/golang/mock/mockgen@v1.6.0
```

`mockgen` é uma ferramenta de linha de comando para `go-mock` que pode analisar as definições de `interface` em nosso código e gerar o código mock correto.


Exemplo:
```bash
mockgen -source=internal/service/user.go -destination mocks/service/user.go
```

O comando acima especifica dois parâmetros: o arquivo fonte da interface e o arquivo de destino onde o código mock gerado será colocado. Colocamos o arquivo alvo no diretório `mocks/service`.

Após gerar o código mock para `UserService`, podemos escrever testes unitários para `UserHandler`.

O código final do teste unitário é o seguinte:

```go

func TestUserHandler_GetProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_service.NewMockUserService(ctrl)
	
	// Key code, define o valor de retorno de mockUserService.GetProfile.
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
	// Adicione asserções para o corpo da resposta, se necessário.
}

```

O código fonte completo está localizado em: https://github.com/go-nunu/nunu-layout-advanced/blob/main/test/server/handler/user_test.go

## sqlmock e redismock

Para testes unitários de `repository`, que dependem não dos nossos próprios módulos de negócios, mas de fontes de dados externas como RPC, Redis e MySQL, é um pouco diferente de testar `handler` e `service`, pois precisamos evitar a conexão com bancos de dados e caches reais para reduzir incertezas nos testes. Portanto, também usamos mocks neste caso.

O código é o seguinte:
```go
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

    // Simular a consulta de dados de teste.
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

O código completo está localizado em: https://github.com/go-nunu/nunu-layout-advanced/blob/main/test/server/repository/user_test.go


## Cobertura de Testes
Golang suporta nativamente a geração de relatórios de cobertura de testes.

```bash
go test -coverpkg=./internal/handler,./internal/service,./internal/repository -coverprofile=./coverage.out ./test/server/...

go tool cover -html=./coverage.out -o coverage.html
```

Os dois comandos acima gerarão um arquivo de relatório de cobertura `coverage.html` em um formato de visualização web, que pode ser aberto diretamente em um navegador.

O efeito é o seguinte:

![coverage](https://github.com/go-nunu/nunu/blob/main/.github/assets/coverage.png)

## Conclusão


Testes unitários são uma prática de desenvolvimento importante em projetos, pois garantem a correção do código e fornecem validação automatizada. Ao conduzir testes unitários, precisamos usar programação orientada a interfaces e objetos mock para isolar as dependências do código testado. Na linguagem Go, podemos usar a biblioteca golang/mock para gerar código mock. Para repositórios que dependem de fontes de dados externas, podemos usar sqlmock e redismock para simular o comportamento de bancos de dados e caches. Ao usar objetos mock, podemos controlar o comportamento de módulos externos, permitindo que o código testado não dependa verdadeiramente e chame módulos externos durante os testes, alcançando assim a isolação do código testado. Isso melhora a confiabilidade, repetibilidade e eficiência dos testes.