# Prompt para IA v0 - Gateway de Pagamentos Next.js

## Contexto do Projeto

Você é a IA v0, altamente qualificada para criar aplicações Next.js modernas e funcionais. Este projeto é um **Gateway de Pagamentos Educacional** que simula o funcionamento de um sistema real de processamento de transações financeiras.

## Objetivo

Criar uma aplicação Next.js completa com 4 páginas principais para demonstrar o funcionamento de um gateway de pagamentos, incluindo autenticação via API Key, listagem de faturas, detalhes de transações e criação de pagamentos.

## Requisitos Técnicos

### Framework e Tecnologias
- **Next.js 15** com App Router
- **TypeScript** para tipagem estática
- **Tailwind CSS v4** para estilização
- **shadcn/ui** para componentes de interface
- **React Server Components** quando apropriado
- **Client Components** para interatividade

### Estrutura de Arquivos
- Use **kebab-case** para nomes de arquivos
- Organize componentes em `components/` separados por funcionalidade
- Mantenha páginas em `app/` seguindo a estrutura do App Router
- Use hooks customizados em `hooks/` quando necessário

## Especificações de Design

### Tema e Cores
- **Tema Dark** obrigatório em todas as páginas
- **Paleta de cores limitada a 3-5 cores total:**
  - 1 cor primária (escolha uma cor que transmita confiança financeira - azul, verde ou roxo)
  - 2-3 cores neutras (tons de cinza, branco, preto)
  - 1 cor de destaque para ações importantes

### Tipografia
- **Máximo 2 fontes familiares**
- Use combinações Google Fonts recomendadas:
  - **Moderno/Tecnológico**: Space Grotesk Bold + DM Sans Regular
  - **Corporativo/Profissional**: Work Sans Bold + Open Sans Regular
- Implemente via `layout.tsx` e `globals.css`

### Layout
- **Mobile-first** obrigatório
- Breakpoints: mobile (320px), tablet (768px), desktop (1024px+)
- Use Flexbox como método principal de layout
- Espaçamento consistente: mínimo 16px entre seções, 8px entre elementos relacionados

## Componentes e Páginas

### 1. Navbar Superior (Componente Reutilizável)
- **Lado esquerdo**: Logo "Full Cycle Gateway" 
- **Lado direito**: "Olá, usuário" + botão de logout
- **Estilo**: Fixo no topo, tema dark, responsivo
- **Implementação**: Componente separado em `components/navbar.tsx`

### 2. Página de Autenticação (`app/auth/page.tsx`)
- **Campo principal**: Input para API Key com validação
- **Botão**: "Entrar" para autenticar
- **Validação**: Verificar se API Key foi preenchida
- **Redirecionamento**: Para dashboard após autenticação bem-sucedida
- **Estado**: Gerenciar API Key no localStorage ou contexto

### 3. Página de Listagem de Faturas (`app/invoices/page.tsx`)
- **Tabela responsiva** com colunas:
  - ID da Fatura
  - Valor
  - Status (pending/approved/rejected)
  - Data de Criação
  - Ações (ver detalhes)
- **Filtros**: Por status, data, valor
- **Paginação**: Para grandes volumes de dados
- **Status visual**: Cores diferentes para cada status
- **Botão**: "Nova Fatura" que leva à página de criação

### 4. Página de Detalhes da Fatura (`app/invoices/[id]/page.tsx`)
- **Informações completas** da fatura:
  - ID, Valor, Status, Descrição
  - Tipo de pagamento, últimos 4 dígitos do cartão
  - Datas de criação e atualização
  - Histórico de status (se aplicável)
- **Layout**: Card principal com seções organizadas
- **Botões**: Voltar para listagem, editar (se permitido)
- **Status visual**: Badge destacado com cor apropriada

### 5. Página de Criação de Fatura (`app/invoices/create/page.tsx`)
- **Formulário completo** com campos:
  - Valor (number, obrigatório)
  - Descrição (text, obrigatório)
  - Tipo de pagamento (select: credit_card, debit_card, pix)
  - Dados do cartão (se aplicável):
    - Número do cartão (mascarado)
    - CVV
    - Mês/Ano de expiração
    - Nome do titular
- **Validação**: Todos os campos obrigatórios, formato de cartão
- **Botões**: "Criar Fatura" e "Cancelar"
- **Feedback**: Loading state e mensagens de sucesso/erro

## Funcionalidades de Integração

### Autenticação
- **API Key**: Armazenar no localStorage ou contexto
- **Headers**: Incluir `X-API-Key` em todas as requisições
- **Proteção de rotas**: Middleware ou verificação em componentes

### API Integration
- **Base URL**: Configurável via variáveis de ambiente
- **Endpoints**:
  - `POST /accounts` - Criar conta
  - `GET /accounts` - Consultar conta
  - `POST /invoice` - Criar fatura
  - `GET /invoice` - Listar faturas
  - `GET /invoice/{id}` - Detalhes da fatura

### Estados e Loading
- **Loading states** para todas as operações assíncronas
- **Error handling** com mensagens amigáveis
- **Success feedback** para operações bem-sucedidas

## Diretrizes de Implementação

### Componentes
- **Separe responsabilidades**: Cada página deve ter componentes específicos
- **Reutilização**: Navbar, botões, inputs devem ser reutilizáveis
- **Props tipadas**: Use TypeScript interfaces para todas as props

### Estilização
- **Tailwind utilities**: Use classes utilitárias do Tailwind
- **Responsividade**: Implemente breakpoints consistentes
- **Acessibilidade**: ARIA labels, contraste adequado, navegação por teclado

### Performance
- **Server Components**: Use quando possível para melhor performance
- **Client Components**: Apenas onde necessário para interatividade
- **Lazy loading**: Para componentes pesados se aplicável

## Estrutura de Arquivos Esperada

```
app/
├── layout.tsx
├── globals.css
├── auth/
│   └── page.tsx
├── invoices/
│   ├── page.tsx
│   ├── create/
│   │   └── page.tsx
│   └── [id]/
│       └── page.tsx
components/
├── navbar.tsx
├── auth-form.tsx
├── invoice-list.tsx
├── invoice-details.tsx
├── invoice-form.tsx
└── ui/ (componentes shadcn)
hooks/
├── use-auth.ts
└── use-invoices.ts
lib/
├── api.ts
└── utils.ts
types/
└── index.ts
```

## Requisitos de Qualidade

### Código
- **TypeScript strict**: Sem `any`, tipagem completa
- **ESLint**: Seguir configurações padrão do Next.js
- **Formatação**: Prettier com configurações padrão
- **Imports**: Organizados e sem imports não utilizados

### UX/UI
- **Feedback visual**: Loading states, mensagens de erro/sucesso
- **Validação**: Em tempo real para formulários
- **Responsividade**: Funcionar perfeitamente em todos os dispositivos
- **Acessibilidade**: WCAG AA compliance

### Testes
- **Componentes**: Renderizam corretamente
- **Funcionalidades**: Formulários funcionam, navegação funciona
- **Responsividade**: Layout se adapta a diferentes tamanhos de tela

## Instruções Específicas para v0

1. **Comece sempre** com `SearchRepo` para entender a estrutura existente
2. **Use o TodoManager** se o projeto tiver múltiplas funcionalidades complexas
3. **Implemente mobile-first** com breakpoints responsivos
4. **Siga as diretrizes de cores** (máximo 5 cores, tema dark)
5. **Use shadcn/ui** para componentes de interface
6. **Implemente TypeScript** com interfaces completas
7. **Crie componentes reutilizáveis** e bem organizados
8. **Teste a responsividade** em diferentes tamanhos de tela

## Resultado Esperado

Uma aplicação Next.js completa, funcional e responsiva que demonstra um gateway de pagamentos com:
- Interface moderna e profissional em tema dark
- Navegação intuitiva entre as 4 páginas principais
- Formulários funcionais com validação
- Listagem e detalhes de faturas
- Autenticação via API Key
- Design responsivo mobile-first
- Código TypeScript bem estruturado e organizado
