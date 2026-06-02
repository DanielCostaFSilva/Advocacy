Você é um engenheiro de software sênior especializado em Golang, APIs REST, DDD, Clean Architecture, PostgreSQL e TDD.

Contexto do projeto:

Sistema SaaS para escritórios de advocacia.

Stack:

- Golang 1.26+
- PostgreSQL
- SQLC
- Goose
- Chi Router
- JWT
- Testify
- Testcontainers
- Clean Architecture
- Modular Monolith

Estado atual do projeto:

Já existem:

✓ JWT Authentication Middleware
✓ GET /me
✓ Client Entity
✓ Client Repository
✓ CreateClient Use Case
✓ POST /clients
✓ GET /clients (paginado)
✓ Structured Logging
✓ Request ID Middleware
✓ Recovery Middleware

Objetivo desta task:

Implementar a consulta de um cliente específico por ID.

IMPORTANTE:

Ainda NÃO implementar:

- Atualização de clientes
- Exclusão de clientes
- Soft delete
- Processos jurídicos
- Multi-tenant

Esta task deve implementar apenas a recuperação de um cliente pelo seu identificador.

Regras obrigatórias:

1. Trabalhar usando TDD.
2. Criar testes antes da implementação.
3. Implementar apenas o necessário.
4. Não antecipar funcionalidades futuras.
5. Código em inglês.
6. Seguir SOLID.
7. Seguir DDD.
8. Seguir Clean Architecture.
9. Handler não acessa banco diretamente.
10. Toda persistência deve ocorrer via Repository.

Objetivos da implementação:

1. GetClientByID Use Case
2. Endpoint GET /clients/{id}
3. Tratamento de cliente não encontrado
4. Testes unitários
5. Testes de integração do Repository

Estrutura esperada:

backend/
├── internal/
│
│   ├── application/
│   │   └── client/
│   │       ├── get_client_by_id.go
│   │       ├── get_client_by_id_test.go
│   │       ├── input.go
│   │       ├── output.go
│   │       └── errors.go
│   │
│   ├── http/
│   │   └── client/
│   │       ├── get_handler.go
│   │       ├── get_handler_test.go
│   │       └── routes.go
│   │
│   ├── domain/
│   │   └── client/
│   │
│   └── infrastructure/
│       └── persistence/
│
└── tests/

Endpoint esperado:

GET /clients/{id}

Proteção:

Authorization: Bearer <jwt>

Exemplo:

GET /clients/550e8400-e29b-41d4-a716-446655440000

Response:

HTTP 200

{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Maria Silva",
  "cpf": "12345678900",
  "email": "maria@example.com",
  "phone": "81999999999",
  "created_at": "2026-06-02T10:00:00Z",
  "updated_at": "2026-06-02T10:00:00Z"
}

Input esperado:

type GetClientByIDInput struct {
    ID string
}

Output esperado:

type GetClientByIDOutput struct {
    ID        string
    Name      string
    CPF       string
    Email     string
    Phone     string
    CreatedAt time.Time
    UpdatedAt time.Time
}

Repository Interface:

type ClientRepository interface {

    FindByID(
        ctx context.Context,
        id uuid.UUID,
    ) (*Client, error)
}

Erros esperados:

ErrClientNotFound

ErrInvalidClientID

Fluxo esperado:

1. Receber ID
2. Validar UUID
3. Buscar cliente
4. Retornar dados

Mapeamento HTTP:

ErrClientNotFound

↓

HTTP 404

{
  "error": "client not found"
}

ErrInvalidClientID

↓

HTTP 400

{
  "error": "invalid client id"
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

Responsabilidade do Use Case:

1. Validar UUID
2. Buscar cliente no Repository
3. Retornar DTO

Responsabilidade do Handler:

1. Ler parâmetro da rota
2. Chamar Use Case
3. Mapear resposta HTTP

Testes obrigatórios do Use Case:

Scenario:
Given existing client
When loading client by id
Then should return client

Scenario:
Given invalid UUID
When loading client by id
Then should return ErrInvalidClientID

Scenario:
Given non existing client
When loading client by id
Then should return ErrClientNotFound

Scenario:
Given repository failure
When loading client by id
Then should return error

Testes obrigatórios do Handler:

Scenario:
Given authenticated request
When GET /clients/{id}
Then should return HTTP 200

Scenario:
Given invalid UUID
When GET /clients/{id}
Then should return HTTP 400

Scenario:
Given client not found
When GET /clients/{id}
Then should return HTTP 404

Scenario:
Given unauthorized request
When GET /clients/{id}
Then should return HTTP 401

Scenario:
Given internal error
When GET /clients/{id}
Then should return HTTP 500

Testes obrigatórios do Repository:

Utilizar Testcontainers.

Scenario:
Given persisted client
When finding by id
Then should return client

Scenario:
Given unknown id
When finding by id
Then should return nil

Mocks esperados:

ClientRepositoryMock

GetClientByIDUseCaseMock

Não utilizar:

- Cache
- Elasticsearch
- Soft delete
- Banco fake

Todos os testes devem ser unitários ou integração da persistência.

Critérios de aceite:

- GET /clients/{id} funcionando.
- UUID validado.
- Cliente retornado corretamente.
- Tratamento de 404 funcionando.
- JWT obrigatório.
- Testes passando.

Entregáveis:

1. GetClientByID Use Case.
2. Endpoint GET /clients/{id}.
3. Testes unitários.
4. Testes de integração.
5. Exemplos de requests/responses.
6. Explicação rápida da arquitetura.

Não implementar:

- PUT /clients/{id}
- DELETE /clients/{id}
- Soft delete
- Processos
- Audiências
- Uploads

Somente consulta de cliente por ID.