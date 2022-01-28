package printer

import (
	errors2 "Src1/client/errors"
	"Src1/client/runner"
	"fmt"
)

type ReportPrinter struct {
	Report *printer.Report
}

func (r *Report) PrinterReport(err *errors2.Errors, cfg *runner.Config) {

	report := fmt.Sprintf("address: %s \n", r.Address)
	report += fmt.Sprintf("port: %d \n", r.Port)
	report += fmt.Sprintf("date: %v \n", r.Date)
	report += "---------------------------------------------- \n"
	report += fmt.Sprint("Options \n")
	report += fmt.Sprintf("blockMode : %v \n", cfg.BlockMode)
	report += fmt.Sprintf("timeOut: %vs \n", cfg.TimeOut)
	report += fmt.Sprintf("TotalCount: %v \n", cfg.TotalCount)
	report += fmt.Sprintf("SkipTLSVerify : %v \n", cfg.SkipTLSVerify)
	report += fmt.Sprintf("CertPem: %v \n", cfg.CertPem)
	report += "---------------------------------------------- \n"

	if len(r.TrackReport.track) > 0 {
		report += "  track name           | Duration           | \n"
		report += "---------------------------------------------- \n"
		for k, v := range r.TrackReport.track {
			report += fmt.Sprintf("  %-18s   |  %-20v \n", k, v)
		}
	}
	if err != nil {
		report += fmt.Sprintf(" Error \n")
		report += fmt.Sprintf(" message : %s \n", err.ErrorString())
		report += fmt.Sprintf(" step : %s \n", err.ErrorStep())
		report += fmt.Sprintf(" error : %s \n", err.Error())
	}

	//report += "  Address            | Date           | Count      | Average             | Fastest             | Slowest             |\n" +
	//	"------------------------------------------------------------------------------------\n"
	//
	//report += fmt.Sprintf("% -20s | % -12s | % -20s | % -8d | % -20v |% -20v |% -20v |\n",
	//	"",
	//	r.Address,
	//	r.Date,
	//)
	//report += "------------------------------------------------------------------------------------\n"

	fmt.Printf(report)

}
