## Documentação
* [Guia do Usuário](https://github.com/go-nunu/nunu/blob/main/docs/pt/guide.md)
* [Arquitetura](https://github.com/go-nunu/nunu/blob/main/docs/pt/architecture.md)
* [Tutorial de Início Rápido](https://github.com/go-nunu/nunu/blob/main/docs/pt/tutorial.md)
* [Teste de Unidade](https://github.com/go-nunu/nunu/blob/main/docs/pt/unit_testing.md)

- [Chinês](https://github.com/go-nunu/nunu/blob/main/docs/zh/architecture.md)
- [Português](https://github.com/go-nunu/nunu/blob/main/docs/pt/architecture.md)

# Explorando a Arquitetura do Nunu

O Nunu adota uma arquitetura em camadas clássica. Além disso, para alcançar uma melhor modularidade e desacoplamento, ele utiliza o framework de injeção de dependência `Wire`.


![Nunu Layout](https://github.com/go-nunu/nunu/blob/main/.github/assets/layout.png)

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


- `cmd`: Ponto de entrada da aplicação, contendo diferentes subcomandos.
- `config`: Arquivos de configuração.
- `deploy`: Arquivos relacionados à implantação, como Dockerfile e docker-compose.yml.
- `internal`: Código principal da aplicação, organizado de acordo com a arquitetura em camadas.
- `pkg`: Código comum, incluindo configuração, registro de eventos e HTTP.
- `scripts`: Arquivos de script para implantação e outras tarefas de automação.
- `storage`: Arquivos de armazenamento, como arquivos de log.
- `test`: Código de teste.
- `web`: Código de front-end.

## internal

- `internal/handler` (ou `controller`): Lida com solicitações HTTP, chama serviços na camada de lógica de negócios e retorna respostas HTTP.
- `internal/server` (ou `router`): Servidor HTTP que inicia o serviço HTTP, ouve portas e lida com solicitações HTTP.
- `internal/service` (ou `logic`): Serviços que implementam lógica de negócios específica e chamam a camada de acesso aos dados (repositório).
- `internal/model` (ou `entity`): Modelos de dados que definem as estruturas de dados necessárias pela camada de lógica de negócios.
- `internal/repository` (ou `dao`): Objetos de acesso aos dados que encapsulam operações de banco de dados e fornecem operações CRUD nos dados.

## Injeção de Dependência

Este projeto utiliza o framework de injeção de dependência `Wire` para alcançar modularidade e desacoplamento. `Wire` gera código de injeção de dependência `wire_gen.go` pré-compilando `wire.go`, simplificando o processo de injeção de dependência.

- `cmd/job/wire.go`: Arquivo de configuração do `Wire` que define as dependências necessárias pelo subcomando `job`.
- `cmd/migration/wire.go`: Arquivo de configuração do `Wire` que define as dependências necessárias pelo subcomando `migration`.
- `cmd/server/wire.go`: Arquivo de configuração do `Wire` que define as dependências necessárias pelo subcomando `server`.

Documentação oficial do Wire: https://github.com/google/wire/blob/main/docs/guide.md

Nota: O arquivo `wire_gen.go` é gerado automaticamente durante a compilação e não deve ser modificado manualmente.

## Código Comum

Para alcançar a reutilização de código e gerenciamento centralizado, este projeto adota uma abordagem de código comum, onde algum código comum é colocado no diretório `pkg`.

- `pkg/config`: Lida com a leitura e análise de arquivos de configuração.
- `pkg/helper`: Contém várias funções utilitárias, como criptografia MD5 e geração de UUID.
- `pkg/http`: Contém código relacionado ao HTTP, como clientes e servidores HTTP.
- `pkg/log`: Contém código relacionado a registro de eventos, como inicialização e escrita de logs.
- `mais...`: Claro, você pode adicionar e expandir livremente mais pacotes conforme necessário.