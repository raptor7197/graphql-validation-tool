# GraphQL Validation Tool (gql-validate)

A powerful command-line tool for validating GraphQL queries against a PostgreSQL database using GraphJin.

## Features

- 🔍 **Validate GraphQL queries** against your actual database schema
- 📁 **Batch validation** of all queries in a directory
- 📄 **Single file validation** for quick checks
- 🔌 **Database connection testing** before running validations
- 🚀 **Project scaffolding** with the `init` command
- 📋 **List queries** with metadata and variable file detection
- 🎨 **Multiple output formats** (text, JSON)
- ⚡ **Fail-fast mode** for CI/CD pipelines
- 🔧 **Environment variable support** for secure credential management
- 🐚 **Shell completion** for bash, zsh, fish, and PowerShell

## Installation

### From Source

```bash
git clone https://github.com/your-repo/graphql-validation-tool.git
cd graphql-validation-tool
go build -o gql-validate .
```

### Move to PATH (optional)

```bash
sudo mv gql-validate /usr/local/bin/
```

## Quick Start

```bash
# Initialize a new project with sample config and queries
gql-validate init

# Edit config.yaml with your database credentials
# Or set environment variables (recommended)
export DB_HOST=localhost
export DB_NAME=mydb
export DB_USER=myuser
export DB_PASSWORD=mypassword

# Check database connection
gql-validate check

# Validate all queries
gql-validate validate
```

## Commands

### `validate` - Validate GraphQL Queries

Validate GraphQL queries against your PostgreSQL database schema.

```bash
# Validate all queries in the default ./queries directory
gql-validate validate

# Validate all queries in a specific directory
gql-validate validate -q ./my-queries

# Validate a single query file
gql-validate validate -f ./queries/get_user.graphql

# Validate with verbose output
gql-validate validate -v

# Stop on first failure (useful for CI/CD)
gql-validate validate --fail-fast

# Output results as JSON
gql-validate validate -j
```

### `check` - Check Database Connection

Verify that the database connection is working correctly.

```bash
# Check connection using default config
gql-validate check

# Check connection with custom config
gql-validate check -c /path/to/config.yaml

# Check with verbose output (shows DB version and table count)
gql-validate check -v
```

### `list` - List Available Queries

List all GraphQL query files in a directory with metadata.

```bash
# List all queries in default directory
gql-validate list

# List queries in a specific directory
gql-validate list -q ./my-queries

# Show full file paths
gql-validate list --full-path

# Output as JSON
gql-validate list -j
```

### `init` - Initialize a New Project

Create a new project with sample configuration and query files.

```bash
# Initialize in current directory
gql-validate init

# Initialize in a specific directory
gql-validate init -d ./my-project

# Overwrite existing files
gql-validate init --overwrite
```

### `completion` - Generate Shell Completion

Generate autocompletion scripts for your shell.

```bash
# Bash
gql-validate completion bash > /etc/bash_completion.d/gql-validate

# Zsh
gql-validate completion zsh > "${fpath[1]}/_gql-validate"

# Fish
gql-validate completion fish > ~/.config/fish/completions/gql-validate.fish

# PowerShell
gql-validate completion powershell > gql-validate.ps1
```

## Configuration

### config.yaml

```yaml
database:
  type: "postgres"
  host: "localhost"
  port: 5432
  dbname: "your_database"
  user: "your_user"
  password: "your_password"
  sslmode: "disable"

production: false
```

### Environment Variables

Environment variables take precedence over config.yaml values:

| Variable      | Description                    |
|---------------|--------------------------------|
| `DB_HOST`     | Database host                  |
| `DB_PORT`     | Database port                  |
| `DB_NAME`     | Database name                  |
| `DB_USER`     | Database user                  |
| `DB_PASSWORD` | Database password              |
| `DB_SSLMODE`  | SSL mode (disable/require/etc) |

**Recommended:** Use environment variables for credentials to avoid storing passwords in files.

```bash
# Create a .env file (don't commit to git!)
export DB_HOST=localhost
export DB_PORT=5432
export DB_NAME=mydb
export DB_USER=myuser
export DB_PASSWORD=secret
export DB_SSLMODE=disable

# Source it before running
source .env
gql-validate validate
```

## Query Files

Place your GraphQL queries in the queries directory with `.graphql` extension:

```
queries/
├── get_user.graphql
├── get_user.json          # Optional: variables for get_user.graphql
├── list_products.graphql
└── list_products.json     # Optional: variables for list_products.graphql
```

### Example Query

**queries/get_user.graphql**
```graphql
query GetUser($id: Int!) {
  users(where: { id: { eq: $id } }) {
    id
    name
    email
    created_at
  }
}
```

**queries/get_user.json** (optional variables)
```json
{
  "id": 1
}
```

If no JSON file is provided, the query will be executed with empty variables `{}`.

## Output Formats

### Text Output (Default)

```
╔══════════════════════════════════════════════════════════════╗
║            GraphQL Query Validation Results                  ║
╚══════════════════════════════════════════════════════════════╝

  ✓ PASS  get_user.graphql                            45ms
  ✗ FAIL  invalid_query.graphql                       12ms
          └─ Execution error: column "nonexistent" does not exist

──────────────────────────────────────────────────────────────────
  Summary: 2 total, 1 passed, 1 failed
```

### JSON Output (`-j` or `--json`)

```json
{
  "total": 2,
  "passed": 1,
  "failed": 1,
  "results": [
    {
      "name": "get_user.graphql",
      "path": "queries/get_user.graphql",
      "passed": true,
      "duration_ms": 45
    },
    {
      "name": "invalid_query.graphql",
      "path": "queries/invalid_query.graphql",
      "passed": false,
      "errors": ["Execution error: column \"nonexistent\" does not exist"],
      "duration_ms": 12
    }
  ]
}
```

## Exit Codes

| Code | Description                         |
|------|-------------------------------------|
| `0`  | All queries passed validation       |
| `1`  | One or more queries failed          |

This makes it easy to integrate into CI/CD pipelines:

```bash
gql-validate validate || exit 1
```

## CI/CD Integration

### GitHub Actions

```yaml
- name: Validate GraphQL Queries
  env:
    DB_HOST: ${{ secrets.DB_HOST }}
    DB_NAME: ${{ secrets.DB_NAME }}
    DB_USER: ${{ secrets.DB_USER }}
    DB_PASSWORD: ${{ secrets.DB_PASSWORD }}
  run: |
    ./gql-validate validate --fail-fast
```

### GitLab CI

```yaml
validate-queries:
  script:
    - ./gql-validate validate --fail-fast
  variables:
    DB_HOST: $DB_HOST
    DB_NAME: $DB_NAME
    DB_USER: $DB_USER
    DB_PASSWORD: $DB_PASSWORD
```

## Global Flags

These flags are available for all commands:

| Flag              | Short | Description                      | Default        |
|-------------------|-------|----------------------------------|----------------|
| `--config`        | `-c`  | Config file path                 | `config.yaml`  |
| `--verbose`       | `-v`  | Enable verbose output            | `false`        |
| `--json`          | `-j`  | Output results as JSON           | `false`        |
| `--help`          | `-h`  | Help for the command             |                |
| `--version`       |       | Version information              |                |

## Troubleshooting

### Database Connection Errors

If you see connection errors:

1. Verify PostgreSQL is running
2. Check credentials in config.yaml or environment variables
3. Ensure the database exists
4. Verify network connectivity (host/port)
5. Check SSL mode settings

```bash
# Test connection manually
psql -h $DB_HOST -U $DB_USER -d $DB_NAME

# Use the check command for detailed diagnostics
gql-validate check -v
```

### Query Validation Errors

If queries fail validation:

- Check that table and column names match your schema
- Verify relationships are properly defined in the database
- Ensure required variables are provided in JSON files
- Use `-v` for more detailed error information

## Dependencies

- [GraphJin](https://github.com/dosco/graphjin) - GraphQL to SQL compiler
- [pgx](https://github.com/jackc/pgx) - PostgreSQL driver
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [yaml.v2](https://gopkg.in/yaml.v2) - YAML parser

## License

MIT License - see LICENSE file for details