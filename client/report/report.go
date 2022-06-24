package report

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

type AuthSuccess struct {
}
type AuthError struct {
}

type Report struct {
	WinnerID   string
	LoserID    string
	ReporterID string
}

func ReportResult(us string, them string, win bool) {

	fmt.Println("Attempting to report game result...")

	var r Report
	if win {
		r = Report{WinnerID: us, LoserID: them, ReporterID: us}
	} else {
		r = Report{WinnerID: them, LoserID: us, ReporterID: us}
	}

	client := resty.New()
	resp, err := client.R().
		SetBody(r).
		SetResult(&AuthSuccess{}).
		SetError(&AuthError{}).
		Post("http://localhost:5000/gameResult")

	printOutput(resp, err)
}

func printOutput(resp *resty.Response, err error) {
	fmt.Println(resp, err)
}
