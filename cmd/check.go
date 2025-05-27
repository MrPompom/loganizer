package cmd

import (
	"errors"
	"fmt"
	"go_loganizer/internal/analyzer"
	"go_loganizer/internal/config"
	"go_loganizer/internal/reporter"
	"sync"

	"github.com/spf13/cobra"
)

var (
	inputFilePath  string
	outputFilePath string
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyse une liste de log",
	Long:  `La commande 'analyze' permet d'analyser un fichier de log et d'en extraire des informations pertinentes.`,
	Run: func(cmd *cobra.Command, args []string) {

		if inputFilePath == "" {
			fmt.Println("Erreur: le chemin du fichier d'entrée (--input) est obligatoire.")
			return
		}

		targets, err := config.LoadTargetsFromFile(inputFilePath)
		if err != nil {
			fmt.Printf("Erreur lors du chargement des fichiers de log: %v\n", err)
			return
		}

		if len(targets) == 0 {
			fmt.Println("Aucune Log à vérifier trouvée dans le fichier d'entrée.")
			return
		}
		var wg sync.WaitGroup
		resultsChan := make(chan analyzer.CheckResult, len(targets))
		wg.Add(len(targets))
		for _, target := range targets {
			go func(t config.InputTarget) {
				result := analyzer.CheckLog(t)
				resultsChan <- result
				defer wg.Done()
			}(target)
		}
		wg.Wait()
		close(resultsChan)

		var finalReport []analyzer.ReportEntry
		for res := range resultsChan {
			reportEntry := analyzer.ConvertToReportEntry(res)
			finalReport = append(finalReport, reportEntry)

			status := "OK"
			msg := res.Message
			errMsg := ""
			if res.Err != nil {
				status = "KO"
				var fileErr *analyzer.NonExistingFileError
				var parseErr *analyzer.ParsingError
				switch {
				case errors.As(res.Err, &fileErr):
					errMsg = fmt.Sprintf("Fichier introuvable/inaccessible: %v", fileErr.Err)
				case errors.As(res.Err, &parseErr):
					errMsg = fmt.Sprintf("Erreur de parsing: %v", parseErr.Err)
				default:
					errMsg = fmt.Sprintf("Autre erreur: %v", res.Err)
				}
			}
			fmt.Printf("Résumé: ID=%s | Chemin=%s | Statut=%s | Message=%s", res.InputTarget.ID, res.InputTarget.Path, status, msg)
			if errMsg != "" {
				fmt.Printf(" | Erreur=%s", errMsg)
			}
			fmt.Println()
		}
		if outputFilePath != "" {
			err := reporter.ExportResultToJsonFile(outputFilePath, finalReport)
			if err != nil {
				fmt.Printf("Erreur lors de l'export des résultats vers le fichier %s: %v\n", outputFilePath, err)
			} else {
				fmt.Printf("✅ Résultats exportés avec succès vers le fichier %s\n", outputFilePath)
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(analyzeCmd)
	analyzeCmd.Flags().StringVarP(&inputFilePath, "config", "c", "", "Chemin vers le fichier JSON d'entrée")
	analyzeCmd.Flags().StringVarP(&outputFilePath, "output", "o", "", "Chemin vers le fichier JSON de sortie pour les résultats")
	analyzeCmd.MarkFlagRequired("config")
}
