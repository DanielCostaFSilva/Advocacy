Você é um engenheiro de software sênior especializado em Golang, Segurança, JWT, Clean Architecture, DDD, Middleware HTTP e TDD.

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
✓ CreateUser Use Case
✓ AuthenticateUser Use Case
✓ JWT Token Provider
✓ POST /users
✓ POST /auth/login
✓ Structured Logging
✓ Request ID Middleware
✓ Recovery Middleware

Objetivo desta task:

Implementar o middleware responsável pela autenticação JWT das rotas protegidas.

IMPORTANTE:

Ainda NÃO implementar:

- Refresh Token
- Logout
- RBAC
- Multi-tenant
- Permissões
- OAuth
- Sessões

Esta task deve apenas validar o Access Token e disponibilizar o usuário autenticado no Context da requisição.

Regras obrigatórias:

1. Trabalhar usando TDD.
2. Criar testes antes da implementação.
3. Implementar apenas o necessário.
4. Não antecipar funcionalidades futuras.
5. Código em inglês.
6. Seguir SOLID.
7. Seguir Clean Architecture.
8. Middleware não pode acessar banco de dados.
9. Middleware deve depender apenas do TokenProvider.
10. Claims devem ser disponibilizadas através do Context.

Objetivos da implementação:

1. JWT Authentication Middleware
2. Context User Claims
3. Helpers para acesso ao usuário autenticado
4. Testes unitários
5. Integração com Chi Router

Estrutura esperada:

backend/
├── internal/
│
│   ├── middleware/
│   │   ├── auth.go
│   │   ├── auth_test.go
│   │   ├── context.go
│   │   └── context_test.go
│   │
│   ├── infrastructure/
│   │   └── auth/
│   │
│   └── shared/
│
└── tests/

Interface esperada:

type TokenProvider interface {
    Validate(
        token string,
    ) (*TokenClaims, error)
}

Claims esperadas:

type TokenClaims struct {
    UserID string
    Email  string
    IssuedAt time.Time
    ExpiresAt time.Time
}

Header esperado:

Authorization

Formato esperado:

Authorization: Bearer <jwt>

Fluxo esperado:

1. Ler header Authorization
2. Validar formato Bearer
3. Extrair token
4. Validar token
5. Extrair claims
6. Inserir claims no Context
7. Continuar fluxo

Context esperado:

ctx
  ↓
User Claims
  ↓
Handler

Helpers esperados:

GetUserID(ctx)

GetUserEmail(ctx)

GetClaims(ctx)

Comportamentos obrigatórios:

Header ausente

→ HTTP 401

Token inválido

→ HTTP 401

Token expirado

→ HTTP 401

Formato inválido

→ HTTP 401

Token válido

→ permitir acesso

Resposta esperada:

HTTP 401

{
  "error": "unauthorized"
}

Não retornar detalhes internos.

Não revelar motivo exato da falha.

Mapeamento de erros:

ErrInvalidToken
ErrExpiredToken

↓

HTTP 401

Testes obrigatórios:

Scenario:
Given valid bearer token
When middleware executes
Then request should continue

Scenario:
Given missing authorization header
When middleware executes
Then should return HTTP 401

Scenario:
Given malformed authorization header
When middleware executes
Then should return HTTP 401

Scenario:
Given invalid token
When middleware executes
Then should return HTTP 401

Scenario:
Given expired token
When middleware executes
Then should return HTTP 401

Scenario:
Given valid token
When middleware executes
Then claims should be available in context

Scenario:
Given valid token
When middleware executes
Then user id should be available

Scenario:
Given valid token
When middleware executes
Then user email should be available

Mocks esperados:

TokenProviderMock

Não utilizar:

- Banco de dados
- Repository
- PostgreSQL

Todos os testes devem ser unitários.

Exemplo de uso esperado:

router.Group(func(r chi.Router) {

    r.Use(AuthMiddleware)

    r.Get("/me", handler.Me)
})

Exemplo de acesso:

userID := middleware.GetUserID(ctx)

email := middleware.GetUserEmail(ctx)

Critérios de aceite:

- Middleware validando JWT.
- Claims disponíveis no Context.
- Rotas protegidas funcionando.
- Tratamento consistente de erro.
- Testes unitários passando.

Entregáveis:

1. Middleware JWT.
2. Helpers de Context.
3. Testes unitários.
4. Exemplo de integração com Chi.
5. Explicação rápida da arquitetura.

Não implementar:

- Refresh Token
- Logout
- RBAC
- Roles
- Permissions
- Tenant
- OAuth
- Sessões

Somente middleware de autenticação JWT.