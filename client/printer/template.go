package printer

var (
	DefaultTmpl = `
Base
	address: %s
	port: %d
	date: %v
Options
	Block: %v
	TotalRequest: %v
	Cert: %s

`

	SummaryTmpl = `
Summary:
Total   : %v
		Slowest : %v
		Fastest : %v
		Average : %v
		Request / sec : %v
		ErrorCount: %v
	
Response time Histogram :
%s
Status Distribute :
%s
Error:
%s
`
)
