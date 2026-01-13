# GraphQL Validation Tool

A command-line tool for validating GraphQL queries against a PostgreSQL database using GraphJin.

## Features

- Validates GraphQL queries against your database schema
- Supports query variables via JSON files
- Multiple output formats (text, JSON)
- Detailed error reporting
- Configurable via YAML and environment variables

## Prerequisites

- Go 1.16 or higher
- PostgreSQL database
- GraphQL queries to validate

## Installation

```bash
go build -o graphql-validator
```

## Configuration

### Database Setup

You can configure the database connection in two ways:

#### Option 1: Using config.yaml

Edit `config.yaml` with your database credentials:

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

#### Option 2: Using Environment Variables (Recommended)

Set environment variables to override the config file:

```bash
export DB_HOST=localhost
export DB_NAME=your_database
export DB_USER=your_user
export DB_PASSWORD=your_password
export DB_SSLMODE=disable
```

Or create a `.env` file (copy from `.env.example`):

```bash
cp .env.example .env
# Edit .env with your credentials
```

Then source it before running:

```bash
source .env
./graphql-validator
```

**Note:** Environment variables take precedence over config.yaml values.

## Usage

### Basic Usage

```bash
./graphql-validator
```

This will use the default configuration (`config.yaml`) and look for queries in the `./queries` directory.

### Command Line Options

```bash
./graphql-validator [options]

Options:
  -config string
        Path to config file (default "config.yaml")
  -queries string
        Path to queries directory (default "./queries")
  -verbose
        Enable verbose output
  -format string
        Output format: text or json (default "text")
```

### Examples

```bash
# Use custom config file
./graphql-validator -config /path/to/config.yaml

# Use custom queries directory
./graphql-validator -queries /path/to/queries

# Verbose output
./graphql-validator -verbose

# JSON output format
./graphql-validator -format json

# Combine options
./graphql-validator -config custom.yaml -queries ./my-queries -verbose -format json
```

## Query Files

Place your GraphQL queries in the queries directory with `.graphql` extension:

```
queries/
├── get_user.graphql
├── get_user.json (optional variables)
├── list_products.graphql
└── list_products.json (optional variables)
```

### Query Example

**queries/get_user.graphql**
```graphql
query GetUser($id: Int!) {
  users(where: { id: { eq: $id } }) {
    id
    name
    email
  }
}
```

**queries/get_user.json** (optional)
```json
{
  "id": 1
}
```

If no JSON file is provided, the query will be executed with empty variables `{}`.

## Output

### Text Format (Default)

```
Query Validation Results:
========================
✓ PASS: get_user.graphql (45ms)
✗ FAIL: invalid_query.graphql (12ms)
    Error: GraphQL execution error: column not found

Summary: 2 tests, 1 passed, 1 failed
```

### JSON Format

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
      "error": "GraphQL execution error: column not found",
      "duration_ms": 12
    }
  ]
}
```

## Exit Codes

- `0`: All queries passed validation
- `1`: One or more queries failed validation

This makes it easy to integrate into CI/CD pipelines:

```bash
./graphql-validator || exit 1
```

## Troubleshooting

### Database Connection Errors

If you see errors like:
```
failed to ping database: failed to connect to `host=localhost user=your_user database=your_database`
```

Make sure:
1. PostgreSQL is running
2. Database credentials are correct
3. Database exists
4. User has access to the database
5. Host/port are correct

### Testing Database Connection

```bash
# Test with psql
psql -h localhost -U your_user -d your_database

# Or using environment variables
psql -h $DB_HOST -U $DB_USER -d $DB_NAME
```

### Query Validation Errors

If queries fail validation:
- Check that your schema matches the query
- Verify table and column names
- Ensure relationships are properly defined
- Check that required variables are provided in JSON files

## Development

### Build

```bash
go build -o graphql-validator
```

### Run Tests

```bash
go test ./...
```

### Dependencies

This tool uses:
- [GraphJin](https://github.com/dosco/graphjin) - GraphQL to SQL compiler
- [pgx](https://github.com/jackc/pgx) - PostgreSQL driver
- [yaml.v2](https://gopkg.in/yaml.v2) - YAML parser

## License

[Add your license here]