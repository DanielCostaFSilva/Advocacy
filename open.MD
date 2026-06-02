# OPEN.md - Instruções permanentes para assistentes de IA

## Projeto
Sistema de gestão para escritório de advocacia (nome temporário: LegalFlow)

## Stack
- **Backend**: Go 1.21+ (Golang)
  - Framework: net/http padrão ou Gin (definir)
  - ORM: GORM ou sqlc
  - Banco: PostgreSQL
  - Autenticação: JWT
  - API: REST (JSON)
- **Frontend**: Angular 17+
  - UI: Angular Material ou PrimeNG
  - Estado: Services com RxJS
  - Build: Angular CLI

## Estrutura de pastas (planejada)
/backend
/cmd/api - entrypoint main.go
/internal
/domain - entidades e regras de negócio
/repository - acesso a dados
/service - lógica de aplicação
/handler - handlers HTTP (controllers)
/pkg - utilidades compartilháveis
/migrations - scripts SQL
/frontend
/src
/app
/core - services, guards, interceptors
/shared - componentes reutilizáveis
/features - módulos funcionais (clientes, processos, etc.)
/layout - menu, header, footer


## Convenções de código
- Go: `gofmt`, `golangci-lint`, nomes em inglês, injeção de dependência
- Angular: ESLint + Prettier, padrão de nomenclatura kebab-case para arquivos, camelCase para variáveis
- Commits: padrão Conventional Commits (feat, fix, docs, chore)

## Regras de negócio do escritório
- **Cliente**: PF ou PJ, com endereço, contatos, documentos (CPF/CNPJ)
- **Processo**: número único, tribunal, assunto, valor da causa, status (andamento, arquivado, sentenciado)
- **Advogado**: usuário do sistema com perfil (admin, advogado, assistente)
- **Prazos**: alertas automáticos (X dias antes do vencimento)
- **Honorários**: contratos, parcelas, recibos
- **Documentos**: upload (PDF, DOCX) com versionamento

## Metodologia de trabalho
- Dividir tudo em **pequenas tasks** de 1 a 4 horas
- Cada task deve resultar em um PR pequeno e testável
- Preferir código simples e legível a superengenharia
- Documentar APIs com Swagger (gerado automaticamente)
- Testes: unitários (Go: `testing` + `testify`, Angular: Jest) e testes de integração nas principais rotas

## Como me usar (a IA)
- Sempre leia este `open.md` antes de sugerir código ou tasks
- Proponha soluções alinhadas com a stack e arquitetura acima
- Para cada task, descreva: objetivo, arquivos alterados, breve passo a passo
- Evite soluções genéricas; personalize para escritório de advocacia