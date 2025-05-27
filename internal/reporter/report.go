package reporter

import (
	"encoding/json"
	"fmt"
	"os"

	"go_loganizer/internal/checker"
)

func ExportResultToJsonFile(FilePath string, result []checker.ReportEntry) error {
	// Convert the result to JSON format
	data, err := json.MarshalIndent(result, "", " ")
	if err != nil {
		return fmt.Errorf("failed to marshal results to JSON: %w", err)
	}

	// Write the JSON data to the specified file
	if err := os.WriteFile(FilePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write results to file %s: %w", FilePath, err)
	}

	return nil
}
