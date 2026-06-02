Você é um engenheiro de software sênior especializado em Golang, TDD, Clean Architecture e Domain-Driven Design.

Contexto do projeto:

Sistema SaaS para escritórios de advocacia.

Stack:

- Golang 1.26+
- PostgreSQL
- Chi Router
- Testify
- Docker
- Clean Architecture
- Modular Monolith
- API REST

Objetivo desta task:

Implementar a estrutura mínima do backend e um endpoint de health check.

Regras obrigatórias:

1. Trabalhar usando TDD.
2. Criar primeiro os testes.
3. Fazer o teste falhar.
4. Implementar apenas o necessário para passar.
5. Refatorar se necessário.
6. Não implementar funcionalidades futuras.
7. Não criar código morto.
8. Seguir princípios SOLID.
9. Código em inglês.
10. Comentários apenas quando estritamente necessários.

Critérios de aceite:

- Projeto deve compilar.
- Deve existir endpoint:

GET /health

Resposta:

HTTP 200

{
  "status": "ok"
}

Estrutura inicial esperada:

backend/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   └── health/
│       ├── handler.go
│       └── handler_test.go
├── go.mod
└── go.sum

Testes obrigatórios:

Scenario 1:
Given the API is running
When GET /health is called
Then should return HTTP 200

Scenario 2:
Given the API is running
When GET /health is called
Then response body must be:

{
  "status": "ok"
}

Entregáveis:

1. Estrutura de pastas.
2. Código completo.
3. Testes.
4. Comandos para execução.
5. Explicação rápida da arquitetura adotada.

Não implemente autenticação.
Não implemente banco de dados.
Não implemente usuários.

Somente a task solicitada.