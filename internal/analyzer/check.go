package analyzer

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
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
	Err         error
}

func CheckLog(target config.InputTarget) CheckResult { // TO DO : modifier pour check les fichiers
	client := http.Client{
		Timeout: time.Duration(rand.Intn(150)+50) * time.Millisecond, // Set a timeout random between 50 and 200 ms
	}

	resp, err := client.Get(target.Path)
	if err != nil {
		return CheckResult{
			InputTarget: target,
			Err:         &NonExistingFileError{Path: target.Path, Err: err},
		}
	}
	defer resp.Body.Close()
	return CheckResult{
		InputTarget: target,
		Status:      resp.Status,
	}
}

func ConvertToReportEntry(res CheckResult) ReportEntry {
	report := ReportEntry{
		log_id:        res.InputTarget.ID,
		file_path:     res.InputTarget.Path,
		status:        res.Status,
		message:       "Analyse terminée avec succès",
		error_details: "",
	}
	if res.Err != nil {
		var NonExistingFile *NonExistingFileError
		if errors.As(res.Err, &NonExistingFile) {
			report.status = "FAILED"
			report.message = "Fichier introuvable."
			report.error_details = fmt.Sprintf("open %v: %v", NonExistingFile.Path, NonExistingFile.Err)
		} else {
			report.status = "FAILED"
			report.message = "Erreur lors de l'analyse."
			report.error_details = fmt.Sprintf("Erreur: %v", res.Err)
		}
	}
	return report
}
