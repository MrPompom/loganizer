package cmd

import (
	"errors"
	"fmt"
	"sync"

	"github.com/MrPompom/gowatcher_tp1/internal/checker"
	"github.com/MrPompom/gowatcher_tp1/internal/config"
	"github.com/MrPompom/gowatcher_tp1/internal/reporter"
	"github.com/spf13/cobra"
)

var (
	inputFilePath  string
	outputFilePath string
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Vérifie l'accessibilité d'une liste d'URLs.",
	Long:  `La commande 'check' parcourt une liste prédéfinie d'URLs et affiche leur statut d'accessibilité en utilisant des goroutines pour la concurrence.`,
	Run: func(cmd *cobra.Command, args []string) {

		if inputFilePath == "" {
			fmt.Println("Erreur: le chemin du fichier d'entrée (--input) est obligatoire.")
			return
		}

		// Charger les "cibles" depuis le fichier JSON d'entrée
		targets, err := config.LoadTargetsFromFile(inputFilePath)
		if err != nil {
			fmt.Printf("Erreur lors du chargement des URLs: %v\n", err)
			return
		}

		if len(targets) == 0 {
			fmt.Println("Aucune URL à vérifier trouvée dans le fichier d'entrée.")
			return
		}
		// Compteur de goroutine en attente
		var wg sync.WaitGroup
		resultsChan := make(chan checker.CheckResult, len(targets)) // Canal pour collecter les résultats
		// On initialise/compte le nombre de goroutines attendues
		wg.Add(len(targets))
		for _, target := range targets {
			// On lance une fonction annonyme qui prend en paramètre une copie de url
			go func(t config.InputTarget) {
				result := checker.CheckURL(t)
				resultsChan <- result // On envoie le resultat au canal
				// Garantit qu'à la fin de la fonction, le compteur wg sera décrémenté de 1, `
				// signalant que la Go routine est terminée
				defer wg.Done()
			}(target)
		}
		// Bloque l'exécution du main() jusqu'à ce que toutes les goroutines aient appelé wg.Done()
		wg.Wait()
		close(resultsChan)

		var finalReport []checker.ReportEntry
		for res := range resultsChan {
			reportEntry := checker.ConvertToReportEntry(res)
			finalReport = append(finalReport, reportEntry)

			if res.Err != nil {
				var unreachable *checker.UnreachableURLError
				if errors.As(res.Err, &unreachable) {
					fmt.Printf("KO %s (%s) est inaccessible : %v\n", res.InputTarget.Name, unreachable.URL, unreachable.Err)
				} else {
					fmt.Printf("KO %s (%s) : erreur - %v\n", res.InputTarget.Name, res.InputTarget.URL, res.Err)
				}
			} else {
				fmt.Printf("OK %s (%s) : OK - %s\n", res.InputTarget.Name, res.InputTarget.URL, res.Status)
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
	// elle "ajoute" la sous-commande `checkCmd` à la commande racine `rootCmd`
	// C'est ainsi que Cobra sait que 'check' est une commande valide sous 'gowatcher'.
	rootCmd.AddCommand(checkCmd)

	checkCmd.Flags().StringVarP(&inputFilePath, "input", "i", "", "Chemin vers le fichier JSON d'entrée contenant les URLS")

	checkCmd.Flags().StringVarP(&outputFilePath, "output", "o", "", "Chemin vers le fichier JSON de sortie pour les résultats")

	// Marquer le drapeau "input" comme obligatoire
	checkCmd.MarkFlagRequired("input")

	checkCmd.MarkFlagRequired("output")
}
