package cmd

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check database connection and configuration",
	Long: `Verify that the database connection is working correctly.

This command tests the database connection using the configured
credentials and reports the connection status.

Examples:
  # Check connection using default config
  gql-validate check
aml
  # Check connection with custom config
  gql-validate check -c /path/to/config.yaml

  # Check with verbose output
  gql-validate check -v`,
	RunE: runCheck,
}

func init() {
	rootCmd.AddCommand(checkCmd)
}

func runCheck(cmd *cobra.Command, args []string) error {
	fmt.Println("Checking configuration and database connection...")
	fmt.Println()

	// Load configuration
	fmt.Printf("  ○ Loading config from: %s\n", cfgFile)
	config, err := LoadConfig(cfgFile)
	if err != nil {
		fmt.Printf("  ✗ Failed to load config: %v\n", err)
		return err
	}
	fmt.Printf("  ✓ Config loaded successfully\n")

	// validating the  configuration
	fmt.Printf("  ○ Validating configuration...\n")
	if err := config.Validate(); err != nil {
		fmt.Printf("  ✗ Invalid configuration: %v\n", err)
		return err
	}
	fmt.Printf("  ✓ Configuration is valid\n")

	// Print connection details (hide password)
	if verbose {
		fmt.Println()
		fmt.Println("  Connection Details:")
		fmt.Printf("    Host:     %s\n", config.Database.Host)
		fmt.Printf("    Port:     %d\n", config.Database.Port)
		fmt.Printf("    Database: %s\n", config.Database.DBName)
		fmt.Printf("    User:     %s\n", config.Database.User)
		fmt.Printf("    SSL Mode: %s\n", config.Database.SSLMode)
		fmt.Println()
	}

	// Test database connection
	fmt.Printf("  ○ Connecting to database...\n")
	start := time.Now()

	db, err := sql.Open("pgx", config.GetDSN())
	if err != nil {
		fmt.Printf("  ✗ Failed to open database connection: %v\n", err)
		return err
	}
	defer db.Close()

	// Ping the database
	if err := db.Ping(); err != nil {
		fmt.Printf("  ✗ Failed to connect to database: %v\n", err)
		return err
	}

	elapsed := time.Since(start)
	fmt.Printf("  ✓ Database connection successful (%dms)\n", elapsed.Milliseconds())

	// Get database version
	var version string
	err = db.QueryRow("SELECT version()").Scan(&version)
	if err == nil && verbose {
		fmt.Printf("  ✓ Database version: %s\n", truncateString(version, 60))
	}

	// Check tables count
	var tableCount int
	err = db.QueryRow(`
		SELECT COUNT(*)
		FROM information_schema.tables
		WHERE table_schema = 'public'
	`).Scan(&tableCount)
	if err == nil {
		fmt.Printf("  ✓ Found %d table(s) in public schema\n", tableCount)
	}

	fmt.Println()
	fmt.Println("All checks passed! Your configuration is ready to use.")
	fmt.Println()

	return nil
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
