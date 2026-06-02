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
✓ POST /clients
✓ GET /clients
✓ GET /clients/{id}
✓ Client Entity
✓ Client Repository
✓ CreateClient Use Case

Objetivo desta task:

Implementar a atualização de clientes.

IMPORTANTE:

Ainda NÃO implementar:

- Exclusão de clientes
- Soft delete
- Processos jurídicos
- Multi-tenant
- Histórico de alterações

Esta task deve implementar apenas a atualização dos dados cadastrais do cliente.

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
10. Regras de negócio ficam no Use Case e Domínio.

Objetivos da implementação:

1. UpdateClient Use Case
2. Endpoint PUT /clients/{id}
3. Atualização completa do recurso
4. Persistência PostgreSQL
5. Testes unitários
6. Testes de integração

Estrutura esperada:

backend/
├── internal/
│
│   ├── application/
│   │   └── client/
│   │       ├── update_client.go
│   │       ├── update_client_test.go
│   │       ├── input.go
│   │       ├── output.go
│   │       └── errors.go
│   │
│   ├── http/
│   │   └── client/
│   │       ├── update_handler.go
│   │       ├── update_handler_test.go
│   │       └── routes.go
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

Estratégia obrigatória:

PUT = atualização completa.

Todos os campos enviados substituem os valores atuais.

Endpoint esperado:

PUT /clients/{id}

Proteção:

Authorization: Bearer <jwt>

Request:

{
  "name": "Maria da Silva",
  "email": "maria.silva@example.com",
  "phone": "81988887777"
}

Observação:

CPF NÃO pode ser alterado.

CPF é imutável após criação.

Response:

HTTP 200

{
  "id": "uuid",
  "name": "Maria da Silva",
  "cpf": "12345678900",
  "email": "maria.silva@example.com",
  "phone": "81988887777"
}

Input esperado:

type UpdateClientInput struct {
    ID    string
    Name  string
    Email string
    Phone string
}

Output esperado:

type UpdateClientOutput struct {
    ID    string
    Name  string
    CPF   string
    Email string
    Phone string
}

Repository Interface esperada:

type ClientRepository interface {

    FindByID(
        ctx context.Context,
        id uuid.UUID,
    ) (*Client, error)

    Update(
        ctx context.Context,
        client *Client,
    ) error
}

SQLC Queries obrigatórias:

UpdateClient

GetClientByID

Erros esperados:

ErrClientNotFound

ErrInvalidClientID

ErrInvalidInput

Fluxo esperado:

1. Receber ID
2. Validar UUID
3. Buscar cliente existente
4. Atualizar dados permitidos
5. Persistir alterações
6. Retornar cliente atualizado

Campos atualizáveis:

Name

Email

Phone

Campos não atualizáveis:

ID

CPF

CreatedAt

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

Responsabilidade do Use Case:

1. Validar UUID
2. Buscar cliente
3. Atualizar dados permitidos
4. Executar validações da entidade
5. Persistir alterações

Responsabilidade do Handler:

1. Ler path param
2. Ler body
3. Executar Use Case
4. Retornar resposta HTTP

Testes obrigatórios do Use Case:

Scenario:
Given existing client
When updating client
Then should update successfully

Scenario:
Given invalid UUID
When updating client
Then should return ErrInvalidClientID

Scenario:
Given non existing client
When updating client
Then should return ErrClientNotFound

Scenario:
Given invalid email
When updating client
Then should return validation error

Scenario:
Given repository update failure
When updating client
Then should return error

Scenario:
Given CPF update attempt
When updating client
Then CPF must remain unchanged

Testes obrigatórios do Handler:

Scenario:
Given authenticated request
When PUT /clients/{id}
Then should return HTTP 200

Scenario:
Given invalid UUID
When PUT /clients/{id}
Then should return HTTP 400

Scenario:
Given client not found
When PUT /clients/{id}
Then should return HTTP 404

Scenario:
Given malformed JSON
When PUT /clients/{id}
Then should return HTTP 400

Scenario:
Given unauthorized request
When PUT /clients/{id}
Then should return HTTP 401

Scenario:
Given internal error
When PUT /clients/{id}
Then should return HTTP 500

Testes obrigatórios do Repository:

Utilizar Testcontainers.

Scenario:
Given persisted client
When updating client
Then changes should be saved

Scenario:
Given updated client
When loading again
Then should return updated values

Mocks esperados:

ClientRepositoryMock

UpdateClientUseCaseMock

Não utilizar:

- PATCH
- Soft delete
- Auditoria
- Event sourcing

Todos os testes devem ser unitários ou integração da persistência.

Critérios de aceite:

- PUT /clients/{id} funcionando.
- CPF protegido contra alteração.
- Atualização persistida corretamente.
- JWT obrigatório.
- Testes passando.

Entregáveis:

1. UpdateClient Use Case.
2. Endpoint PUT /clients/{id}.
3. Queries SQLC.
4. Repository atualizado.
5. Testes unitários.
6. Testes de integração.
7. Exemplos de requests/responses.
8. Explicação rápida da arquitetura.

Não implementar:

- DELETE /clients/{id}
- Soft delete
- Histórico
- Processos
- Audiências

Somente atualização de cliente.