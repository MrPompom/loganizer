package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd représente la commande de base lorsqu'elle est appelée sans sous-commande.
// C'est une variable globale (dans le package cmd) qui est un pointeur vers une instance de cobra.Command.
var rootCmd = &cobra.Command{
	Use:   "gowatcher",
	Short: "Gowatcher est un outil pour vérifier l'accessibilité des URLs.",
	Long:  "Un outil CLI en Go pour vérifier l'état d'URLs, gérer la concurrence et exporter les résultats.`",
}

// Execute ajoute toutes les commandes enfants à la commande racine et définit également les drapeaux.
// C'est la fonction principale du package `cmd` qui doit être appelée depuis `main.go`.
// Elle est responsable de l'analyse des arguments de la ligne de commande et de l'exécution de la commande appropriée.
func Execute() {
	// Tente d'exécuter la commande racine. Cobra va automatiquement analyser les arguments
	// fournis sur la ligne de commande et déclencher la fonction `Run` de la commande ou sous-commande correspondante.
	if err := rootCmd.Execute(); err != nil {
		// Si une erreur se produit pendant l'exécution (par exemple, argument invalide, commande inconnue),
		// l'erreur est affichée sur la sortie d'erreur standard (stderr).
		fmt.Fprintf(os.Stderr, "Erreur: %v\n", err)
		// Le programme se termine avec un code d'erreur non nul (1) pour indiquer un échec.
		os.Exit(1)
	}
}

// init() est une fonction spéciale en Go qui est exécutée automatiquement
// une fois que tous les imports du package ont été initialisés.
// Elle est souvent utilisée pour configurer des éléments avant que le code principal ne s'exécute.
func init() {
	// Ici, vous pouvez définir des drapeaux (flags) persistants.
	// Un drapeau persistant est un drapeau qui est disponible pour la commande racine
	// et toutes ses sous-commandes.
	// Exemple (commenté) : rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "fichier de configuration (par défaut est $HOME/.gowatcher.yaml)")
	// Pour notre cas, la commande racine n'a pas de drapeaux persistants pour l'instant.
}
