Você é um engenheiro de software sênior especializado em Golang, APIs REST, Segurança, JWT, Clean Architecture, DDD e TDD.

Contexto do projeto:

Sistema SaaS para escritórios de advocacia.

Stack:

- Golang 1.26+
- Chi Router
- PostgreSQL
- SQLC
- Goose
- JWT
- Testify
- Clean Architecture
- Modular Monolith

Estado atual do projeto:

Já existem:

✓ User Entity
✓ User Repository
✓ PostgreSQL Persistence
✓ CreateUser Use Case
✓ AuthenticateUser Use Case
✓ JWT Token Provider
✓ POST /users
✓ POST /auth/login
✓ JWT Authentication Middleware
✓ Structured Logging
✓ Request ID Middleware
✓ Recovery Middleware

Objetivo desta task:

Implementar o primeiro endpoint protegido da aplicação: GET /me.

IMPORTANTE:

Ainda NÃO implementar:

- Refresh Token
- Logout
- RBAC
- Multi-tenant
- Perfis
- Permissões
- OAuth

Esta task deve validar o fluxo completo de autenticação da aplicação.

Regras obrigatórias:

1. Trabalhar usando TDD.
2. Criar testes antes da implementação.
3. Implementar apenas o necessário.
4. Não antecipar funcionalidades futuras.
5. Código em inglês.
6. Seguir SOLID.
7. Seguir DDD.
8. Seguir Clean Architecture.
9. Handler não pode acessar banco diretamente.
10. Handler deve depender apenas de um Use Case.

Objetivos da implementação:

1. GetCurrentUser Use Case
2. Endpoint GET /me
3. Integração com JWT Middleware
4. Leitura das Claims do Context
5. Busca do usuário autenticado
6. Testes unitários
7. Testes HTTP

Estrutura esperada:

backend/
├── internal/
│
│   ├── application/
│   │   └── auth/
│   │       ├── get_current_user.go
│   │       ├── get_current_user_test.go
│   │       ├── input.go
│   │       └── output.go
│   │
│   ├── http/
│   │   └── auth/
│   │       ├── me_handler.go
│   │       ├── me_handler_test.go
│   │       └── routes.go
│   │
│   ├── middleware/
│   │
│   ├── domain/
│   │   └── user/
│   │
│   └── infrastructure/
│
└── tests/

Endpoint esperado:

GET /me

Headers:

Authorization: Bearer <jwt>

Response:

HTTP 200

{
  "id": "uuid",
  "name": "John Doe",
  "email": "john@example.com"
}

Fluxo esperado:

Request
   ↓
Auth Middleware
   ↓
Claims Context
   ↓
GetCurrentUser Use Case
   ↓
UserRepository
   ↓
Response

Interface esperada:

type UserRepository interface {
    FindByID(
        ctx context.Context,
        id string,
    ) (*User, error)
}

Input esperado:

type GetCurrentUserInput struct {
    UserID string
}

Output esperado:

type GetCurrentUserOutput struct {
    ID    string
    Name  string
    Email string
}

Responsabilidade do Use Case:

1. Receber UserID
2. Buscar usuário
3. Retornar dados públicos

Responsabilidade do Handler:

1. Ler UserID do Context
2. Executar Use Case
3. Retornar resposta HTTP

Mapeamento de erros:

Usuário não encontrado

HTTP 404

{
  "error": "user not found"
}

Não autenticado

HTTP 401

{
  "error": "unauthorized"
}

Erro interno

HTTP 500

{
  "error": "internal server error"
}

Testes obrigatórios do Use Case:

Scenario:
Given existing user
When loading current user
Then should return user

Scenario:
Given non-existing user
When loading current user
Then should return not found

Scenario:
Given repository failure
When loading current user
Then should return error

Testes obrigatórios do Handler:

Scenario:
Given authenticated request
When GET /me
Then should return HTTP 200

Scenario:
Given authenticated request
When GET /me
Then should return user information

Scenario:
Given missing claims
When GET /me
Then should return HTTP 401

Scenario:
Given user not found
When GET /me
Then should return HTTP 404

Scenario:
Given internal error
When GET /me
Then should return HTTP 500

Mocks esperados:

UserRepositoryMock

GetCurrentUserUseCaseMock

Não utilizar:

- Banco real
- PostgreSQL
- Integrações externas

Todos os testes devem ser unitários.

Exemplo de integração:

router.Group(func(r chi.Router) {

    r.Use(AuthMiddleware)

    r.Get("/me", authHandler.Me)
})

Critérios de aceite:

- Endpoint protegido funcionando.
- Claims lidas corretamente.
- Use Case desacoplado.
- Handler desacoplado.
- Testes passando.
- Fluxo de autenticação validado ponta a ponta.

Entregáveis:

1. GetCurrentUser Use Case.
2. Handler GET /me.
3. DTOs Input/Output.
4. Testes unitários.
5. Testes HTTP.
6. Exemplos de request/response.
7. Explicação rápida da arquitetura.

Não implementar:

- Refresh Token
- Logout
- RBAC
- Roles
- Permissions
- Tenant
- OAuth
- Perfil de usuário editável

Somente endpoint GET /me.