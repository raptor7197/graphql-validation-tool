package cmd

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	graphjin "github.com/dosco/graphjin/core"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/spf13/cobra"
)

var (
	queriesDir string
	queryFile  string
	failFast   bool
)

// TestResult represents the result of validating a single query
type TestResult struct {
	Name     string   `json:"name"`
	Path     string   `json:"path"`
	Passed   bool     `json:"passed"`
	Errors   []string `json:"errors,omitempty"`
	Duration int64    `json:"duration_ms"`
}

// ValidationSummary represents the overall validation results
type ValidationSummary struct {
	Total   int          `json:"total"`
	Passed  int          `json:"passed"`
	Failed  int          `json:"failed"`
	Results []TestResult `json:"results"`
}

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate GraphQL queries against the database",
	Long: `Validate GraphQL queries against your PostgreSQL database schema.

This command reads .graphql files and validates them using GraphJin,
which compiles the queries against your actual database schema.

Examples:
  # Validate all queries in the default directory
  gql-validate validate

  # Validate all queries in a specific directory
  gql-validate validate -q ./my-queries

  # Validate a single query file
  gql-validate validate -f ./queries/get_user.graphql

  # Validate with verbose output
  gql-validate validate -v

  # Stop on first failure
  gql-validate validate --fail-fast`,
	RunE: runValidate,
}

func init() {
	rootCmd.AddCommand(validateCmd)

	validateCmd.Flags().StringVarP(&queriesDir, "queries", "q", "./queries", "directory containing GraphQL query files")
	validateCmd.Flags().StringVarP(&queryFile, "file", "f", "", "single GraphQL file to validate")
	validateCmd.Flags().BoolVar(&failFast, "fail-fast", false, "stop on first validation failure")
}

func runValidate(cmd *cobra.Command, args []string) error {
	// Load configuration
	config, err := LoadConfig(cfgFile)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if err := config.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// Initialize GraphJin
	gj, db, err := initializeGraphJin(config)
	if err != nil {
		return fmt.Errorf("failed to initialize GraphJin: %w", err)
	}
	defer db.Close()

	// Find query files to validate
	var queryFiles []string

	if queryFile != "" {
		// Validate single file
		if _, err := os.Stat(queryFile); os.IsNotExist(err) {
			return fmt.Errorf("query file not found: %s", queryFile)
		}
		queryFiles = []string{queryFile}
	} else {
		// Find all query files in directory
		queryFiles, err = findQueryFiles(queriesDir)
		if err != nil {
			return fmt.Errorf("failed to find query files: %w", err)
		}
	}

	if len(queryFiles) == 0 {
		fmt.Println("No query files found")
		return nil
	}

	if verbose {
		fmt.Printf("Found %d query file(s) to validate\n\n", len(queryFiles))
	}

	// Run validation
	results := validateQueries(gj, queryFiles)

	// Print results
	printResults(results)

	// Return error if any tests failed
	if results.Failed > 0 {
		return fmt.Errorf("%d validation(s) failed", results.Failed)
	}

	return nil
}

func initializeGraphJin(config *Config) (*graphjin.GraphJin, *sql.DB, error) {
	// Connect to database
	db, err := sql.Open("pgx", config.GetDSN())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Create GraphJin configuration
	gjConfig := &graphjin.Config{
		Debug:            verbose,
		Production:       config.Production,
		DisableAllowList: true,
		DefaultBlock:     false,
	}

	// Initialize GraphJin
	gj, err := graphjin.NewGraphJin(gjConfig, db)
	if err != nil {
		db.Close()
		return nil, nil, fmt.Errorf("failed to create GraphJin instance: %w", err)
	}

	return gj, db, nil
}

func findQueryFiles(dir string) ([]string, error) {
	var queryFiles []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".graphql") {
			queryFiles = append(queryFiles, path)
		}

		return nil
	})

	return queryFiles, err
}

func validateQueries(gj *graphjin.GraphJin, queryFiles []string) ValidationSummary {
	summary := ValidationSummary{
		Total:   len(queryFiles),
		Results: make([]TestResult, 0, len(queryFiles)),
	}

	for _, qf := range queryFiles {
		result := validateSingleQuery(gj, qf)
		summary.Results = append(summary.Results, result)

		if result.Passed {
			summary.Passed++
		} else {
			summary.Failed++
			if failFast {
				break
			}
		}
	}

	return summary
}

func validateSingleQuery(gj *graphjin.GraphJin, queryPath string) TestResult {
	result := TestResult{
		Name:   filepath.Base(queryPath),
		Path:   queryPath,
		Passed: false,
		Errors: []string{},
	}

	start := time.Now()

	// Read query file
	query, err := os.ReadFile(queryPath)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to read query file: %v", err))
		result.Duration = time.Since(start).Milliseconds()
		return result
	}

	// Look for corresponding JSON file with variables
	jsonFile := strings.TrimSuffix(queryPath, ".graphql") + ".json"
	var variables json.RawMessage

	if _, err := os.Stat(jsonFile); err == nil {
		jsonData, err := os.ReadFile(jsonFile)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Failed to read variables file: %v", err))
			result.Duration = time.Since(start).Milliseconds()
			return result
		}
		variables = json.RawMessage(jsonData)

		if verbose {
			fmt.Printf("  Using variables from: %s\n", filepath.Base(jsonFile))
		}
	} else {
		variables = json.RawMessage("{}")
	}

	// Execute query
	ctx := context.Background()
	res, err := gj.GraphQL(ctx, string(query), variables, nil)

	result.Duration = time.Since(start).Milliseconds()

	// Check for execution errors
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Execution error: %v", err))
	}

	// Check for GraphQL errors in the response
	if res != nil && len(res.Errors) > 0 {
		for _, gjErr := range res.Errors {
			result.Errors = append(result.Errors, gjErr.Message)
		}
	}

	// Check for nested errors in the response data
	if res != nil && len(res.Data) > 0 {
		nestedErrors := findNestedErrors(res.Data)
		if len(nestedErrors) > 0 {
			result.Errors = append(result.Errors, nestedErrors...)
		}
	}

	// Query passes only if there are no errors at any level
	if len(result.Errors) == 0 {
		result.Passed = true
	}

	return result
}

// findNestedErrors recursively searches for error fields in the GraphQL response data
func findNestedErrors(data json.RawMessage) []string {
	var errors []string

	if len(data) == 0 {
		return errors
	}

	var result interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return errors
	}

	collectErrors(result, &errors, "")
	return errors
}

// collectErrors recursively walks through the data structure looking for error indicators
func collectErrors(data interface{}, errors *[]string, path string) {
	switch v := data.(type) {
	case map[string]interface{}:
		// Check for "errors" key (array of errors)
		if errs, ok := v["errors"]; ok && errs != nil {
			if errArray, ok := errs.([]interface{}); ok && len(errArray) > 0 {
				for i, e := range errArray {
					if errMap, ok := e.(map[string]interface{}); ok {
						if msg, ok := errMap["message"].(string); ok {
							location := path
							if location == "" {
								location = "root"
							}
							*errors = append(*errors, fmt.Sprintf("Error at %s[%d]: %s", location, i, msg))
						} else {
							*errors = append(*errors, fmt.Sprintf("Error at %s[%d]: %v", path, i, e))
						}
					}
				}
			}
		}

		// Check for "error" key (single error)
		if errVal, ok := v["error"]; ok && errVal != nil {
			switch errStr := errVal.(type) {
			case string:
				if errStr != "" {
					location := path
					if location == "" {
						location = "root"
					}
					*errors = append(*errors, fmt.Sprintf("Error at %s: %s", location, errStr))
				}
			case map[string]interface{}:
				if msg, ok := errStr["message"].(string); ok {
					location := path
					if location == "" {
						location = "root"
					}
					*errors = append(*errors, fmt.Sprintf("Error at %s: %s", location, msg))
				}
			}
		}

		// Recursively check nested objects
		for key, value := range v {
			newPath := key
			if path != "" {
				newPath = path + "." + key
			}
			collectErrors(value, errors, newPath)
		}

	case []interface{}:
		// Recursively check arrays
		for i, item := range v {
			newPath := fmt.Sprintf("%s[%d]", path, i)
			if path == "" {
				newPath = fmt.Sprintf("[%d]", i)
			}
			collectErrors(item, errors, newPath)
		}
	}
}

func printResults(summary ValidationSummary) {
	if jsonOutput {
		jsonData, _ := json.MarshalIndent(summary, "", "  ")
		fmt.Println(string(jsonData))
		return
	}

	// Text output
	fmt.Println()
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║            GraphQL Query Validation Results                  ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	fmt.Println()

	for _, result := range summary.Results {
		if result.Passed {
			fmt.Printf("  ✓ PASS  %-40s %4dms\n", result.Name, result.Duration)
		} else {
			fmt.Printf("  ✗ FAIL  %-40s %4dms\n", result.Name, result.Duration)
			for _, err := range result.Errors {
				fmt.Printf("          └─ %s\n", err)
			}
		}
	}

	fmt.Println()
	fmt.Println("──────────────────────────────────────────────────────────────────")

	if summary.Failed == 0 {
		fmt.Printf("  ✓ All %d queries passed validation\n", summary.Total)
	} else {
		fmt.Printf("  Summary: %d total, %d passed, %d failed\n",
			summary.Total, summary.Passed, summary.Failed)
	}
	fmt.Println()
}
