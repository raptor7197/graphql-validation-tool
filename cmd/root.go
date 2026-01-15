package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Global flags
	cfgFile    string
	verbose    bool
	jsonOutput bool

	// Version info
	Version   = "1.0.0"
	BuildDate = "unknown"
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "gql-validate",
	Short: "A GraphQL query validation tool",
	Long: `GraphQL Validation Tool - Validate your GraphQL queries against a PostgreSQL database.

This tool uses GraphJin to compile and validate GraphQL queries against your
actual database schema, ensuring your queries are valid before deployment.

Examples:
  # Validate all queries in a directory
  gql-validate validate -q ./queries

  # Validate a single query file
  gql-validate validate -f ./queries/get_user.graphql

  # Check database connection
  gql-validate check

  # Initialize a new project with sample config
  gql-validate init`,
	Version: Version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Global flags available to all commands
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "config.yaml", "config file path")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose output")
	rootCmd.PersistentFlags().BoolVarP(&jsonOutput, "json", "j", false, "output results as JSON")

	// Set version template
	rootCmd.SetVersionTemplate(`{{printf "gql-validate version %s\n" .Version}}`)
}
