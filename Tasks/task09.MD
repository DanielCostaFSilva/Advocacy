Você é um engenheiro de software sênior especializado em Golang, APIs REST, Segurança, JWT, Clean Architecture, DDD e TDD.

Contexto do projeto:

Sistema SaaS para escritórios de advocacia.

Stack:

- Golang 1.26+
- Chi Router
- PostgreSQL
- SQLC
- Goose
- Testify
- JWT
- Clean Architecture
- Modular Monolith

Estado atual do projeto:

Já existem:

✓ User Entity
✓ User Repository
✓ CreateUser Use Case
✓ AuthenticateUser Use Case
✓ JWT Token Provider
✓ Structured Logging
✓ Request ID Middleware
✓ Recovery Middleware
✓ POST /users

Objetivo desta task:

Implementar o endpoint HTTP responsável pelo login de usuários.

IMPORTANTE:

Ainda NÃO implementar:

- Refresh Token
- Logout
- Middleware de autenticação
- RBAC
- Multi-tenant
- OAuth

Esta task deve apenas autenticar o usuário e retornar um Access Token JWT.

Regras obrigatórias:

1. Trabalhar usando TDD.
2. Criar testes antes da implementação.
3. Implementar apenas o necessário.
4. Não antecipar funcionalidades futuras.
5. Código em inglês.
6. Seguir SOLID.
7. Seguir Clean Architecture.
8. Handler não pode acessar banco.
9. Handler não pode gerar JWT diretamente.
10. Toda lógica deve ser delegada para Use Cases e Interfaces.

Objetivos da implementação:

1. Endpoint POST /auth/login
2. Request DTO
3. Response DTO
4. Login Handler
5. Integração com AuthenticateUser Use Case
6. Integração com TokenProvider
7. Testes unitários do Handler

Estrutura esperada:

backend/
├── internal/
│
│   ├── http/
│   │   └── auth/
│   │       ├── handler.go
│   │       ├── handler_test.go
│   │       ├── request.go
│   │       ├── response.go
│   │       └── routes.go
│   │
│   ├── application/
│   │   └── auth/
│   │
│   ├── infrastructure/
│   │   └── auth/
│   │
│   └── domain/
│
└── cmd/
    └── api/

Endpoint esperado:

POST /auth/login

Request:

{
  "email": "john@example.com",
  "password": "12345678"
}

Response:

HTTP 200

{
  "access_token": "jwt-token",
  "token_type": "Bearer",
  "expires_in": 3600,
  "user": {
      "id": "uuid",
      "name": "John Doe",
      "email": "john@example.com"
  }
}

Headers:

Content-Type: application/json

Interfaces esperadas:

type AuthenticateUserUseCase interface {
    Execute(
        ctx context.Context,
        input AuthenticateUserInput,
    ) (
        AuthenticateUserOutput,
        error,
    )
}

type TokenProvider interface {
    Generate(
        userID string,
        email string,
    ) (string, error)
}

Fluxo esperado:

1. Receber email e senha
2. Validar request
3. Chamar AuthenticateUser
4. Receber usuário autenticado
5. Gerar JWT
6. Retornar token

Mapeamento de erros:

ErrInvalidCredentials

HTTP 401

{
  "error": "invalid credentials"
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

Validações HTTP:

- body obrigatório
- JSON válido
- email obrigatório
- password obrigatória

Testes obrigatórios:

Scenario:
Given valid credentials
When POST /auth/login
Then should return HTTP 200

Scenario:
Given valid credentials
When POST /auth/login
Then should return access token

Scenario:
Given invalid credentials
When POST /auth/login
Then should return HTTP 401

Scenario:
Given malformed JSON
When POST /auth/login
Then should return HTTP 400

Scenario:
Given missing email
When POST /auth/login
Then should return HTTP 400

Scenario:
Given missing password
When POST /auth/login
Then should return HTTP 400

Scenario:
Given token generation failure
When POST /auth/login
Then should return HTTP 500

Scenario:
Given valid request
When POST /auth/login
Then should invoke AuthenticateUser

Scenario:
Given authenticated user
When POST /auth/login
Then should invoke TokenProvider

Mocks esperados:

AuthenticateUserUseCaseMock

TokenProviderMock

Não utilizar:

- Banco de dados
- PostgreSQL
- Integração real

Todos os testes devem ser unitários.

Fluxo arquitetural esperado:

HTTP Request
      ↓
Login Handler
      ↓
AuthenticateUser Use Case
      ↓
TokenProvider
      ↓
HTTP Response

Critérios de aceite:

- Endpoint funcional.
- JWT retornado corretamente.
- Handler desacoplado.
- Tratamento de erros consistente.
- Testes unitários passando.

Entregáveis:

1. Handler completo.
2. DTOs Request/Response.
3. Rotas.
4. Testes unitários.
5. Exemplos de requests/responses.
6. Explicação rápida da arquitetura.

Não implementar:

- Middleware Auth
- Refresh Token
- Logout
- Sessões
- RBAC
- Multi-tenant
- OAuth
- Swagger

Somente endpoint HTTP de login.