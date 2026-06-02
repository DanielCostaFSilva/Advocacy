Você é um engenheiro de software sênior especializado em Golang, Segurança, JWT, DDD, Clean Architecture e TDD.

Contexto do projeto:

Sistema SaaS para escritórios de advocacia.

Stack:

- Golang 1.26+
- PostgreSQL
- SQLC
- Goose
- Chi Router
- Testify
- Clean Architecture
- Modular Monolith

Estado atual do projeto:

Já existem:

✓ User Entity
✓ User Repository
✓ CreateUser Use Case
✓ AuthenticateUser Use Case
✓ Password Hasher
✓ Password Verifier
✓ Structured Logging
✓ Request ID Middleware

Objetivo desta task:

Implementar a infraestrutura responsável pela geração e validação de JWT.

IMPORTANTE:

Ainda NÃO implementar:

- Endpoint de Login
- Middleware de Autenticação
- Refresh Token
- Logout
- Sessões
- RBAC
- Multi-tenant

Esta task deve implementar apenas o mecanismo de geração e validação de Access Tokens.

Regras obrigatórias:

1. Trabalhar usando TDD.
2. Criar testes antes da implementação.
3. Implementar apenas o necessário.
4. Não antecipar funcionalidades futuras.
5. Código em inglês.
6. Seguir SOLID.
7. Seguir Clean Architecture.
8. Nenhuma dependência direta do JWT dentro dos Use Cases.
9. JWT deve ficar isolado na infraestrutura.

Objetivos da implementação:

1. JWT Provider Interface
2. JWT Service
3. Token Generation
4. Token Validation
5. Claims Model
6. Configuração de Expiração
7. Testes unitários

Estrutura esperada:

backend/
├── internal/
│
│   ├── application/
│   │
│   ├── domain/
│   │
│   ├── infrastructure/
│   │   └── auth/
│   │       ├── jwt_service.go
│   │       ├── jwt_service_test.go
│   │       ├── claims.go
│   │       └── errors.go
│   │
│   ├── config/
│   │
│   └── shared/
│
└── tests/

Interface esperada:

type TokenProvider interface {
    Generate(
        userID string,
        email string,
    ) (string, error)

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

Claims JWT obrigatórias:

sub → user id

email → user email

iat → issued at

exp → expiration

Configuração esperada:

JWT_SECRET

JWT_EXPIRATION_MINUTES

Exemplo:

JWT_SECRET=my-secret-key

JWT_EXPIRATION_MINUTES=60

Algoritmo obrigatório:

HS256

Fluxo Generate:

1. Receber user id
2. Receber email
3. Gerar claims
4. Assinar token
5. Retornar JWT

Fluxo Validate:

1. Receber token
2. Validar assinatura
3. Validar expiração
4. Extrair claims
5. Retornar TokenClaims

Erros esperados:

ErrInvalidToken

ErrExpiredToken

ErrGenerateToken

ErrMissingSecret

Comportamentos obrigatórios:

Token válido

→ retornar claims

Token expirado

→ retornar ErrExpiredToken

Token inválido

→ retornar ErrInvalidToken

Assinatura inválida

→ retornar ErrInvalidToken

Configuração ausente

→ retornar ErrMissingSecret

Testes obrigatórios:

Scenario:
Given valid user data
When generating token
Then should return valid JWT

Scenario:
Given valid token
When validating token
Then should return claims

Scenario:
Given expired token
When validating token
Then should return ErrExpiredToken

Scenario:
Given malformed token
When validating token
Then should return ErrInvalidToken

Scenario:
Given invalid signature
When validating token
Then should return ErrInvalidToken

Scenario:
Given missing secret
When generating token
Then should return ErrMissingSecret

Scenario:
Given generated token
When validating token
Then user id should match

Scenario:
Given generated token
When validating token
Then email should match

Não utilizar:

- Banco de dados
- HTTP
- Middleware
- Context

Todos os testes devem ser unitários.

Dependência recomendada:

github.com/golang-jwt/jwt/v5

Critérios de aceite:

- Tokens gerados corretamente.
- Tokens validados corretamente.
- Expiração funcionando.
- Assinatura funcionando.
- Erros padronizados.
- Testes unitários passando.

Entregáveis:

1. JWT Provider Interface.
2. JWT Service.
3. Claims Model.
4. Errors.
5. Configuração necessária.
6. Testes unitários.
7. Exemplos de uso.

Não implementar:

- Endpoint Login
- Middleware Auth
- Refresh Token
- Logout
- Sessões
- RBAC
- Multi-tenant

Somente geração e validação de Access Tokens JWT.