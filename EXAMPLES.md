# GraphQL Validator - Usage Examples

## Overview

This document provides comprehensive examples of using the GraphQL validation tool with GraphJin to validate queries against your PostgreSQL database.

## Basic Examples

### Example 1: Simple Query

**queries/get_all_users.graphql**
```graphql
query GetAllUsers {
  users {
    id
    name
    email
    created_at
  }
}
```

**Result:**
```
✓ PASS  get_all_users.graphql (2ms)
```

---

### Example 2: Query with Variables

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

**queries/get_user.json**
```json
{
  "id": 1
}
```

**Result:**
```
✓ PASS  get_user.graphql (1ms)
```

---

### Example 3: Query with Relationships

**queries/get_posts_with_authors.graphql**
```graphql
query GetPostsWithAuthors {
  posts {
    id
    title
    content
    created_at
    user {
      id
      name
      email
    }
  }
}
```

**Result:**
```
✓ PASS  get_posts_with_authors.graphql (3ms)
```

---

## Advanced Examples

### Example 4: Filtered Query with Variables

**queries/get_user_posts.graphql**
```graphql
query GetUserPosts($userId: Int!, $limit: Int) {
  posts(
    where: { user_id: { eq: $userId } }
    order_by: { created_at: desc }
    limit: $limit
  ) {
    id
    title
    content
    created_at
  }
}
```

**queries/get_user_posts.json**
```json
{
  "userId": 1,
  "limit": 10
}
```

---

### Example 5: Complex Filtering

**queries/search_posts.graphql**
```graphql
query SearchPosts($searchTerm: String!) {
  posts(where: { 
    title: { ilike: $searchTerm }
  }) {
    id
    title
    content
    user {
      name
    }
  }
}
```

**queries/search_posts.json**
```json
{
  "searchTerm": "%First%"
}
```

---

### Example 6: Multiple Filters

**queries/get_recent_posts_by_user.graphql**
```graphql
query GetRecentPostsByUser($userId: Int!, $since: String!) {
  posts(
    where: { 
      and: [
        { user_id: { eq: $userId } },
        { created_at: { gt: $since } }
      ]
    }
    order_by: { created_at: desc }
  ) {
    id
    title
    created_at
  }
}
```

**queries/get_recent_posts_by_user.json**
```json
{
  "userId": 1,
  "since": "2024-01-01"
}
```

---

## Error Examples

### Example 7: Invalid Column Name

**queries/invalid_column.graphql**
```graphql
query InvalidColumn {
  users {
    id
    name
    nonexistent_column
  }
}
```

**Result:**
```
✗ FAIL  invalid_column.graphql (0ms)
        • Execution error: column: 'users.nonexistent_column' not found
        • column: 'users.nonexistent_column' not found
```

---

### Example 8: Invalid Table Name

**queries/invalid_table.graphql**
```graphql
query InvalidTable {
  nonexistent_table {
    id
    name
  }
}
```

**Result:**
```
✗ FAIL  invalid_table.graphql (0ms)
        • Execution error: table not found: public.nonexistent_table
        • table not found: public.nonexistent_table
```

---

### Example 9: Missing Required Variable

**queries/requires_variable.graphql**
```graphql
query RequiresVariable($id: Int!) {
  users(where: { id: { eq: $id } }) {
    id
    name
  }
}
```

**queries/requires_variable.json**
```json
{}
```

**Result:**
```
✗ FAIL  requires_variable.graphql (0ms)
        • Execution error: variable 'id' is required but not provided
```

---

## Command Line Usage Examples

### Example 10: Basic Validation

```bash
./graphql-validator
```

Output:
```
Query Validation Results
========================

✓ PASS  get_posts_with_users.graphql (3ms)
✓ PASS  get_user_by_id.graphql (2ms)
✓ PASS  get_users.graphql (1ms)

Summary: 3 tests, 3 passed, 0 failed
```

Exit code: `0`

---

### Example 11: Verbose Mode

```bash
./graphql-validator -verbose
```

Output:
```
✓ PASS: get_posts_with_users.graphql (took 3ms)
✓ PASS: get_user_by_id.graphql (took 2ms)
✓ PASS: get_users.graphql (took 1ms)

Query Validation Results
========================

✓ PASS  get_posts_with_users.graphql (3ms)
✓ PASS  get_user_by_id.graphql (2ms)
✓ PASS  get_users.graphql (1ms)

Summary: 3 tests, 3 passed, 0 failed
```

---

### Example 12: JSON Output

```bash
./graphql-validator -format json
```

Output:
```json
{
  "total": 3,
  "passed": 3,
  "failed": 0,
  "results": [
    {
      "name": "get_posts_with_users.graphql",
      "path": "queries/get_posts_with_users.graphql",
      "passed": true,
      "duration_ms": 3
    },
    {
      "name": "get_user_by_id.graphql",
      "path": "queries/get_user_by_id.graphql",
      "passed": true,
      "duration_ms": 2
    },
    {
      "name": "get_users.graphql",
      "path": "queries/get_users.graphql",
      "passed": true,
      "duration_ms": 1
    }
  ]
}
```

---

### Example 13: Custom Queries Directory

```bash
./graphql-validator -queries ./my-queries
```

---

### Example 14: Custom Config File

```bash
./graphql-validator -config ./staging-config.yaml
```

---

## CI/CD Integration Examples

### Example 15: GitHub Actions

```yaml
name: Validate GraphQL Queries

on: [push, pull_request]

jobs:
  validate:
    runs-on: ubuntu-latest
    
    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_DB: testdb
          POSTGRES_USER: testuser
          POSTGRES_PASSWORD: testpass
        ports:
          - 5432:5432
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
    
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Build validator
        run: go build -o graphql-validator
      
      - name: Run migrations
        run: |
          # Run your database migrations here
          psql -h localhost -U testuser -d testdb < schema.sql
        env:
          PGPASSWORD: testpass
      
      - name: Validate queries
        run: ./graphql-validator -format json
        env:
          DB_HOST: localhost
          DB_NAME: testdb
          DB_USER: testuser
          DB_PASSWORD: testpass
```

---

### Example 16: GitLab CI

```yaml
validate-queries:
  stage: test
  image: golang:1.21
  
  services:
    - postgres:15
  
  variables:
    POSTGRES_DB: testdb
    POSTGRES_USER: testuser
    POSTGRES_PASSWORD: testpass
    DB_HOST: postgres
    DB_NAME: testdb
    DB_USER: testuser
    DB_PASSWORD: testpass
  
  script:
    - go build -o graphql-validator
    - psql -h postgres -U testuser -d testdb < schema.sql
    - ./graphql-validator -format json
```

---

### Example 17: Jenkins Pipeline

```groovy
pipeline {
    agent any
    
    environment {
        DB_HOST = 'localhost'
        DB_NAME = 'testdb'
        DB_USER = 'testuser'
        DB_PASSWORD = credentials('db-password')
    }
    
    stages {
        stage('Build') {
            steps {
                sh 'go build -o graphql-validator'
            }
        }
        
        stage('Validate Queries') {
            steps {
                sh './graphql-validator -format json'
            }
        }
    }
    
    post {
        always {
            archiveArtifacts artifacts: 'queries/**/*.graphql', allowEmptyArchive: true
        }
    }
}
```

---

### Example 18: Pre-commit Hook

**.git/hooks/pre-commit**
```bash
#!/bin/bash

echo "Validating GraphQL queries..."

# Build if necessary
if [ ! -f ./graphql-validator ]; then
    go build -o graphql-validator
fi

# Run validator
./graphql-validator

if [ $? -ne 0 ]; then
    echo "❌ GraphQL validation failed!"
    echo "Please fix the errors before committing."
    exit 1
fi

echo "✅ GraphQL validation passed!"
exit 0
```

Make it executable:
```bash
chmod +x .git/hooks/pre-commit
```

---

## Environment-Specific Examples

### Example 19: Development Environment

```bash
export DB_HOST=localhost
export DB_NAME=myapp_dev
export DB_USER=dev_user
export DB_PASSWORD=dev_password
export DB_SSLMODE=disable

./graphql-validator -verbose
```

---

### Example 20: Staging Environment

```bash
export DB_HOST=staging-db.example.com
export DB_NAME=myapp_staging
export DB_USER=staging_user
export DB_PASSWORD=$STAGING_DB_PASSWORD
export DB_SSLMODE=require

./graphql-validator -config staging-config.yaml
```

---

### Example 21: Using .env File

**.env**
```bash
DB_HOST=localhost
DB_NAME=testdb
DB_USER=testuser
DB_PASSWORD=testpass
DB_SSLMODE=disable
```

**Usage:**
```bash
source .env
./graphql-validator
```

Or with direnv:
```bash
# Install direnv first
brew install direnv  # macOS
# or: apt-get install direnv  # Linux

# Create .envrc
echo 'source .env' > .envrc
direnv allow

# Now just run
./graphql-validator
```

---

## Testing Patterns

### Example 22: Test Organization

```
queries/
├── users/
│   ├── get_user.graphql
│   ├── get_user.json
│   ├── list_users.graphql
│   └── search_users.graphql
├── posts/
│   ├── get_post.graphql
│   ├── get_post.json
│   ├── list_posts.graphql
│   └── user_posts.graphql
└── auth/
    ├── login.graphql
    ├── login.json
    └── verify_token.graphql
```

---

### Example 23: Shared Variables

Multiple queries can share the same variable file:

**queries/common_user.json**
```json
{
  "userId": 1
}
```

Use symbolic links:
```bash
ln -s common_user.json get_user_posts.json
ln -s common_user.json get_user_profile.json
```

---

### Example 24: Test Data Setup

**scripts/setup-test-data.sh**
```bash
#!/bin/bash

PGPASSWORD=testpass psql -h localhost -U testuser -d testdb <<EOF
TRUNCATE users, posts CASCADE;

INSERT INTO users (name, email) VALUES
  ('Alice', 'alice@example.com'),
  ('Bob', 'bob@example.com'),
  ('Charlie', 'charlie@example.com');

INSERT INTO posts (user_id, title, content) VALUES
  (1, 'First Post', 'Content 1'),
  (1, 'Second Post', 'Content 2'),
  (2, 'Bob Post', 'Content 3');
EOF

echo "Test data setup complete!"
```

Run before validation:
```bash
./scripts/setup-test-data.sh
./graphql-validator
```

---

## Filtering and Sorting Examples

### Example 25: Equality Filter

```graphql
query GetUserById($id: Int!) {
  users(where: { id: { eq: $id } }) {
    id
    name
  }
}
```

### Example 26: Comparison Filters

```graphql
query GetRecentUsers($since: String!) {
  users(where: { created_at: { gt: $since } }) {
    id
    name
    created_at
  }
}
```

Supported operators:
- `eq` - equals
- `neq` - not equals
- `gt` - greater than
- `gte` - greater than or equal
- `lt` - less than
- `lte` - less than or equal
- `like` - SQL LIKE
- `ilike` - case-insensitive LIKE
- `in` - in array
- `is_null` - is NULL

### Example 27: Sorting

```graphql
query GetUsersSorted {
  users(order_by: { created_at: desc, name: asc }) {
    id
    name
    created_at
  }
}
```

### Example 28: Pagination

```graphql
query GetUsersPaginated($limit: Int!, $offset: Int!) {
  users(limit: $limit, offset: $offset, order_by: { id: asc }) {
    id
    name
    email
  }
}
```

**variables:**
```json
{
  "limit": 10,
  "offset": 0
}
```

---

## Summary

This document covered:

- ✅ Basic query validation
- ✅ Queries with variables
- ✅ Relationship queries
- ✅ Error detection and reporting
- ✅ Command-line usage
- ✅ CI/CD integration
- ✅ Environment-specific configurations
- ✅ Testing patterns and organization
- ✅ Filtering, sorting, and pagination

For more information, see:
- [QUICKSTART.md](QUICKSTART.md) - Get started in 5 minutes
- [README.md](README.md) - Complete documentation
- [GraphJin Documentation](https://graphjin.com/docs) - GraphJin features

---

**Happy Validating! 🚀**