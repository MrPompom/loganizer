package analyzer

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	"go_loganizer/internal/config"
)

type ReportEntry struct {
	log_id        string
	file_path     string
	status        string
	message       string
	error_details string
}
type CheckResult struct {
	InputTarget config.InputTarget
	Status      string
	Message     string
	Err         error
}

func CheckLog(target config.InputTarget) CheckResult { // TO DO : modifier pour check les fichiers
	// Vérifier si le fichier existe et est lisible
	_, err := os.Stat(target.Path)
	if os.IsNotExist(err) {
		return CheckResult{
			InputTarget: target,
			Err:         &NonExistingFileError{Path: target.Path, Err: err},
			Message:     "Fichier introuvable.",
		}
	}

	// Simuler l'analyse avec un time.Sleep aléatoire (50 à 200 ms)
	sleepDuration := time.Duration(rand.Intn(150)+50) * time.Millisecond
	time.Sleep(sleepDuration)

	// Erreur aléatoire simulée : 10% de chance de générer une erreur de parsing
	if rand.Float64() < 0.1 {
		return CheckResult{
			InputTarget: target,
			Status:      "FAILED",
			Message:     "Erreur de parsing simulée.",
			Err:         errors.New("erreur de parsing simulée"),
		}
	}

	return CheckResult{
		InputTarget: target,
		Status:      "SUCCESS",
		Message:     "Analyse terminée avec succès.",
	}
}

func ConvertToReportEntry(res CheckResult) ReportEntry {
	report := ReportEntry{
		log_id:        res.InputTarget.ID,
		file_path:     res.InputTarget.Path,
		status:        res.Status,
		message:       res.Message,
		error_details: "",
	}
	if res.Err != nil {
		var nonExistingFileErr *NonExistingFileError
		var parsingErr *ParsingError
		switch {
		case errors.As(res.Err, &nonExistingFileErr):
			report.status = "FAILED"
			report.message = "Fichier introuvable ou inaccessible."
			report.error_details = fmt.Sprintf("open %v: %v", nonExistingFileErr.Path, nonExistingFileErr.Err)
		case errors.As(res.Err, &parsingErr):
			report.status = "FAILED"
			report.message = "Erreur de parsing."
			report.error_details = fmt.Sprintf("Erreur de parsing: %v", parsingErr.Err)
		default:
			report.status = "FAILED"
			report.message = "Erreur lors de l'analyse."
			report.error_details = fmt.Sprintf("Erreur: %v", res.Err)
		}
	}
	return report
}
