# Quick Start Guide

## 🚀 Getting Started in 5 Minutes

This tool validates GraphQL queries against your PostgreSQL database using GraphJin's internal query compilation and execution pipeline. It checks for errors at every level of nesting in the response data.

### 1. Your Database is Ready!

Your PostgreSQL Docker container is already running:
- **Host:** localhost:5432
- **Database:** testdb
- **User:** testuser
- **Password:** testpass

### 2. Database Schema

Two tables are set up with sample data:

**users table:**
- id (serial primary key)
- name (varchar)
- email (varchar, unique)
- created_at (timestamp)

**posts table:**
- id (serial primary key)
- user_id (integer, foreign key)
- title (varchar)
- content (text)
- created_at (timestamp)

### 3. Run the Validator

```bash
./graphql-validator
```

That's it! The tool will:
1. Load all `.graphql` files from the `queries/` directory
2. Read corresponding `.json` files for variables (if they exist)
3. Use GraphJin's compile and execute pipeline to validate each query
4. Check for errors at all levels (compilation, execution, and nested response errors)
5. Report which queries pass or fail

### 4. Sample Queries Included

Four sample queries are ready to test:

**queries/get_users.graphql** - Get all users
```graphql
query GetUsers {
  users {
    id
    name
    email
    created_at
  }
}
```

**queries/get_user_by_id.graphql** - Get specific user (with variables)
```graphql
query GetUserById($id: Int!) {
  users(where: { id: { eq: $id } }) {
    id
    name
    email
    created_at
  }
}
```

**queries/get_posts_with_users.graphql** - Get posts with related user data
```graphql
query GetPostsWithUsers {
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

**queries/get_posts_by_user.graphql** - Filtered query with variables
```graphql
query GetPostsByUserId($userId: Int!) {
  posts(where: { user_id: { eq: $userId } }) {
    id
    title
    content
    user {
      id
      name
    }
  }
}
```

### 5. Command Options

```bash
# Verbose output
./graphql-validator -verbose

# JSON output (great for CI/CD)
./graphql-validator -format json

# Custom queries directory
./graphql-validator -queries /path/to/queries

# Custom config file
./graphql-validator -config /path/to/config.yaml
```

### 6. Adding Your Own Queries

1. Create a `.graphql` file in the `queries/` directory
2. Optionally create a `.json` file with the same name for variables
3. Run the validator

**Example:**

`queries/my_query.graphql`:
```graphql
query GetUserPosts($userId: Int!) {
  posts(where: { user_id: { eq: $userId } }) {
    title
    content
  }
}
```

`queries/my_query.json`:
```json
{
  "userId": 1
}
```

### 7. Understanding Output

**✓ PASS** - Query compiled and executed successfully with no errors at any level
**✗ FAIL** - Query failed due to:
  - Compilation errors (invalid syntax)
  - Execution errors (invalid columns, tables, or relationships)
  - Nested errors in the response data

Example output:
```
Query Validation Results
========================

✓ PASS  get_users.graphql (2ms)
✗ FAIL  invalid_query.graphql (0ms)
        • Execution error: column: 'users.invalid_field' not found
        • column: 'users.invalid_field' not found

Summary: 2 tests, 1 passed, 1 failed
```

Exit codes:
- `0` - All queries passed validation
- `1` - One or more queries failed (useful for CI/CD)

### 8. CI/CD Integration

The tool is designed for automated testing:

```bash
# Basic usage - fails if any query fails
./graphql-validator || exit 1

# JSON output for parsing results
./graphql-validator -format json

# Example JSON output
{
  "total": 4,
  "passed": 3,
  "failed": 1,
  "results": [
    {
      "name": "get_users.graphql",
      "path": "queries/get_users.graphql",
      "passed": true,
      "duration_ms": 2
    },
    {
      "name": "invalid_query.graphql",
      "path": "queries/invalid_query.graphql",
      "passed": false,
      "errors": [
        "Execution error: column not found",
        "column not found"
      ],
      "duration_ms": 0
    }
  ]
}
```

### 9. Error Detection

The validator detects multiple types of errors:

**Compilation Errors:**
- Invalid GraphQL syntax
- Malformed queries

**Execution Errors:**
- Non-existent tables: `table not found: public.invalid_table`
- Invalid columns: `column: 'users.invalid_field' not found`
- Invalid relationships or foreign keys
- Type mismatches

**Nested Errors:**
- Errors embedded in response data at any nesting level
- Error objects in arrays
- Both `error` and `errors` fields

### 10. Query Capabilities (GraphJin)

GraphJin automatically supports:
- **Filtering:** `where: { field: { eq: value } }`
- **Sorting:** `order_by: { field: asc }`
- **Pagination:** `limit: 10, offset: 0`
- **Relationships:** Automatic joins based on foreign keys
- **Aggregations:** `count`, `sum`, `avg`, `min`, `max`

**Example with filters and sorting:**
```graphql
query GetRecentPosts {
  posts(
    where: { user_id: { eq: 1 } }
    order_by: { created_at: desc }
    limit: 5
  ) {
    title
    created_at
  }
}
```

### 11. Database Connection via Environment Variables

Override config.yaml with environment variables:

```bash
export DB_HOST=localhost
export DB_NAME=testdb
export DB_USER=testuser
export DB_PASSWORD=testpass
export DB_SSLMODE=disable

./graphql-validator
```

## 🔧 Troubleshooting

**Query validation failing unexpectedly?**
Run with verbose mode to see detailed output:
```bash
./graphql-validator -verbose
```

**Database connection failed?**
```bash
# Test connection
PGPASSWORD=testpass psql -h localhost -U testuser -d testdb -c "\dt"
```

**No tables found?**
```bash
# Check if tables exist
PGPASSWORD=testpass psql -h localhost -U testuser -d testdb -c "SELECT table_name FROM information_schema.tables WHERE table_schema='public';"
```

**Query validation failed?**
- Check column names match your database schema exactly (case-sensitive)
- Verify table relationships and foreign keys are defined
- Ensure variables in `.json` file match query parameters
- Confirm table names exist in the database
- Check that the user has SELECT permissions on tables

**Understanding error messages:**
- `column: 'table.column' not found` - Column doesn't exist or is misspelled
- `table not found: public.table_name` - Table doesn't exist in schema
- `Error at root[0]` - Error in the first element of result array
- `Error at users.posts[1]` - Error in nested relationship

## 📚 How It Works

The validator uses GraphJin's internal pipeline:

1. **Load Queries**: Reads all `.graphql` files from the queries directory
2. **Load Variables**: Reads corresponding `.json` files (or uses `{}` if not found)
3. **Compile**: GraphJin compiles the GraphQL query into SQL
4. **Execute**: Runs the compiled SQL against your database
5. **Validate Response**: Checks for errors at three levels:
   - Top-level execution errors
   - GraphQL errors in Result.Errors
   - Nested errors in the response data
6. **Report**: Shows which queries passed or failed with detailed error messages

## 📚 Next Steps

1. Add your actual database schema to the testdb (or point to your real database)
2. Create `.graphql` files for your real queries
3. Add `.json` files for queries that need variables
4. Run validation as part of your development workflow
5. Integrate into CI/CD pipeline to catch breaking changes early

## 🎯 Common Use Cases

### Pre-deployment Validation
```bash
# Validate all queries before deploying
./graphql-validator || (echo "Query validation failed!" && exit 1)
```

### Development Workflow
```bash
# Watch mode (using entr or similar)
ls queries/*.graphql | entr -c ./graphql-validator -verbose
```

### Testing Schema Changes
```bash
# After altering database schema, validate all queries still work
./graphql-validator -verbose
```

Happy validating! 🚀