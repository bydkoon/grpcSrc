package printer

import (
	"Src1/client/runner"
	"bytes"
	"fmt"
	"strings"
	"time"
)

func (rp *ReportPrinter) Print() {
	BasicReport(rp)
	summaryReport(rp)

}

func BasicReport(rp *ReportPrinter) {
	date := time.Now().Format("2006-01-02")
	fmt.Printf(
		DefaultTmpl,
		rp.Config.Host,
		rp.Config.Port,
		date,
		rp.Config.Block,
		rp.Config.TotalRequest,
		rp.Config.Cert,
	)

}
func summaryReport(rp *ReportPrinter) {

	fmt.Printf(SummaryTmpl,
		rp.Report.MainWorker.TotalCount,
		rp.Report.MainWorker.Slowest,
		rp.Report.MainWorker.Fastest,
		rp.Report.MainWorker.Average,
		rp.Report.MainWorker.Rps,
		rp.Report.MainWorker.ErrorCount,
		histogramPrintString(rp.Report.MainWorker.Histogram),
		LatencyDistributionString(rp.Report.MainWorker.LatencyDistribution),
		ErrorString(rp.Report.MainWorker.ErrorReport),
	)

}

const (
	barChar = "⬛"
)

func LatencyDistributionString(latencyDistribution []runner.LatencyDistribution) string {
	report := ""
	for _, o := range latencyDistribution {
		report += fmt.Sprintf("  %v  in %v\n", o.Percentage, o.Latency)
	}
	return report
}

func ErrorString(e []runner.Error) string {
	report := ""
	for _, o := range e {
		report += fmt.Sprintf("  %s  , %s,  %v \n", o.Error, o.ErrorCode, o.ErrorMessage)
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
