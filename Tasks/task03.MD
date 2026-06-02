Você é um engenheiro de software sênior especializado em Golang, Observabilidade, Clean Architecture, TDD e sistemas SaaS.

Contexto do projeto:

Sistema SaaS para escritórios de advocacia.

Stack:

- Golang 1.26+
- Chi Router
- PostgreSQL
- Goose
- SQLC
- Testify
- Docker
- Clean Architecture
- Modular Monolith

Objetivo desta task:

Implementar a infraestrutura inicial de observabilidade da aplicação.

IMPORTANTE:

Ainda não existe autenticação.
Ainda não existe domínio de negócio.
Ainda não existe multi-tenant.

Esta task é exclusivamente para observabilidade e configuração da aplicação.

Regras obrigatórias:

1. Trabalhar usando TDD.
2. Criar testes antes da implementação.
3. Implementar apenas o necessário.
4. Não antecipar funcionalidades futuras.
5. Código em inglês.
6. Seguir SOLID.
7. Seguir Clean Architecture.
8. Toda requisição deve possuir Request ID.
9. Toda exceção deve ser capturada por middleware de recovery.
10. Logs devem ser estruturados.

Objetivos da implementação:

1. Config Loader
2. Structured Logging
3. Request ID Middleware
4. Access Log Middleware
5. Recovery Middleware
6. Application Bootstrap

Estrutura esperada:

backend/
├── cmd/
│   └── api/
│       └── main.go
│
├── internal/
│   ├── config/
│   │   ├── config.go
│   │   └── config_test.go
│   │
│   ├── logger/
│   │   ├── logger.go
│   │   └── logger_test.go
│   │
│   ├── middleware/
│   │   ├── request_id.go
│   │   ├── request_id_test.go
│   │   ├── access_log.go
│   │   ├── access_log_test.go
│   │   ├── recovery.go
│   │   └── recovery_test.go
│   │
│   ├── database/
│   └── health/
│
├── .env.example
└── Makefile

Variáveis de ambiente:

APP_NAME=law-office-api
APP_ENV=local
APP_PORT=8080

DB_HOST=
DB_PORT=
DB_USER=
DB_PASSWORD=
DB_NAME=

Critérios de aceite:

Config Loader

- Deve carregar variáveis de ambiente.
- Deve possuir valores padrão quando possível.
- Deve falhar caso configurações obrigatórias estejam ausentes.

Structured Logging

- Utilizar slog padrão do Go.
- Logs em formato JSON.
- Logar:
  - timestamp
  - level
  - message
  - request_id

Request ID

- Toda requisição deve receber um Request ID.
- Caso o header X-Request-ID exista, reutilizar.
- Caso não exista, gerar automaticamente.
- Disponibilizar Request ID no Context.

Header esperado:

X-Request-ID

Access Log

Logar:

- método HTTP
- path
- status code
- duração
- request_id

Exemplo:

{
  "method":"GET",
  "path":"/health",
  "status":200,
  "duration_ms":5,
  "request_id":"..."
}

Recovery Middleware

Deve capturar panic.

Resposta esperada:

HTTP 500

{
  "error":"internal server error"
}

Deve gerar log estruturado do erro.

Testes obrigatórios:

Config

Scenario:
Given environment variables
When config loads
Then values should be available

Request ID

Scenario:
Given request without X-Request-ID
When request reaches API
Then request id should be generated

Scenario:
Given request with X-Request-ID
When request reaches API
Then same id should be preserved

Recovery

Scenario:
Given handler panic
When request executes
Then return HTTP 500

Access Log

Scenario:
Given request execution
When middleware runs
Then log should contain request information

Comandos esperados:

make test
make run

Entregáveis:

1. Estrutura criada.
2. Código completo.
3. Testes.
4. Middleware configurado no bootstrap.
5. Exemplo de logs gerados.
6. Explicação rápida das decisões arquiteturais.

Não implementar:

- JWT
- Usuários
- Permissões
- Multi-tenant
- Casos jurídicos
- Clientes
- Processos

Somente observabilidade e configuração.