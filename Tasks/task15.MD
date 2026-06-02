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
✓ Structured Logging
✓ Request ID Middleware
✓ Recovery Middleware

Objetivo desta task:

Implementar a listagem paginada de clientes.

IMPORTANTE:

Ainda NÃO implementar:

- Detalhes do cliente
- Atualização de clientes
- Exclusão de clientes
- Soft delete
- Processos jurídicos
- Multi-tenant

Esta task deve implementar apenas a consulta/listagem paginada de clientes.

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
10. Paginação deve ser implementada no Repository.

Objetivos da implementação:

1. ListClients Use Case
2. Endpoint GET /clients
3. Paginação
4. Filtros básicos
5. Ordenação
6. SQLC Queries
7. Testes unitários
8. Testes de integração do Repository

Estrutura esperada:

backend/
├── internal/
│
│   ├── application/
│   │   └── client/
│   │       ├── list_clients.go
│   │       ├── list_clients_test.go
│   │       ├── input.go
│   │       └── output.go
│   │
│   ├── http/
│   │   └── client/
│   │       ├── list_handler.go
│   │       └── list_handler_test.go
│   │
│   ├── domain/
│   │   └── client/
│   │
│   └── infrastructure/
│       └── persistence/
│
├── sql/
│   └── clients.sql
│
└── tests/

Endpoint esperado:

GET /clients

Proteção:

Authorization: Bearer <jwt>

Query Params suportados:

?page=1
&page_size=20

Opcionalmente:

?name=maria

?cpf=12345678900

?sort=name

?order=asc

Exemplo:

GET /clients?page=1&page_size=20

Response:

HTTP 200

{
  "data": [
    {
      "id": "uuid",
      "name": "Maria Silva",
      "cpf": "12345678900",
      "email": "maria@example.com",
      "phone": "81999999999"
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total": 85,
    "total_pages": 5
  }
}

Input esperado:

type ListClientsInput struct {
    Page     int
    PageSize int
    Name     string
    CPF      string
    Sort     string
    Order    string
}

Output esperado:

type ListClientsOutput struct {
    Data       []ClientDTO
    Page       int
    PageSize   int
    Total      int64
    TotalPages int
}

Repository Interface esperada:

type ClientRepository interface {

    List(
        ctx context.Context,
        params ListClientsParams,
    ) ([]Client, int64, error)
}

ListClientsParams:

type ListClientsParams struct {
    Offset   int
    Limit    int
    Name     string
    CPF      string
    Sort     string
    Order    string
}

Comportamentos obrigatórios:

Valores padrão:

page=1

page_size=20

Limites:

page_size máximo = 100

Ordenação permitida:

name

created_at

Order permitido:

asc

desc

Qualquer valor inválido:

→ utilizar padrão seguro

Filtros suportados:

Nome parcial

CPF exato

SQLC Queries obrigatórias:

ListClients

CountClients

Use Case responsabilidades:

1. Validar paginação
2. Calcular offset
3. Chamar repository
4. Calcular total_pages
5. Retornar resultado

Handler responsabilidades:

1. Ler query params
2. Montar input
3. Executar use case
4. Retornar response

Testes obrigatórios do Use Case:

Scenario:
Given clients exist
When listing clients
Then should return paginated result

Scenario:
Given page size greater than 100
When listing clients
Then should use maximum allowed

Scenario:
Given invalid page
When listing clients
Then should use default value

Scenario:
Given repository failure
When listing clients
Then should return error

Scenario:
Given filtered name
When listing clients
Then should return matching clients

Scenario:
Given filtered CPF
When listing clients
Then should return matching client

Testes obrigatórios do Handler:

Scenario:
Given authenticated request
When GET /clients
Then should return HTTP 200

Scenario:
Given invalid query params
When GET /clients
Then should use defaults

Scenario:
Given unauthorized request
When GET /clients
Then should return HTTP 401

Scenario:
Given internal error
When GET /clients
Then should return HTTP 500

Testes obrigatórios do Repository:

Utilizar Testcontainers.

Scenario:
Given persisted clients
When listing first page
Then should return paginated result

Scenario:
Given name filter
When listing
Then should filter correctly

Scenario:
Given CPF filter
When listing
Then should filter correctly

Scenario:
Given sorting by name
When listing
Then should return ordered data

Mocks esperados:

ClientRepositoryMock

ListClientsUseCaseMock

Não utilizar:

- Soft delete
- Cache
- Elasticsearch
- Full Text Search

Todos os testes devem ser unitários ou integração da persistência.

Critérios de aceite:

- GET /clients funcionando.
- Paginação funcionando.
- Filtros funcionando.
- Ordenação funcionando.
- JWT obrigatório.
- Testes passando.

Entregáveis:

1. ListClients Use Case.
2. Endpoint GET /clients.
3. Queries SQLC.
4. Repository atualizado.
5. Testes unitários.
6. Testes de integração.
7. Exemplos de requests/responses.
8. Explicação rápida da arquitetura.

Não implementar:

- GET /clients/{id}
- PUT /clients/{id}
- DELETE /clients/{id}
- Soft delete
- Processos
- Audiências

Somente listagem paginada de clientes.