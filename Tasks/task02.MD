Você é um engenheiro de software sênior especializado em Golang, PostgreSQL, TDD, Clean Architecture e sistemas SaaS multi-tenant.

Contexto do projeto:

Sistema SaaS para escritórios de advocacia.

Stack:

- Golang 1.26+
- PostgreSQL 17
- Chi Router
- SQLC
- Goose Migrations
- Testify
- Docker
- Docker Compose
- Clean Architecture
- Modular Monolith
- API REST

Objetivo desta task:

Adicionar infraestrutura de banco de dados ao projeto.

IMPORTANTE:

Ainda não existe autenticação.
Ainda não existe módulo de usuários.
Ainda não existe regra de negócio.

Esta task é exclusivamente para criar a infraestrutura de persistência.

Regras obrigatórias:

1. Trabalhar usando TDD.
2. Criar testes primeiro.
3. Implementar apenas o necessário para passar.
4. Não antecipar funcionalidades futuras.
5. Código em inglês.
6. Seguir SOLID.
7. Seguir Clean Architecture.

Critérios de aceite:

- PostgreSQL deve subir via Docker Compose.
- Aplicação deve conectar ao banco.
- Deve existir sistema de migrations usando Goose.
- Deve existir teste de conexão.
- Deve existir tabela de controle de migrations.

Estrutura esperada:

backend/
├── cmd/
│   └── api/
├── internal/
│   ├── database/
│   │   ├── connection.go
│   │   └── connection_test.go
│   └── health/
├── migrations/
│   └── 000001_initial.sql
├── docker-compose.yml
├── .env.example
└── Makefile

Variáveis esperadas:

DB_HOST
DB_PORT
DB_USER
DB_PASSWORD
DB_NAME

Testes obrigatórios:

Scenario 1:
Given database is running
When application starts
Then database connection should succeed

Scenario 2:
Given migration command is executed
When migrations run
Then schema version should be registered

Comandos esperados:

make up
make down
make migrate-up
make migrate-down
make test

Entregáveis:

1. Estrutura criada.
2. Docker Compose.
3. Configuração PostgreSQL.
4. Configuração Goose.
5. Testes.
6. Instruções de execução.

Não criar tabela users.
Não criar autenticação.
Não criar domínio de negócio.

Somente infraestrutura de banco.