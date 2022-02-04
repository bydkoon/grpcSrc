package printer

var (
	SummaryTmpl = `
Summary:
Total   : %-v
		Slowest : %-v
		Fastest : %-v
		Average : %-v
		Request / sec : %v
	
	Response time Histogram :
%s
	Status Distribute :
%s
`
)
