Você é um engenheiro de software sênior especializado em Golang, DDD, Clean Architecture, PostgreSQL e TDD.

Contexto do projeto:

Sistema SaaS para escritórios de advocacia.

Stack:

- Golang 1.26+
- PostgreSQL
- SQLC
- Goose
- Chi Router
- Testify
- Testcontainers
- Clean Architecture
- Modular Monolith

Estado atual do projeto:

Já existem:

✓ Autenticação completa
✓ User Entity
✓ User Repository
✓ JWT Authentication
✓ GET /me
✓ Observabilidade básica
✓ PostgreSQL configurado
✓ SQLC configurado

Objetivo desta task:

Criar a fundação do módulo de Clientes.

IMPORTANTE:

Ainda NÃO implementar:

- Endpoints HTTP
- Casos de uso
- Upload de documentos
- Processos jurídicos
- Multi-tenant
- Soft delete

Esta task deve implementar apenas:

- Entidade Client
- Regras de domínio
- Migration
- Repository
- Persistência PostgreSQL

Regras obrigatórias:

1. Trabalhar usando TDD.
2. Criar testes antes da implementação.
3. Implementar apenas o necessário.
4. Não antecipar funcionalidades futuras.
5. Código em inglês.
6. Seguir SOLID.
7. Seguir DDD.
8. Seguir Clean Architecture.
9. Regras de negócio ficam no domínio.
10. Repository não contém regras de negócio.

Objetivos da implementação:

1. Client Entity
2. Domain Validations
3. Migration clients
4. SQLC Queries
5. Repository Interface
6. PostgreSQL Repository
7. Testes de domínio
8. Testes de integração

Estrutura esperada:

backend/
├── internal/
│
│   ├── domain/
│   │   └── client/
│   │       ├── entity.go
│   │       ├── errors.go
│   │       ├── repository.go
│   │       ├── entity_test.go
│   │       └── repository_mock.go
│   │
│   ├── infrastructure/
│   │   └── persistence/
│   │       └── postgres/
│   │           └── client_repository.go
│   │
│   └── database/
│
├── migrations/
│   └── 000003_create_clients.sql
│
├── sql/
│   └── clients.sql
│
└── tests/

Entidade esperada:

type Client struct {
    ID          uuid.UUID
    Name        string
    CPF         string
    Email       string
    Phone       string
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

Factory esperada:

NewClient(
    name string,
    cpf string,
    email string,
    phone string,
) (*Client, error)

Validações obrigatórias:

Name

- obrigatório
- mínimo 3 caracteres

CPF

- obrigatório
- deve possuir 11 dígitos
- armazenar apenas números

Email

- opcional
- se informado deve possuir formato válido

Phone

- opcional

Erros esperados:

ErrInvalidName

ErrInvalidCPF

ErrInvalidEmail

Migration esperada:

Tabela:

clients

Campos:

id uuid primary key

name varchar(255)

cpf varchar(11)

email varchar(255)

phone varchar(30)

created_at timestamp

updated_at timestamp

Restrições:

cpf unique

Índices:

pk_clients

uk_clients_cpf

Repository Interface:

type ClientRepository interface {

    Create(
        ctx context.Context,
        client *Client,
    ) error

    FindByID(
        ctx context.Context,
        id uuid.UUID,
    ) (*Client, error)

    FindByCPF(
        ctx context.Context,
        cpf string,
    ) (*Client, error)
}

Persistência:

- PostgreSQL
- SQLC
- Context obrigatório
- Sem ORM

Queries SQLC obrigatórias:

CreateClient

GetClientByID

GetClientByCPF

Testes de domínio obrigatórios:

Scenario:
Given valid client data
When creating client
Then should create successfully

Scenario:
Given invalid CPF
When creating client
Then should return ErrInvalidCPF

Scenario:
Given invalid email
When creating client
Then should return ErrInvalidEmail

Scenario:
Given short name
When creating client
Then should return ErrInvalidName

Scenario:
Given CPF with punctuation
When creating client
Then should normalize CPF

Exemplo:

Input:

123.456.789-00

Output:

12345678900

Testes de integração obrigatórios:

Utilizar Testcontainers.

Scenario:
Given valid client
When persisting
Then should save successfully

Scenario:
Given duplicated CPF
When persisting
Then should return duplicate error

Scenario:
Given existing client
When searching by CPF
Then should return client

Scenario:
Given existing client
When searching by ID
Then should return client

Mocks esperados:

ClientRepositoryMock

Não utilizar:

- HTTP
- Use Cases
- JWT
- Middleware

Todos os testes devem ser unitários ou integração de persistência.

Critérios de aceite:

- Entidade Client criada.
- Regras de domínio implementadas.
- Migration criada.
- SQLC funcionando.
- Repository funcionando.
- Testes de domínio passando.
- Testes de integração passando.

Entregáveis:

1. Migration completa.
2. Entidade Client.
3. Regras de domínio.
4. Repository Interface.
5. Implementação PostgreSQL.
6. Queries SQLC.
7. Testes de domínio.
8. Testes de integração.
9. Explicação rápida da arquitetura.

Não implementar:

- CreateClient Use Case
- Endpoints HTTP
- Paginação
- Atualização
- Exclusão
- Processos
- Audiências

Somente domínio e persistência do módulo Client.