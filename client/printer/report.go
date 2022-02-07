package printer

import (
	"Src1/client/runner"
)

type ReportPrinter struct {
	//Out    io.Writer
	Report *runner.Reporter
	Config *runner.RunConfig
}
