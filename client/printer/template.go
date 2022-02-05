package printer

var (
	DefaultTmpl = `
Base
	address: %s
	port: %d
	date: %v
Options
	Block: %v
	totalCount: %v
	cert: %s
`

	SummaryTmpl = `
Summary:
Total   : %v
		Slowest : %v
		Fastest : %v
		Average : %v
		Request / sec : %v
	
	Response time Histogram :
%s
	Status Distribute :
%s
	Error
%s
`
)
