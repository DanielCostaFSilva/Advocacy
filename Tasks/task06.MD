Você é um engenheiro de software sênior especializado em Golang, Clean Architecture, DDD, TDD e APIs REST.

Contexto do projeto:

Sistema SaaS para escritórios de advocacia.

Stack:

- Golang 1.26+
- Chi Router
- PostgreSQL
- SQLC
- Goose
- Testify
- Clean Architecture
- Modular Monolith

Estado atual do projeto:

Já existem:

✓ User Entity
✓ User Repository
✓ PostgreSQL Persistence
✓ CreateUser Use Case
✓ Password Hasher Interface
✓ Structured Logging
✓ Request ID Middleware
✓ Recovery Middleware

Objetivo desta task:

Expor o caso de uso CreateUser através de uma API REST.

IMPORTANTE:

Ainda NÃO implementar:

- Login
- JWT
- Middleware de autenticação
- Refresh Token
- RBAC
- Multi-tenant

Esta task deve implementar apenas a camada HTTP responsável pelo cadastro de usuários.

Regras obrigatórias:

1. Trabalhar usando TDD.
2. Criar testes antes da implementação.
3. Implementar apenas o necessário.
4. Não antecipar funcionalidades futuras.
5. Código em inglês.
6. Seguir SOLID.
7. Seguir Clean Architecture.
8. Handler não pode acessar banco diretamente.
9. Handler deve depender apenas do Use Case.
10. Não colocar regra de negócio no Handler.

Objetivos da implementação:

1. Endpoint POST /users
2. Request DTO
3. Response DTO
4. Handler HTTP
5. Registro de rota
6. Testes de Handler
7. Tratamento padronizado de erros

Estrutura esperada:

backend/
├── internal/
│
│   ├── http/
│   │   └── user/
│   │       ├── handler.go
│   │       ├── handler_test.go
│   │       ├── request.go
│   │       ├── response.go
│   │       └── routes.go
│   │
│   ├── application/
│   │   └── user/
│   │
│   └── domain/
│
└── cmd/
    └── api/
        └── main.go

Endpoint esperado:

POST /users

Request:

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "12345678"
}

Response:

HTTP 201

{
  "id": "uuid",
  "name": "John Doe",
  "email": "john@example.com"
}

Headers:

Content-Type: application/json

Validações HTTP:

- body obrigatório
- JSON válido
- name obrigatório
- email obrigatório
- password obrigatório

As regras de negócio continuam no Use Case.

O Handler deve apenas:

1. Ler request
2. Validar estrutura mínima
3. Chamar Use Case
4. Traduzir resultado para HTTP

Mapeamento de erros:

ErrEmailAlreadyExists

HTTP 409

{
  "error": "email already exists"
}

ErrInvalidInput

HTTP 400

{
  "error": "invalid input"
}

Erro inesperado

HTTP 500

{
  "error": "internal server error"
}

Interface esperada:

type CreateUserUseCase interface {
    Execute(
        ctx context.Context,
        input CreateUserInput,
    ) (
        CreateUserOutput,
        error,
    )
}

Testes obrigatórios:

Scenario:
Given valid request
When POST /users
Then should return 201

Scenario:
Given duplicated email
When POST /users
Then should return 409

Scenario:
Given invalid body
When POST /users
Then should return 400

Scenario:
Given malformed JSON
When POST /users
Then should return 400

Scenario:
Given internal error
When POST /users
Then should return 500

Scenario:
Given valid request
When POST /users
Then should invoke CreateUser Use Case

Mocks esperados:

CreateUserUseCaseMock

Não utilizar banco de dados.

Não utilizar PostgreSQL.

Todos os testes devem ser unitários.

Bootstrap esperado:

main.go
  ↓
router
  ↓
user routes
  ↓
handler
  ↓
use case

Critérios de aceite:

- Endpoint funcional.
- Testes passando.
- Respostas padronizadas.
- Handler sem regra de negócio.
- Dependência invertida corretamente.

Entregáveis:

1. Handler completo.
2. DTOs request/response.
3. Registro de rotas.
4. Testes unitários.
5. Exemplo de requests/responses.
6. Explicação rápida da arquitetura.

Não implementar:

- Login
- JWT
- Refresh Token
- Middleware Auth
- Tenant
- Frontend
- Swagger

Somente o endpoint de cadastro de usuário.