package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	showFullPath bool
)

// QueryInfo represents information about a query file
type QueryInfo struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	HasVars     bool   `json:"has_variables"`
	VarsFile    string `json:"variables_file,omitempty"`
	SizeBytes   int64  `json:"size_bytes"`
	Description string `json:"description,omitempty"`
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available GraphQL query files",
	Long: `List all GraphQL query files in the specified directory.

This command scans the queries directory and displays information
about each .graphql file found, including whether it has an
accompanying variables file.

Examples:
  # List all queries in default directory
  gql-validate list

  # List queries in a specific directory
  gql-validate list -q ./my-queries

  # Show full file paths
  gql-validate list --full-path

  # Output as JSON
  gql-validate list -j`,
	RunE: runList,
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(&queriesDir, "queries", "q", "./queries", "directory containing GraphQL query files")
	listCmd.Flags().BoolVar(&showFullPath, "full-path", false, "show full file paths")
}

func runList(cmd *cobra.Command, args []string) error {
	// Check if directory exists
	if _, err := os.Stat(queriesDir); os.IsNotExist(err) {
		return fmt.Errorf("queries directory not found: %s", queriesDir)
	}

	// Find all query files
	var queries []QueryInfo

	err := filepath.Walk(queriesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".graphql") {
			query := QueryInfo{
				Name:      info.Name(),
				Path:      path,
				SizeBytes: info.Size(),
			}

			// Check for corresponding JSON file
			jsonFile := strings.TrimSuffix(path, ".graphql") + ".json"
			if _, err := os.Stat(jsonFile); err == nil {
				query.HasVars = true
				query.VarsFile = jsonFile
			}

			// Try to extract description from first comment line
			if content, err := os.ReadFile(path); err == nil {
				query.Description = extractDescription(string(content))
			}

			queries = append(queries, query)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to scan directory: %w", err)
	}

	if len(queries) == 0 {
		fmt.Printf("No GraphQL query files found in: %s\n", queriesDir)
		return nil
	}

	// Output results
	if jsonOutput {
		return printListJSON(queries)
	}

	return printListText(queries)
}

func extractDescription(content string) string {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") {
			// Remove the # and leading space
			desc := strings.TrimPrefix(line, "#")
			desc = strings.TrimSpace(desc)
			if desc != "" {
				return desc
			}
		} else if line != "" {
			// Stop at first non-comment, non-empty line
			break
		}
	}
	return ""
}

func printListJSON(queries []QueryInfo) error {
	output := map[string]interface{}{
		"directory":   queriesDir,
		"total_files": len(queries),
		"queries":     queries,
	}

	jsonData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(jsonData))
	return nil
}

func printListText(queries []QueryInfo) error {
	fmt.Println()
	fmt.Printf("GraphQL Queries in: %s\n", queriesDir)
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	for i, q := range queries {
		displayPath := q.Name
		if showFullPath {
			displayPath = q.Path
		}

		fmt.Printf("  %d. %s\n", i+1, displayPath)

		if q.Description != "" {
			fmt.Printf("     │ %s\n", q.Description)
		}

		if q.HasVars {
			varsDisplay := filepath.Base(q.VarsFile)
			if showFullPath {
				varsDisplay = q.VarsFile
			}
			fmt.Printf("     └─ Variables: %s\n", varsDisplay)
		}

		if verbose {
			fmt.Printf("     └─ Size: %d bytes\n", q.SizeBytes)
		}

		fmt.Println()
	}

	fmt.Printf("Total: %d query file(s)\n", len(queries))

	// Count files with variables
	withVars := 0
	for _, q := range queries {
		if q.HasVars {
			withVars++
		}
	}
	if withVars > 0 {
		fmt.Printf("       %d with variables file(s)\n", withVars)
	}

	fmt.Println()
	return nil
}
