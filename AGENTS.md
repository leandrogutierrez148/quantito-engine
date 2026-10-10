# AGENTS.md - Quantito Engine Development & Implementation Standards

This file serves as the central reference for AI agents and developers working on the **quantito-engine** repository. It documents the system architecture, programming standards, testing procedures, and workflow policies.

---

## 1. Role in the Quantito Ecosystem

**`quantito-engine`** is the execution, orchestration, and benchmarking service of the **Quantito Financial Platform**:
- Built with **Go 1.22+**, **Echo v4**, and **MySQL 8.0**.
- Orchestrates quantitative experiments under a **parallel contraposition model**, evaluating multiple strategy variants concurrently across identical asset universes.
- Manages the lifecycle of **executions** (`executions`) and per-variant execution progress (`execution_variants`).
- Calls the stateless quantitative evaluation microservice **`quantito`** (`POST /api/v1/strategies/:name/evaluate`).
- Maintains an immutable signal audit log (`signals`) and dispatches emitted trade recommendations to configured destinations (`ConsoleSink`, brokers, ledgers).

### Directory Structure

```
quantito-engine/
├── cmd/
│   └── api-quantito-engine/     # Service entrypoint, router initialization, CLI config
├── db/
│   ├── 00-schema.sql            # MySQL table DDLs, foreign keys, stored procedures
│   └── 01-data.sql              # Seed baseline experiments, strategies, and variants
├── internal/
│   ├── config/                  # TOML configuration decoders (API, DB, Log)
│   ├── controllers/             # Echo HTTP request handlers
│   ├── infraestructures/       # MySQL connection, REST API client, mutex locks, logger
│   ├── interfaces/              # Decoupled domain interfaces & Mockery mocks
│   │   └── mocks/               # Generated mock implementations
│   ├── models/                  # Domain entities (Experiment, Variant, Execution, Signal, etc.)
│   ├── repositories/            # MySQL data access layer (Stored procedure wrappers)
│   ├── services/                # Business logic & execution orchestration
│   ├── sinks/                   # Pluggable signal destinations (Console, Dispatcher)
│   ├── strategies/              # Strategy adapters calling quantito API
│   └── utils/                   # Shared validation and parsing utilities
├── tests/
│   └── integration/             # Integration test suite with Docker Compose
│       ├── configs/             # Test TOML configurations
│       ├── db/                  # Test DB connection and cleanup helpers
│       ├── tests/               # End-to-end integration tests
│       └── docker-compose.yml   # Multi-container setup (mysql + quantito + engine)
├── Dockerfile
├── Makefile
├── go.mod
└── README.md
```

---

## 2. Go Standards & Architectural Guidelines

### 2.1. Clean Architecture & Interfaces
- Keep logic decoupled by relying on interfaces defined in `internal/interfaces/`.
- Every external dependency (DB, HTTP client, Clock, Mutex) must be injected via its interface.
- Run `mockery` whenever an interface is added or updated in `internal/interfaces/`.

### 2.2. Error Handling & Logging
- Wrap errors with descriptive context using `fmt.Errorf("...: %w", err)`.
- Use structured logging with **`zerolog`** (`github.com/rs/zerolog/log`).
- Include relevant contextual fields in log statements: `experimentId`, `executionId`, `variantId`, `asset`, `status`.

### 2.3. Temporal & Concurrency Protection
- Always use UTC timestamps (`time.Now().UTC()`) for executions, signals, and database records.
- Protect concurrent evaluations per asset/variant using `IExecutionMutex` (`executionMutex.Lock(experimentID, variantID, asset)`).

---

## 3. Quality Verification & Testing Commands

Run before every commit:

```bash
# 1. Run unit tests
go test -v ./internal/... ./cmd/...

# 2. Run integration tests (requires MySQL running)
go test -v ./tests/...

# 3. Or run all tests via Makefile
make test

# 4. Regenerate mocks if interfaces changed
mockery
```

---

## 4. Git Workflow & Branch Management

All AI agents and developers must adhere to the mandatory branch policy:
- **Dedicated Branching**: Create a new branch `<type>/<short-kebab-case-description>` (e.g., `feat/execution-orchestration`, `fix/mysql-connection-leak`, `test/engine-failure-cases`).
- **No Direct Commits**: Direct commits to `main` or `master` are strictly forbidden.
- **Commit Messages**: Plain, concise descriptions without `<type>(scope):` prefixes.
- **Full Specification**: Detailed rule definition is located in [`.agents/rules/git-workflow.md`](file:///.agents/rules/git-workflow.md).
