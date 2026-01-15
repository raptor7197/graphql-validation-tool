package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	initDir   string
	overwrite bool
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new GraphQL validation project",
	Long: `Create a new project with sample configuration and query files.

This command creates the necessary directory structure and files
to get started with GraphQL query validation.

Examples:
  # Initialize in current directory
  gql-validate init

  # Initialize in a specific directory
  gql-validate init -d ./my-project

  # Overwrite existing files
  gql-validate init --overwrite`,
	RunE: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringVarP(&initDir, "dir", "d", ".", "directory to initialize the project in")
	initCmd.Flags().BoolVar(&overwrite, "overwrite", false, "overwrite existing files")
}

func runInit(cmd *cobra.Command, args []string) error {
	fmt.Printf("Initializing GraphQL validation project in: %s\n\n", initDir)

	// Create directories
	queriesDir := filepath.Join(initDir, "queries")
	if err := os.MkdirAll(queriesDir, 0755); err != nil {
		return fmt.Errorf("failed to create queries directory: %w", err)
	}
	fmt.Printf("  ✓ Created directory: %s\n", queriesDir)

	// Create config.yaml
	configPath := filepath.Join(initDir, "config.yaml")
	if err := writeFileIfNotExists(configPath, sampleConfig, overwrite); err != nil {
		return err
	}

	// Create .env.example
	envPath := filepath.Join(initDir, ".env.example")
	if err := writeFileIfNotExists(envPath, sampleEnv, overwrite); err != nil {
		return err
	}

	// Create sample query
	queryPath := filepath.Join(queriesDir, "get_users.graphql")
	if err := writeFileIfNotExists(queryPath, sampleQuery, overwrite); err != nil {
		return err
	}

	// Create sample query with variables
	queryWithVarsPath := filepath.Join(queriesDir, "get_user_by_id.graphql")
	if err := writeFileIfNotExists(queryWithVarsPath, sampleQueryWithVars, overwrite); err != nil {
		return err
	}

	// Create variables file
	varsPath := filepath.Join(queriesDir, "get_user_by_id.json")
	if err := writeFileIfNotExists(varsPath, sampleVars, overwrite); err != nil {
		return err
	}

	// Create .gitignore
	gitignorePath := filepath.Join(initDir, ".gitignore")
	if err := writeFileIfNotExists(gitignorePath, sampleGitignore, overwrite); err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("Project initialized successfully!")
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  1. Edit config.yaml with your database credentials")
	fmt.Println("     Or set environment variables: DB_HOST, DB_NAME, DB_USER, DB_PASSWORD")
	fmt.Println()
	fmt.Println("  2. Check your database connection:")
	fmt.Println("     gql-validate check")
	fmt.Println()
	fmt.Println("  3. Add your GraphQL queries to the queries/ directory")
	fmt.Println()
	fmt.Println("  4. Run validation:")
	fmt.Println("     gql-validate validate")
	fmt.Println()

	return nil
}

func writeFileIfNotExists(path, content string, overwrite bool) error {
	if _, err := os.Stat(path); err == nil && !overwrite {
		fmt.Printf("  ○ Skipped (exists): %s\n", path)
		return nil
	}

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", path, err)
	}

	fmt.Printf("  ✓ Created: %s\n", path)
	return nil
}

const sampleConfig = `# GraphQL Validation Tool Configuration
# Database credentials can be overridden with environment variables:
# DB_HOST, DB_PORT, DB_NAME, DB_USER, DB_PASSWORD, DB_SSLMODE

database:
  type: "postgres"
  host: "localhost"
  port: 5432
  dbname: "your_database"
  user: "your_user"
  password: "your_password"
  sslmode: "disable"

# Set to true for production mode (disables debug output)
production: false
`

const sampleEnv = `# Database Configuration
# Copy this file to .env and fill in your values
# Then run: source .env

export DB_HOST=localhost
export DB_PORT=5432
export DB_NAME=your_database
export DB_USER=your_user
export DB_PASSWORD=your_password
export DB_SSLMODE=disable
`

const sampleQuery = `# Sample query to fetch all users
# Modify this to match your database schema

query GetUsers {
  users {
    id
    name
    email
    created_at
  }
}
`

const sampleQueryWithVars = `# Sample query with variables
# Variables are provided in the corresponding .json file

query GetUserById($id: Int!) {
  users(where: { id: { eq: $id } }) {
    id
    name
    email
    created_at
  }
}
`

const sampleVars = `{
  "id": 1
}
`

const sampleGitignore = `# Environment files with secrets
.env

# Binary
gql-validate

# OS files
.DS_Store
Thumbs.db

# IDE
.idea/
.vscode/
*.swp
*.swo
`
