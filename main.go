package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dosco/graphjin/core"
	_ "github.com/jackc/pgx/v5/stdlib"
	"gopkg.in/yaml.v2"
)

// TestResult represents the result of validating a single query
type TestResult struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	Passed   bool   `json:"passed"`
	Error    string `json:"error,omitempty"`
	Duration int64  `json:"duration_ms"`
}

// Config holds the configuration for GraphJin
type Config struct {
	Database struct {
		Type     string `yaml:"type"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		DBName   string `yaml:"dbname"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		SSLMode  string `yaml:"sslmode"`
	} `yaml:"database"`
	Production bool `yaml:"production"`
}

func main() {
	// Parse command line flags
	configPath := flag.String("config", "config.yaml", "Path to config file")
	queriesDir := flag.String("queries", "./queries", "Path to queries directory")
	verbose := flag.Bool("verbose", false, "Enable verbose output")
	outputFormat := flag.String("format", "text", "Output format: text or json")
	flag.Parse()

	// Load configuration
	config, err := loadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize GraphJin
	gj, err := initializeGraphJin(config)
	if err != nil {
		log.Fatalf("Failed to initialize GraphJin: %v", err)
	}

	// Find all query files
	queryFiles, err := findQueryFiles(*queriesDir)
	if err != nil {
		log.Fatalf("Failed to find query files: %v", err)
	}

	if len(queryFiles) == 0 {
		fmt.Println("No query files found in", *queriesDir)
		os.Exit(0)
	}

	// Test all queries
	results := testQueries(gj, queryFiles, *verbose)

	// Print results
	printResults(results, *outputFormat)

	// Determine exit code
	exitCode := 0
	for _, result := range results {
		if !result.Passed {
			exitCode = 1
			break
		}
	}

	os.Exit(exitCode)
}

func loadConfig(configPath string) (*Config, error) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("could not read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("could not parse config file: %w", err)
	}

	return &config, nil
}

func initializeGraphJin(config *Config) (*core.GraphJin, error) {
	// Build DSN
	dsn := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=%s",
		config.Database.Host,
		config.Database.Port,
		config.Database.DBName,
		config.Database.User,
		config.Database.Password,
		config.Database.SSLMode,
	)

	// Connect to database
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Create GraphJin configuration
	gjConfig := &core.Config{
		Debug:            true,
		Production:       config.Production,
		DisableAllowList: true,  // This allows all queries without an allow list
		DefaultBlock:     false, // Don't block tables by default
	}

	// Initialize GraphJin
	gj, err := core.NewGraphJin(gjConfig, db)
	if err != nil {
		return nil, fmt.Errorf("failed to create GraphJin instance: %w", err)
	}

	return gj, nil
}

func findQueryFiles(queriesDir string) ([]string, error) {
	var queryFiles []string

	err := filepath.Walk(queriesDir, func(path string, info os.FileInfo, err error) error {
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

func testQueries(gj *core.GraphJin, queryFiles []string, verbose bool) []TestResult {
	var results []TestResult

	for _, queryFile := range queryFiles {
		result := TestResult{
			Name:   filepath.Base(queryFile),
			Path:   queryFile,
			Passed: false,
		}

		start := time.Now()

		// Read query
		query, err := ioutil.ReadFile(queryFile)
		if err != nil {
			result.Error = fmt.Sprintf("Failed to read query file: %v", err)
			result.Duration = time.Since(start).Milliseconds()
			results = append(results, result)
			continue
		}

		// Look for corresponding JSON file
		jsonFile := strings.TrimSuffix(queryFile, ".graphql") + ".json"
		var variables json.RawMessage

		if _, err := os.Stat(jsonFile); err == nil {
			jsonData, err := ioutil.ReadFile(jsonFile)
			if err == nil {
				variables = json.RawMessage(jsonData)
			}
		} else {
			// If no JSON file, use empty object
			variables = json.RawMessage("{}")
		}

		// Execute query
		ctx := context.Background()
		res, err := gj.GraphQL(ctx, string(query), variables, nil)

		result.Duration = time.Since(start).Milliseconds()

		if err != nil {
			result.Error = fmt.Sprintf("GraphQL execution error: %v", err)
		} else if hasNestedErrors(res.Data) {
			result.Error = "Query returned nested errors in data"
		} else {
			result.Passed = true
		}

		if verbose {
			if result.Passed {
				fmt.Printf("✓ PASS: %s (took %dms)\n", result.Name, result.Duration)
			} else {
				fmt.Printf("✗ FAIL: %s - %s\n", result.Name, result.Error)
			}
		}

		results = append(results, result)
	}

	return results
}

func hasNestedErrors(data []byte) bool {
	if len(data) == 0 {
		return false
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return false
	}

	return checkForErrors(result)
}

func checkForErrors(data interface{}) bool {
	switch v := data.(type) {
	case map[string]interface{}:
		// Check for "errors" key
		if errors, ok := v["errors"]; ok {
			if errors != nil {
				return true
			}
		}

		// Check for "error" key
		if errorVal, ok := v["error"]; ok {
			if errorVal != nil && errorVal != "" {
				return true
			}
		}

		// Recursively check nested objects
		for _, value := range v {
			if checkForErrors(value) {
				return true
			}
		}
	case []interface{}:
		// Recursively check arrays
		for _, item := range v {
			if checkForErrors(item) {
				return true
			}
		}
	}

	return false
}

func printResults(results []TestResult, format string) {
	passed := 0
	failed := 0

	for _, result := range results {
		if result.Passed {
			passed++
		} else {
			failed++
		}
	}

	if format == "json" {
		jsonOutput := map[string]interface{}{
			"total":   len(results),
			"passed":  passed,
			"failed":  failed,
			"results": results,
		}

		jsonData, _ := json.MarshalIndent(jsonOutput, "", "  ")
		fmt.Println(string(jsonData))
		return
	}

	// Text output (default)
	fmt.Printf("\nQuery Validation Results:\n")
	fmt.Printf("========================\n")

	for _, result := range results {
		status := "✓ PASS"
		if !result.Passed {
			status = "✗ FAIL"
		}
		fmt.Printf("%s: %s (%dms)\n", status, result.Name, result.Duration)
		if !result.Passed && result.Error != "" {
			fmt.Printf("    Error: %s\n", result.Error)
		}
	}

	fmt.Printf("\nSummary: %d tests, %d passed, %d failed\n", len(results), passed, failed)
}
