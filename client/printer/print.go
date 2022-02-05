package printer

import (
	"Src1/client/runner"
	"bytes"
	"fmt"
	"strings"
	"time"
)

func (rp *ReportPrinter) Print() {
	var config runner.Config
	BasicReport(config)
	summaryReport(rp)

	//rp.C.BlockMode,
	//rp.C.TotalRequest,
	//rp.C.Cert)
	//outputTmpl := defaultTmpl
	//
	//return rp.print(buf.String())
}

func BasicReport(config runner.Config) {
	date := time.Now().Format("2006-01-02")
	fmt.Printf(
		DefaultTmpl,
		config.Host,
		config.Port,
		date,
		config.BlockMode,
		config.TotalRequest,
		config.Cert,
	)

}
func summaryReport(rp *ReportPrinter) {

	fmt.Printf(SummaryTmpl,
		rp.Report.MainWorker.TotalCount,
		rp.Report.MainWorker.Slowest,
		rp.Report.MainWorker.Fastest,
		rp.Report.MainWorker.Average,
		rp.Report.MainWorker.Rps,
		histogramPrintString(rp.Report.MainWorker.Histogram),
		LatencyDistributionString(rp.Report.MainWorker.LatencyDistribution),
	)

}

//report := fmt.Sprintf("address: %s \n", r.Address)
//report += fmt.Sprintf("port: %d \n", r.Port)
//report += fmt.Sprintf("date: %v \n", r.Date)
//report += "---------------------------------------------- \n"
//report += fmt.Sprint("Options \n")
//report += fmt.Sprintf("blockMode : %v \n", r.BlockMode)
//report += fmt.Sprintf("timeOut: %vs \n", r.TimeOut)
//report += fmt.Sprintf("TotalCount: %v \n", r.TotalCount)
//report += fmt.Sprintf("SkipTLSVerify : %v \n", r.SkipTLSVerify)
//report += fmt.Sprintf("CertPem: %v \n", r.CertPem)
//report += "---------------------------------------------- \n"
//
//if len(r.TrackReport.track) > 0 {
//	report += "  track name           | Duration           | \n"
//	report += "---------------------------------------------- \n"
//	for k, v := range r.TrackReport.track {
//		report += fmt.Sprintf("  %-18s   |  %-20v \n", k, v)
//	}
//}
//if err != nil {
//	report += fmt.Sprintf(" Error \n")
//	report += fmt.Sprintf(" message : %s \n", err.ErrorString())
//	report += fmt.Sprintf(" step : %s \n", err.ErrorStep())
//	report += fmt.Sprintf(" error : %s \n", err.Error())
//}

//report += "  Address            | Date           | Count      | Average             | Fastest             | Slowest             |\n" +
//	"------------------------------------------------------------------------------------\n"
//
//report += fmt.Sprintf("% -20s | % -12s | % -20s | % -8d | % -20v |% -20v |% -20v |\n",
//	"",
//	r.Address,
//	r.Date,
//)
//report += "------------------------------------------------------------------------------------\n"

//fmt.Printf(report)

const (
	barChar = "⁕"
)

func LatencyDistributionString(latencyDistribution []runner.LatencyDistribution) string {
	report := ""
	for _, o := range latencyDistribution {
		report += fmt.Sprintf("  %v  in %v\n", o.Percentage, o.Latency)
	}
	return report
}

func histogramPrintString(buckets []runner.Bucket) string {
	maxMark := 0.0
	maxCount := 0
	for _, b := range buckets {
		if v := b.Mark; v > maxMark {
			maxMark = v
		}
		if v := b.Count; v > maxCount {
			maxCount = v
		}
	}

	formatMark := func(mark float64) string {
		return fmt.Sprintf("%.3f", mark*1000)
	}
	formatCount := func(count int) string {
		return fmt.Sprintf("%v", count)
	}

	maxMarkLen := len(formatMark(maxMark))
	maxCountLen := len(formatCount(maxCount))
	res := new(bytes.Buffer)
	for i := 0; i < len(buckets); i++ {
		// Normalize bar lengths.
		var barLen int
		if maxCount > 0 {
			barLen = (buckets[i].Count*40 + maxCount/2) / maxCount
		}
		markStr := formatMark(buckets[i].Mark)
		countStr := formatCount(buckets[i].Count)
		res.WriteString(fmt.Sprintf(
			"  %s%s [%v]%s |%v\n",
			markStr,
			strings.Repeat(" ", maxMarkLen-len(markStr)),
			countStr,
			strings.Repeat(" ", maxCountLen-len(countStr)),
			strings.Repeat(barChar, barLen),
		))
	}

	return res.String()
}
