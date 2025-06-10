package reporter

import (
	"encoding/json"
	"fmt"
	"go_loganizer/internal/analyzer"
	"os"
)

func ExportResultToJsonFile(FilePath string, result []analyzer.ReportEntry) error {
	// Convert the result to JSON format
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal results to JSON: %w", err)
	}

	if err := os.WriteFile(FilePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write results to file %s: %w", FilePath, err)
	}
	return nil
}
