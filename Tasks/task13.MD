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
- Clean Architecture
- Modular Monolith

Estado atual do projeto:

Já existem:

✓ Client Entity
✓ Client Repository
✓ PostgreSQL Persistence
✓ SQLC
✓ Migration clients
✓ Autenticação completa
✓ JWT Middleware

Objetivo desta task:

Implementar o caso de uso responsável pelo cadastro de clientes.

IMPORTANTE:

Ainda NÃO implementar:

- Endpoints HTTP
- Atualização de clientes
- Exclusão de clientes
- Paginação
- Processos jurídicos
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
10. Toda persistência deve ocorrer através de interfaces.

Objetivos da implementação:

1. CreateClient Use Case
2. DTOs Input/Output
3. Verificação de CPF duplicado
4. Integração com Client Entity
5. Testes unitários

Estrutura esperada:

backend/
├── internal/
│
│   ├── application/
│   │   └── client/
│   │       ├── create_client.go
│   │       ├── create_client_test.go
│   │       ├── input.go
│   │       ├── output.go
│   │       └── errors.go
│   │
│   ├── domain/
│   │   └── client/
│   │
│   └── infrastructure/
│
└── tests/

Input esperado:

type CreateClientInput struct {
    Name  string
    CPF   string
    Email string
    Phone string
}

Output esperado:

type CreateClientOutput struct {
    ID    string
    Name  string
    CPF   string
    Email string
    Phone string
}

Interface esperada:

type ClientRepository interface {

    Create(
        ctx context.Context,
        client *Client,
    ) error

    FindByCPF(
        ctx context.Context,
        cpf string,
    ) (*Client, error)
}

Fluxo esperado:

1. Receber dados
2. Verificar CPF existente
3. Criar entidade Client
4. Persistir cliente
5. Retornar resultado

Erros esperados:

ErrClientAlreadyExists

ErrInvalidInput

Validações da camada de aplicação:

Name

- obrigatório

CPF

- obrigatório

As demais validações continuam sendo responsabilidade da entidade Client.

Comportamentos obrigatórios:

CPF inexistente

→ cria cliente

CPF existente

→ retorna ErrClientAlreadyExists

Erro de persistência

→ retorna erro

Testes obrigatórios:

Scenario:
Given valid input
When creating client
Then should create successfully

Scenario:
Given duplicated CPF
When creating client
Then should return ErrClientAlreadyExists

Scenario:
Given invalid CPF
When creating client
Then should return domain error

Scenario:
Given repository create failure
When creating client
Then should return error

Scenario:
Given repository search failure
When creating client
Then should return error

Scenario:
Given valid input
When creating client
Then should invoke repository create

Mocks esperados:

ClientRepositoryMock

Não utilizar:

- Banco de dados
- PostgreSQL
- HTTP
- JWT

Todos os testes devem ser unitários.

Fluxo arquitetural esperado:

CreateClientUseCase
        ↓
ClientRepository
        ↓
Client Entity
        ↓
Persistência

Critérios de aceite:

- CreateClient Use Case implementado.
- CPF duplicado tratado.
- Dependências invertidas corretamente.
- Sem acesso direto ao banco.
- Testes unitários passando.

Entregáveis:

1. CreateClient Use Case.
2. DTOs Input/Output.
3. Errors padronizados.
4. Mocks.
5. Testes unitários.
6. Explicação rápida da arquitetura.

Não implementar:

- HTTP
- Paginação
- Atualização
- Exclusão
- Processos
- Audiências
- Uploads

Somente o caso de uso CreateClient.