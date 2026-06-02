Você é um engenheiro de software sênior especializado em Golang, DDD, TDD, Clean Architecture e sistemas SaaS.

Contexto do projeto:

Sistema SaaS para escritórios de advocacia.

Stack:

- Golang 1.26+
- PostgreSQL
- SQLC
- Goose
- Chi Router
- Testify
- Docker
- Clean Architecture
- Modular Monolith

Objetivo desta task:

Implementar o módulo de Usuários no nível de domínio e persistência.

IMPORTANTE:

Ainda NÃO implementar:

- Login
- JWT
- Refresh Token
- Middleware de autenticação
- RBAC
- Multi-tenant

Esta task existe apenas para criar a fundação do módulo de usuários.

Regras obrigatórias:

1. Trabalhar usando TDD.
2. Criar testes antes da implementação.
3. Implementar apenas o necessário.
4. Não antecipar funcionalidades futuras.
5. Código em inglês.
6. Seguir SOLID.
7. Seguir DDD.
8. Seguir Clean Architecture.
9. Regras de negócio devem ficar no domínio.
10. Repositórios não devem conter regras de negócio.

Objetivos da implementação:

1. Entidade User
2. Regras de domínio
3. User Repository
4. Persistência PostgreSQL
5. Migration da tabela users
6. Testes de domínio
7. Testes de integração do repository

Estrutura esperada:

backend/
├── internal/
│
│   ├── user/
│   │
│   ├── domain/
│   │   └── user/
│   │       ├── entity.go
│   │       ├── errors.go
│   │       ├── repository.go
│   │       ├── entity_test.go
│   │       └── repository_mock.go
│   │
│   ├── infrastructure/
│   │   └── persistence/
│   │       └── postgres/
│   │           └── user_repository.go
│   │
│   └── database/
│
├── migrations/
│   └── 000002_create_users.sql
│
├── sql/
│   └── users.sql
│
├── sqlc.yaml
│
└── tests/

Entidade User:

type User struct {
    ID
    Name
    Email
    PasswordHash
    CreatedAt
    UpdatedAt
}

Regras de domínio:

Nome

- obrigatório
- mínimo 3 caracteres

Email

- obrigatório
- formato válido

PasswordHash

- obrigatório
- não pode ser vazio

ID

- UUID

Campos auditáveis

- CreatedAt
- UpdatedAt

Factory esperada:

NewUser(
    name,
    email,
    passwordHash,
)

Erros esperados:

ErrInvalidName
ErrInvalidEmail
ErrInvalidPasswordHash

Migration esperada:

Tabela:

users

Campos:

id uuid pk
name varchar
email varchar unique
password_hash varchar
created_at timestamp
updated_at timestamp

Índices:

pk_users
uk_users_email

Repository Interface:

type UserRepository interface {
    Create(ctx context.Context, user *User) error

    FindByID(
        ctx context.Context,
        id uuid.UUID,
    ) (*User, error)

    FindByEmail(
        ctx context.Context,
        email string,
    ) (*User, error)
}

Persistência:

- PostgreSQL
- SQLC
- Context obrigatório
- Sem ORM

Testes de domínio obrigatórios:

Scenario:
Given valid user data
When creating user
Then user should be created

Scenario:
Given invalid email
When creating user
Then should return error

Scenario:
Given empty password hash
When creating user
Then should return error

Scenario:
Given short name
When creating user
Then should return error

Testes de repository obrigatórios:

Scenario:
Given valid user
When persisting user
Then should save successfully

Scenario:
Given existing email
When persisting another user
Then should return duplicate error

Scenario:
Given existing user
When searching by email
Then should return user

Scenario:
Given existing user
When searching by id
Then should return user

Testes de integração:

Utilizar Testcontainers.

Não utilizar mocks para persistência.

Comandos esperados:

make migrate-up
make migrate-down
make test

Entregáveis:

1. Migration completa.
2. Entidade User.
3. Regras de domínio.
4. Repository interface.
5. Implementação PostgreSQL.
6. Queries SQLC.
7. Testes de domínio.
8. Testes de integração.
9. Explicação rápida da arquitetura.

Não implementar:

- Endpoints HTTP
- Login
- JWT
- Hash de senha
- Cadastro via API
- Casos de uso
- Services

Somente domínio e persistência.