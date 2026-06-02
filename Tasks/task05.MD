Você é um engenheiro de software sênior especializado em Golang, DDD, Clean Architecture, TDD e sistemas SaaS.

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

Estado atual do projeto:

Já existem:

✓ Entidade User
✓ Regras de domínio
✓ Migration users
✓ UserRepository
✓ Implementação PostgreSQL
✓ SQLC
✓ Testes de domínio
✓ Testes de integração

Objetivo desta task:

Implementar o caso de uso responsável pela criação de usuários.

IMPORTANTE:

Ainda NÃO implementar:

- Login
- JWT
- Middleware Auth
- Endpoints HTTP
- Controllers
- RBAC
- Multi-tenant

Esta task deve implementar apenas a camada de aplicação (Use Case).

Regras obrigatórias:

1. Trabalhar usando TDD.
2. Criar testes antes da implementação.
3. Implementar apenas o necessário.
4. Não antecipar funcionalidades futuras.
5. Código em inglês.
6. Seguir SOLID.
7. Seguir DDD.
8. Seguir Clean Architecture.
9. Não acessar banco diretamente no Use Case.
10. Todo acesso a dados deve ocorrer via interfaces.

Objetivos da implementação:

1. CreateUser Use Case
2. Password Hasher Interface
3. Verificação de e-mail duplicado
4. Hash da senha
5. Testes unitários

Estrutura esperada:

backend/
├── internal/
│
│   ├── application/
│   │   └── user/
│   │       ├── create_user.go
│   │       ├── create_user_test.go
│   │       ├── input.go
│   │       └── output.go
│   │
│   ├── domain/
│   │   └── user/
│   │
│   └── infrastructure/
│
└── tests/

Input esperado:

type CreateUserInput struct {
    Name     string
    Email    string
    Password string
}

Output esperado:

type CreateUserOutput struct {
    ID    string
    Name  string
    Email string
}

Dependências do Use Case:

type UserRepository interface {
    Create(...)
    FindByEmail(...)
}

type PasswordHasher interface {
    Hash(password string) (string, error)
}

Fluxo esperado:

1. Receber dados
2. Validar entrada
3. Verificar se email já existe
4. Gerar hash da senha
5. Criar entidade User
6. Persistir usuário
7. Retornar resultado

Erros esperados:

ErrEmailAlreadyExists
ErrInvalidInput
ErrHashPassword

Validações:

Nome

- obrigatório

Email

- obrigatório

Senha

- obrigatória
- mínimo 8 caracteres

Não validar regras de domínio duplicadas.
O domínio continua responsável pelas validações da entidade User.

Testes obrigatórios:

Scenario:
Given valid input
When creating user
Then user should be created

Scenario:
Given duplicated email
When creating user
Then should return ErrEmailAlreadyExists

Scenario:
Given hash failure
When creating user
Then should return ErrHashPassword

Scenario:
Given invalid password
When creating user
Then should return validation error

Scenario:
Given repository create failure
When creating user
Then should return error

Mocks esperados:

UserRepositoryMock

PasswordHasherMock

Não utilizar banco de dados nos testes desta task.

Todos os testes devem ser unitários.

Exemplo de fluxo:

input
 ↓
Use Case
 ↓
Repository Interface
 ↓
PasswordHasher Interface

Entregáveis:

1. CreateUser Use Case.
2. DTOs Input/Output.
3. Interfaces necessárias.
4. Mocks.
5. Testes unitários.
6. Explicação rápida da arquitetura.

Não implementar:

- HTTP
- Login
- JWT
- Refresh Token
- Middleware
- Tenant
- Frontend

Somente o caso de uso CreateUser.