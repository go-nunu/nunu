# Nunu - Uma ferramenta de linha de comando (CLI) para construir aplicações em Go.

Nunu é uma ferramenta de geração de estrutura (scaffolding) para construir aplicações em Go. Seu nome vem de um personagem de um jogo chamado League of Legends, um garotinho montado nos ombros de um Yeti. Assim como Nunu, este projeto se apoia nos ombros de gigantes, pois é construído sobre uma combinação de bibliotecas populares do ecossistema Go. Essa combinação permite que você construa rapidamente aplicações eficientes e confiáveis.

🚀Dicas: Este projeto é muito completo, então as atualizações não serão muito frequentes. Sinta-se à vontade para utilizá-lo.

- [简体中文介绍](https://github.com/go-nunu/nunu/blob/main/README_zh.md)
- [Português](https://github.com/go-nunu/nunu/blob/main/README_pt.md)

![Nunu](https://github.com/go-nunu/nunu/blob/main/.github/assets/banner.png)

## Documentação
* [Guia do Usuário](https://github.com/go-nunu/nunu/blob/main/docs/pt/guide.md)
* [Arquitetura](https://github.com/go-nunu/nunu/blob/main/docs/pt/architecture.md)
* [Tutorial de Início Rápido](https://github.com/go-nunu/nunu/blob/main/docs/pt/tutorial.md)
* [Teste de Unidade](https://github.com/go-nunu/nunu/blob/main/docs/pt/unit_testing.md)
* [MCP Server](https://github.com/go-nunu/nunu-layout-mcp/blob/main/README.md)
* [Monorepo Layout](https://github.com/go-nunu/nunu-layout-monorepo)



## Funcionalidades
- **Gin**: https://github.com/gin-gonic/gin
- **Gorm**: https://github.com/go-gorm/gorm
- **Wire**: https://github.com/google/wire
- **Viper**: https://github.com/spf13/viper
- **Zap**: https://github.com/uber-go/zap
- **Golang-jwt**: https://github.com/golang-jwt/jwt
- **Go-redis**: https://github.com/go-redis/redis
- **Testify**: https://github.com/stretchr/testify
- **Sonyflake**: https://github.com/sony/sonyflake
- **Gocron**:  https://github.com/go-co-op/gocron
- **Go-sqlmock**:  https://github.com/DATA-DOG/go-sqlmock
- **Gomock**:  https://github.com/golang/mock
- **Swaggo**:  https://github.com/swaggo/swag
- **Casbin**:  https://github.com/casbin/casbin
- **Pitaya**:  https://github.com/topfreegames/pitaya
- **MCP-GO**:  https://github.com/mark3labs/mcp-go

- Mais...

## Funcionalidades Principais
* **Curva de Aprendizado Baixa e Personalização**: Nunu encapsula bibliotecas populares que os Gophers estão familiarizados, permitindo que você customize facilmente a aplicação para atender a requisitos específicos.
* **Alto Desempenho e Escalabilidade**: Nunu tem como objetivo ser de alto desempenho e escalável. Ele utiliza as tecnologias mais recentes e as melhores práticas para garantir que sua aplicação possa lidar com alto tráfego e grandes quantidades de dados.
* **Segurança e Confiabilidade**: O Nunu utiliza bibliotecas de terceiros estáveis e confiáveis para garantir a segurança e confiabilidade da sua aplicação.
* **Modular e Extensível**: O Nunu foi projetado para ser modular e extensível. Você pode facilmente adicionar novos recursos e funcionalidades usando bibliotecas de terceiros ou escrevendo seus próprios módulos.
* **Documentação Completa e Testes**: O Nunu possui documentação completa e testes abrangentes. Ele oferece documentação extensa e exemplos para ajudá-lo a começar rapidamente. Também inclui um conjunto de testes para garantir que sua aplicação funcione conforme o esperado.

## Arquitetura em Camadas Concisa
O Nunu adota uma arquitetura em camadas clássica. Para alcançar modularidade e desacoplamento, ele utiliza o framework de injeção de dependência `Wire`.

![Nunu Layout](https://github.com/go-nunu/nunu/blob/main/.github/assets/layout.png)

## Nunu CLI

![Nunu](https://github.com/go-nunu/nunu/blob/main/.github/assets/screenshot.jpg)


## Estrutura de Diretórios
```
.
├── api
│   └── v1
├── cmd
│   ├── migration
│   ├── server
│   │   ├── wire
│   │   │   ├── wire.go
│   │   │   └── wire_gen.go
│   │   └── main.go
│   └── task
├── config
├── deploy
├── docs
├── internal
│   ├── handler
│   ├── middleware
│   ├── model
│   ├── repository
│   ├── server
│   └── service
├── pkg
├── scripts
├── test
│   ├── mocks
│   └── server
├── web
├── Makefile
├── go.mod
└── go.sum

```

A arquitetura do projeto segue uma estrutura em camadas típica, consistindo nos seguintes módulos:

* `cmd`: Este módulo contém os pontos de entrada da aplicação, que realizam diferentes operações com base em comandos diferentes, como iniciar o servidor ou executar migrações de banco de dados. Cada submódulo tem um arquivo `main.go` como arquivo de entrada, além dos arquivos `wire.go` e `wire_gen.go` para injeção de dependência.

* `config`: Este módulo contém os arquivos de configuração da aplicação, fornecendo diferentes configurações para ambientes diferentes, como desenvolvimento e produção.

* `deploy`: Este módulo é usado para implantar a aplicação e inclui scripts de implantação e arquivos de configuração.

* `internal`: Este módulo é o módulo central da aplicação e contém a implementação de várias lógicas de negócios.

  - `handler`: Este submódulo contém os manipuladores para lidar com solicitações HTTP, responsáveis por receber solicitações e invocar os serviços correspondentes para processamento.

  - `job`: Este submódulo contém a lógica para tarefas em segundo plano.

  - `model`: Este submódulo contém a definição de modelos de dados.

  - `repository`: Este submódulo contém a implementação da camada de acesso a dados, responsável por interagir com o banco de dados.

  - `server`: Este submódulo contém a implementação do servidor HTTP.

  - `service`: Este submódulo contém a implementação da lógica de negócios, responsável por lidar com operações de negócios específicas.

* `pkg`: Este módulo contém algumas utilidades e funções comuns.

* `scripts`: Este módulo contém scripts usados para compilação, teste e operações de implantação do projeto.

* `storage`: Este módulo é usado para armazenar arquivos ou outros recursos estáticos.

* `test`: Este módulo contém testes unitários para diversos módulos, organizados em subdiretórios com base nos módulos.

* `web`: Este módulo contém os arquivos relacionados ao frontend, como HTML, CSS e JavaScript.

Além disso, existem outros arquivos e diretórios, como arquivos de licença, arquivos de construção e README. No geral, a arquitetura do projeto é clara, com responsabilidades claras para cada módulo, facilitando o entendimento e a manutenção.

## Requisitos
Para usar o Nunu, você precisa ter o seguinte software instalado em seu sistema:

* Go 1.19 ou superior
* Git
* Docker (opcional)
* MySQL 5.7 ou superior (opcional)
* Redis (opcional)

### Instalação

Você pode instalar o Nunu com o seguinte comando:

```bash
go install github.com/go-nunu/nunu@latest
```

> Dicas: Se `go install` for bem-sucedido, mas o comando `nunu` não é reconhecido, é porque a variável de ambiente não está configurada. Você pode adicionar o diretório GOBIN à variável de ambiente.

## Criar um Novo Projeto

Você pode criar um novo projeto em Go com o seguinte comando:

```bash
nunu new projectName
```

Por padrão, ele busca no repositório do GitHub, mas você também pode usar um repositório acelerado na China:

```bash
# Usar o modelo básico
nunu new nomeDoProjeto -r https://gitee.com/go-nunu/nunu-layout-basic.git
# Usar o modelo avançado
nunu new nomeDoProjeto -r https://gitee.com/go-nunu/nunu-layout-advanced.git
# Usar o modelo monorepo
nunu new nomeDoProjeto -r https://github.com/go-nunu/nunu-layout-monorepo.git
```

Este comando criará um diretório chamado `nomeDoProjeto` e gerará uma estrutura de projeto elegante em Go dentro dele.

### Criar Componentes

Você pode criar handlers, services, repositories e models para o seu projeto usando os seguintes comandos:

```bash
nunu create handler user
nunu create service user
nunu create repository user
nunu create model user
```
ou
```bash
nunu create all user
```

Estes comandos criarão os componentes com os nomes `UserHandler`, `UserService`, `UserRepository` e `UserModel`, respectivamente, e os colocarão nos diretórios corretos.

### Executar o Projeto

Você pode iniciar rapidamente o projeto com o seguinte comando:

```bash
nunu run
```

Este comando iniciará o seu projeto em Go e oferecerá suporte a recarregamento automático (hot-reloading) quando os arquivos forem atualizados.

### Compilar wire.go

Você pode compilar rapidamente o arquivo wire.go com o seguinte comando:

```bash
nunu wire
```

Este comando irá compilar o seu `arquivo` wire.go e gerar as dependências necessárias.

## Contribuição

Se encontrar algum problema ou tiver sugestões de melhoria, sinta-se à vontade para abrir um problema ou enviar um pull request. Suas contribuições são altamente apreciadas!

## Licença

O Nunu é lançado sob a Licença MIT. Para mais informações, consulte o arquivo de [LICENSE](LICENSE).

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=go-nunu/nunu&type=Date)](https://star-history.com/#go-nunu/nunu&Date)
