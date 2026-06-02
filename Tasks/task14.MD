Você é um engenheiro de software sênior especializado em Golang, APIs REST, DDD, Clean Architecture, Segurança e TDD.

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

✓ JWT Authentication Middleware
✓ GET /me
✓ Client Entity
✓ Client Repository
✓ CreateClient Use Case
✓ Structured Logging
✓ Request ID Middleware
✓ Recovery Middleware

Objetivo desta task:

Expor o CreateClient Use Case através de uma API REST protegida.

IMPORTANTE:

Ainda NÃO implementar:

- Listagem de clientes
- Atualização de clientes
- Exclusão de clientes
- Paginação
- Busca avançada
- Processos jurídicos
- Multi-tenant

Esta task deve implementar apenas o endpoint de criação de clientes.

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
10. Endpoint deve ser protegido por JWT.

Objetivos da implementação:

1. Endpoint POST /clients
2. Request DTO
3. Response DTO
4. Client Handler
5. Registro de rota protegida
6. Testes unitários do Handler
7. Integração com JWT Middleware

Estrutura esperada:

backend/
├── internal/
│
│   ├── http/
│   │   └── client/
│   │       ├── handler.go
│   │       ├── handler_test.go
│   │       ├── request.go
│   │       ├── response.go
│   │       └── routes.go
│   │
│   ├── application/
│   │   └── client/
│   │
│   ├── middleware/
│   │
│   └── domain/
│
└── tests/

Endpoint esperado:

POST /clients

Proteção:

Authorization: Bearer <jwt>

Request:

{
  "name": "Maria Silva",
  "cpf": "12345678900",
  "email": "maria@example.com",
  "phone": "81999999999"
}

Response:

HTTP 201

{
  "id": "uuid",
  "name": "Maria Silva",
  "cpf": "12345678900",
  "email": "maria@example.com",
  "phone": "81999999999"
}

Headers:

Content-Type: application/json

Interface esperada:

type CreateClientUseCase interface {
    Execute(
        ctx context.Context,
        input CreateClientInput,
    ) (
        CreateClientOutput,
        error,
    )
}

Responsabilidade do Handler:

1. Ler request JSON
2. Validar estrutura mínima
3. Executar CreateClient Use Case
4. Retornar resposta HTTP

Responsabilidade do Middleware:

1. Garantir usuário autenticado
2. Bloquear acesso sem JWT válido

Mapeamento de erros:

ErrClientAlreadyExists

↓

HTTP 409

{
  "error": "client already exists"
}

ErrInvalidInput

↓

HTTP 400

{
  "error": "invalid input"
}

Unauthorized

↓

HTTP 401

{
  "error": "unauthorized"
}

Erro interno

↓

HTTP 500

{
  "error": "internal server error"
}

Validações HTTP:

Body obrigatório

JSON válido

Name obrigatório

CPF obrigatório

As regras de negócio continuam sendo responsabilidade do domínio e do Use Case.

Testes obrigatórios:

Scenario:
Given authenticated request
When POST /clients
Then should return HTTP 201

Scenario:
Given authenticated request
When POST /clients
Then should return created client

Scenario:
Given duplicated CPF
When POST /clients
Then should return HTTP 409

Scenario:
Given invalid request body
When POST /clients
Then should return HTTP 400

Scenario:
Given malformed JSON
When POST /clients
Then should return HTTP 400

Scenario:
Given missing JWT
When POST /clients
Then should return HTTP 401

Scenario:
Given invalid JWT
When POST /clients
Then should return HTTP 401

Scenario:
Given internal error
When POST /clients
Then should return HTTP 500

Scenario:
Given valid request
When POST /clients
Then should invoke CreateClient Use Case

Mocks esperados:

CreateClientUseCaseMock

TokenProviderMock (quando necessário para testes da rota protegida)

Não utilizar:

- Banco de dados
- PostgreSQL
- Integrações reais

Todos os testes devem ser unitários.

Fluxo arquitetural esperado:

HTTP Request
      ↓
JWT Middleware
      ↓
Client Handler
      ↓
CreateClient Use Case
      ↓
Repository

Critérios de aceite:

- Endpoint protegido por JWT.
- Cliente criado com sucesso.
- Tratamento de erros consistente.
- Handler desacoplado da persistência.
- Testes unitários passando.

Entregáveis:

1. Handler completo.
2. DTOs Request/Response.
3. Rotas protegidas.
4. Testes unitários.
5. Exemplos de requests/responses.
6. Explicação rápida da arquitetura.

Não implementar:

- GET /clients
- GET /clients/{id}
- PUT /clients/{id}
- DELETE /clients/{id}
- Processos
- Audiências
- Uploads
- Paginação

Somente endpoint POST /clients.