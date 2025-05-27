package checker

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/MrPompom/gowatcher_tp1/internal/config"
)

type ReportEntry struct {
	Name   string
	URL    string
	Owner  string
	Status string
	ErrMsg string
}
type CheckResult struct {
	InputTarget config.InputTarget
	Status      string
	Err         error
}

func CheckURL(target config.InputTarget) CheckResult {
	client := http.Client{
		Timeout: time.Second * 3, // Set a timeout of 3 seconds
	}

	resp, err := client.Get(target.URL)
	if err != nil {
		return CheckResult{
			InputTarget: target,
			Err:         &UnreachableURLError{URL: target.URL, Err: err},
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
		Name:   res.InputTarget.Name,
		URL:    res.InputTarget.URL,
		Owner:  res.InputTarget.Owner,
		Status: res.Status,
	}
	if res.Err != nil {
		var UnreachableURL *UnreachableURLError
		if errors.As(res.Err, &UnreachableURL) {
			report.Status = "Unreachable"
			report.ErrMsg = fmt.Sprintf("Unreachable URL: %v", UnreachableURL.Err)
		} else {
			report.Status = "Error"
			report.ErrMsg = fmt.Sprintf("Error generique: %v", res.Err)
		}
	}
	return report
}
