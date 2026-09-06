# Quantito Engine

**Quantito Engine** is the execution and orchestration engine of **Quantito**, an automated trading and financial strategy platform.

It serves as the central hub for evaluating strategies, generating discrete trading signals, and dispatching signals to active destinations (consoles, brokers, ledgers) under a parallel contraposition model.

---

## Core Architecture & Responsibilities

1. **Contraposition Execution**: Evaluates all strategy variants of an experiment concurrently across the same asset universe, guaranteeing unbiased comparative performance tracking.
2. **Experiments & Experiment Assets**: Defines and manages quantitative experiments and assigns the specific assets (`ExperimentAssets`) participating in each study.
3. **Strategy Evaluation & Signal Generation**: Evaluates quantitative strategies (e.g. `fifty_day_avg`, `price_momentum`) for each asset to emit standardized trading signals: **`BUY`**, **`HOLD`**, or **`SELL`**.
4. **Signal Sinks & Dispatching**: Dispatches generated signals through a pluggable `SignalDispatcher` to active sinks (ledger, brokers, Telegram, BigQuery).
5. **Execution Ledger**: Persists the lifecycle, prices, confidence levels, and historical outcomes of every strategy evaluation in the database.

---

## Key Domain Concepts

### Experiments

- **Experiment**: Orchestration unit in charge of:
  - Ensuring concurrent execution of the strategies to compare.
  - Defining participating assets of the experiment, via `experimentAssets`.
  - Evaluating signals emitted by the participating strategies.

- **ExperimentAsset**: Associative entity linking an asset (`AssetID`) to an active experiment (`ExperimentID`), determining eligibility.

### Variants & Strategies

- **Variant**: A variant within an experiment representing a distinct strategy configuration in competition.
- **Strategy**: The pluggable quantitative algorithm (e.g. `fifty_day_avg`, `price_momentum`) evaluating market data to produce a **`Signal`**.
- **Signal**: Standardized recommendation produced by a strategy:
  - `Type`: **`BUY`**, **`HOLD`**, or **`SELL`**.
  - `Price`: Asset price at the evaluation moment.
  - `Confidence`: Confidence score from $0.0$ to $1.0$.
  - `Timestamp`: UTC timestamp of the signal.

### Signal Sinks & Dispatcher

- **ISignalSink**: Pluggable destination for emitted signals (e.g. `ConsoleSink`, `BrokerSink`, `BigQuerySink`).
- **ISignalDispatcher**: Central router dispatching signals to all registered sinks.

### Signals (Execution History)

- **SignalRecord**: Historical ledger entry persisted in the `signals` table recording an asset's evaluation under an experiment variant along with its emitted signal, price, confidence, and timestamp.

---

## Dependencies

This project uses [mockery](https://vektra.github.io/mockery/latest/installation/) for interface mock generation and [husky](https://github.com/automation-co/husky#docs) for managing Git hooks:

```bash
mockery
husky install
```

---

## Configuration

The application is configured via `config.toml`:

```toml
# Web server configuration
[api]
listen = ":8080"

# Database connection
[database]
host = "dbhost:3306"
schema = "quantito-engine"
username = "quantito_user"
password = "quantito_password"

# Logging configuration
[log]
level = "debug" | "info" | "error"
logger = "console" | "split-json" | "json"
```

---

## Local Development & Makefile Commands

A standard [Makefile](./Makefile) is provided to simplify development, testing, and service orchestration:

### Quick Reference

| Command | Description |
| :--- | :--- |
| `make up` | Spin up all services (`mysql` and `quantito-engine`) with Docker Compose in the background. |
| `make up-local` | Run the `quantito-engine` service directly in the console using local Go runtime. |
| `make down` | Stop and remove Docker Compose containers. |
| `make prepare-integration-tests` | Spin up local environment for integration tests (starts the `mysql` container). |
| `make integration-tests` | Run the integration test suite (`go test -v ./tests/...`). |
| `make unit-tests` | Run all unit tests (`go test -v ./internal/... ./cmd/...`). |
| `make test` | Run both unit and integration test suites. |
| `make help` | Show available Makefile targets and descriptions. |

### Usage Examples

```bash
# Start all containers in background
make up

# Run locally in terminal (reads .dev/config.toml)
make up-local

# Prepare database container for running integration tests
make prepare-integration-tests

# Run test suites
make unit-tests
make integration-tests

# Teardown background containers
make down
```


