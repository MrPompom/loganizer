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
			fmt.Printf("Erreur lors du chargement des URLs: %v\n", err)
			return
		}

		if len(targets) == 0 {
			fmt.Println("Aucune URL à vérifier trouvée dans le fichier d'entrée.")
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

			if res.Err != nil {
				var unreachable *analyzer.NonExistingFileError
				if errors.As(res.Err, &unreachable) {
					fmt.Printf("KO %s (%s) est inaccessible : %v\n", res.InputTarget.ID, unreachable.Path, unreachable.Err)
				} else {
					fmt.Printf("KO %s (%s) : erreur - %v\n", res.InputTarget.ID, res.InputTarget.Path, res.Err)
				}
			} else {
				fmt.Printf("OK %s (%s) : OK - %s\n", res.InputTarget.ID, res.InputTarget.Path, res.Status)
			}
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
	analyzeCmd.MarkFlagRequired("output")
}
