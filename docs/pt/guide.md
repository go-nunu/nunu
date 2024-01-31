## Documentação
* [Guia do Usuário](https://github.com/go-nunu/nunu/blob/main/docs/pt/guide.md)
* [Arquitetura](https://github.com/go-nunu/nunu/blob/main/docs/pt/architecture.md)
* [Tutorial de Início Rápido](https://github.com/go-nunu/nunu/blob/main/docs/pt/tutorial.md)
* [Teste de Unidade](https://github.com/go-nunu/nunu/blob/main/docs/pt/unit_testing.md)

- [Chinês](https://github.com/go-nunu/nunu/blob/main/docs/zh/guide.md)
- [Português](https://github.com/go-nunu/nunu/blob/main/docs/pt/guide.md)

# Guia do Usuário Nunu

Nunu é uma estrutura de aplicação baseada em Golang que ajuda você a construir aplicações eficientes e confiáveis rapidamente. Este guia mostrará como usar o Nunu para criar e desenvolver suas aplicações.

## Instalação

Você pode instalar o Nunu usando o seguinte comando:

```bash
go install github.com/go-nunu/nunu@latest
```

> Dica: Se `go install` for bem-sucedido, mas o comando nunu não for encontrado, é porque a variável de ambiente não está configurada. Você pode configurar o diretório GOBIN na variável de ambiente.

## Criando um Novo Projeto

Você pode usar o seguinte comando para criar um novo projeto em Golang:

```bash
nunu new nomeDoProjeto
```

Este comando criará um diretório chamado `nomeDoProjeto` e gerará uma estrutura elegante de projeto em Golang dentro dele.

**Fonte Acelerada na China:**

Por padrão, `nunu new` busca a partir da fonte do GitHub. No entanto, você também pode usar um repositório acelerado na China:
```bash
# Usando o template básico
nunu new nomeDoProjeto  -r https://gitee.com/go-nunu/nunu-layout-basic.git
# Usando o template avançado (recomendado)
nunu new nomeDoProjeto  -r https://gitee.com/go-nunu/nunu-layout-advanced.git
```

> O Nunu oferece dois tipos de layouts:

* **Layout Básico**

O Layout Básico contém uma estrutura de diretórios minimalista e é adequado para desenvolvedores que já estão familiarizados com projetos Nunu.

* **Layout Avançado**

**Recomendação: Recomendamos que iniciantes escolham o Layout Avançado primeiro.**

O Layout Avançado inclui muitos exemplos de uso do Nunu (por exemplo, db, redis, jwt, cron, migração, etc.), o que é adequado para desenvolvedores aprenderem rapidamente e entenderem as ideias arquitetônicas do Nunu.


## Início Rápido com Docker

Se você quiser experimentar rapidamente o layout avançado do Nunu, recomendamos usar os seguintes comandos para iniciar o projeto rapidamente:

```bash
cd ./deploy/docker-compose && docker compose up -d && cd ../../

go run ./cmd/migration

nunu run ./cmd/server
```

Alternativamente, você pode usar diretamente o comando `make`:

```bash
make bootstrap
```

## Criando Componentes

Você pode usar os seguintes comandos para criar componentes como handler, service, repository e model para o seu projeto:

```bash
nunu create handler user
nunu create service user
nunu create repository user
nunu create model user
```

Esses comandos criarão componentes chamados `UserHandler`, `UserService`, `UserRepository` e `UserModel`, e os colocarão nos diretórios corretos.

Se você quiser criar os componentes correspondentes em um diretório personalizado, você pode fazer isso da seguinte forma:
```bash
nunu create handler internal/handler/user/center
nunu create service internal/service/user/center
nunu create repository internal/repository/user/center
nunu create model internal/model/user/center
```

Você também pode usar o seguinte comando para criar todos os componentes (`handler`, `service`, `repository` e `model`) de uma vez:

```bash
nunu create all user
```

## Iniciando o Projeto

Você pode iniciar rapidamente o projeto usando o seguinte comando:

```bash
# Equivalente a go run ./cmd/server

nunu run
```

Este comando iniciará seu projeto Golang e suportará recarga automática quando os arquivos forem atualizados.

## Compilando wire.go

Você pode compilar rapidamente `wire.go` usando o seguinte comando:

```bash
# Equivalente a cd cmd/server && wire
nunu wire
```

Este comando procurará automaticamente pelo arquivo `wire.go` em seu projeto e compilará as dependências necessárias.

## Arquivo de Configuração

### Iniciando com um Arquivo de Configuração Específico
O Nunu usa a biblioteca Viper para gerenciar arquivos de configuração.

Por padrão, ele carrega `config/local.yml`, mas você pode especificar o caminho do arquivo de configuração usando variáveis de ambiente ou parâmetros.

```bash
# Linux ou MacOS
APP_CONF=config/prod.yml nunu run

# Windows
set APP_CONF=config\prod.yml && nunu run
```
Alternativamente, você pode usar a abordagem de parâmetro: `go run ./cmd/server -conf=config/prod.yml`

### Lendo Itens de Configuração

Você pode criar um arquivo chamado `local.yaml` no diretório `config` para armazenar suas informações de configuração. Por exemplo:

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

Você pode usar a injeção de dependência `conf *viper.Viper` para ler as informações de configuração no seu código:

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
Dica: Após realizar a injeção de dependência através de parâmetros, não se esqueça de executar o comando nunu wire para gerar o arquivo de dependência.

## Registro de Logs

O Nunu utiliza a biblioteca Zap para gerenciar logs. Você pode configurar o log no diretório `config`. Por exemplo:

```yaml
log:
  log_level: info
  encoding: json           	   # json ou console
  log_file_name: "./storage/logs/server.log"
  max_backups: 30              # Número máximo de backups de arquivos de log
  max_age: 7                   # Número máximo de dias para manter os arquivos
  max_size: 1024               # Tamanho máximo de cada arquivo de log em MB
  compress: true               # Se os arquivos de log devem ser comprimidos
```

Você pode usar o seguinte método para registrar logs no seu código:

```go
package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-nunu/nunu-layout-basic/internal/service"
	"github.com/go-nunu/nunu-layout-basic/pkg/helper/resp"
	"go.uber.org/zap"
	"net/http"
)

// ...

func (h *userHandler) GetUserById(ctx *gin.Context) {
	h.logger.Info("GetUserByID", zap.Any("user", user))
	// ...
}

// ...

```

## Banco de Dados

O Nunu utiliza a biblioteca GORM para gerenciar bancos de dados. Você pode configurar o banco de dados no diretório `config`. Por exemplo:

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

Você pode se conectar ao banco de dados usando o seguinte código:

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


É importante notar que `xxxRepository`, `xxxService`, `xxxHandler`, etc., no Nunu são implementados com base em interfaces. Isso é conhecido como **programação orientada a interfaces**, que pode melhorar a flexibilidade, escalabilidade, testabilidade e manutenibilidade do código. É um estilo de programação altamente recomendado na linguagem Go.

No código acima, escrevemos:

```go
type UserRepository interface {
	FirstById(id int64) (*model.User, error)
}
type userRepository struct {
	*Repository
}
```
em vez de escrever diretamente:
```go
type UserRepository struct {
	*Repository
}
```
> Dica: Os testes unitários no layout avançado do Nunu são baseados nas características de `interface` para operações de mock.

## Testes

O Nunu utiliza bibliotecas como testify, redismock, gomock e go-sqlmock para escrever testes.

Você pode executar o teste usando o seguinte comando:

```bash
go test -coverpkg=./internal/handler,./internal/service,./internal/repository -coverprofile=./.nunu/coverage.out ./test/server/...
go tool cover -html=./.nunu/coverage.out -o coverage.html

```

O comando acima irá gerar um arquivo HTML chamado `coverage.html`. Você pode abri-lo diretamente em um navegador para visualizar a cobertura detalhada dos testes unitários.

## Conclusão

O Nunu é um esqueleto de aplicação prático em Golang que ajuda você a construir rapidamente aplicações eficientes e confiáveis. Esperamos que este guia possa ajudá-lo a fazer um melhor uso do Nunu.