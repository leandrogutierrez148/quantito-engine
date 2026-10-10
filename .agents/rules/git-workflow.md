# Git Workflow & Branch Management Standards

This rule defines mandatory Git procedures for all AI agents and developers working on the **quantito-engine** repository.

---

## 1. Mandatory Branching Policy

Before starting work on any user requirement, feature, bug fix, or refactoring task, you **MUST** ensure you are working on a dedicated Git branch.

> [!IMPORTANT]
> **Never commit or push directly to `main` or `master`.** All changes must originate from a dedicated task branch.

---

## 2. Branch Naming Convention

All branch names must strictly follow the format:

```text
<type>/<short-kebab-case-description>
```

### Allowed Types:
- `fix/`: Bug fixes and error resolutions (e.g., `fix/mutex-unlock-timeout`, `fix/mysql-connection-leak`).
- `feat/`: New features, strategies, or capabilities (e.g., `feat/execution-orchestration`, `feat/broker-sink`).
- `chore/`: Maintenance, dependency updates, and configuration tasks (e.g., `chore/cleanup-logs`, `chore/docker-compose`).
- `refactor/`: Code refactoring without behavioral changes (e.g., `refactor/signal-repository`).
- `test/`: Adding or updating test suites (e.g., `test/engine-failure-cases`).
- `docs/`: Documentation updates (e.g., `docs/architecture-guide`).

### Rules for the Descriptor:
1. Must use **lowercase** and **kebab-case** (hyphen-separated words).
2. Must be concise, self-explanatory, and directly related to the task.
3. Example: `feat/experiment-executions-model`

---

## 3. Standard Git Workflow

1. **Verify Current State**:
   ```bash
   git status
   git branch
   ```

2. **Create and Switch to the New Branch**:
   ```bash
   git checkout -b <type>/<short-description>
   ```

3. **Verify Quality Before Committing**:
   ```bash
   # Run unit & integration test suites:
   make test
   ```

4. **Commit with Plain Message**:
   - Use a clear, concise, plain description without `<type>(scope):` prefixes.
   - Example:
     ```bash
     git add <modified-files>
     git commit -m "Add execution tracking repository"
     ```
