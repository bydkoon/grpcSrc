package runner

import (
	"sync"
	"time"
)

type SubWorker struct {
	workerID string
	Start    time.Time
	End      time.Time

	StartTime time.Duration
	EndTime   time.Duration

	Error Error
}

type Error struct {
	ErrorCode    string
	ErrorMessage string
	Error        error
}

type MainWorker struct {
	Workers []SubWorker

	Name string    `json:"name,omitempty"`
	Date time.Time `json:"date"`

	StartTime time.Time
	EndTime   time.Time

	TotalCount uint64 `json:"totalCount"`

	TotalDuration time.Duration `json:"totalDuration"`

	Average time.Duration `json:"average"`
	Fastest time.Duration `json:"fastest"`
	Slowest time.Duration `json:"slowest"`

	Rps                 int `json:"rps"`
	lock                sync.Mutex
	LatencyDistribution []LatencyDistribution `json:"latencyDistribution"`
	Histogram           []Bucket              `json:"histogram"`

	ErrorCount  int `json:"errorCount"`
	ErrorReport []Error
}

func newReporter(wID string) *SubWorker {

	return &SubWorker{
		workerID: wID,
	}
}

func (ltr *MainWorker) addSimpleReports(target SubWorker) {
	ltr.lock.Lock()
	ltr.Workers = append(ltr.Workers, target)
	ltr.lock.Unlock()
}
func newLoadReporter() *MainWorker {

	return &MainWorker{}
}

type Reporter struct {
	MainWorker *MainWorker
	config     *RunConfig

	id         string
	totalCount uint64
}
type LatencyDistribution struct {
	Percentage int           `json:"percentage"`
	Latency    time.Duration `json:"latency"`
}

type Bucket struct {
	// The Mark for histogram bucket in seconds
	Mark float64 `json:"mark"`

	// The count in the bucket
	Count int `json:"count"`

	// The frequency of results in the bucket as a decimal percentage
	Frequency float64 `json:"frequency"`
}

type ResultDetail struct {
	Timestamp time.Time     `json:"timestamp"`
	Latency   time.Duration `json:"latency"`
	Error     string        `json:"error"`
	Status    string        `json:"status"`
}

func latencies(latencies []float64) []LatencyDistribution {
	pctls := []int{10, 25, 50, 75, 90, 95, 99}
	data := make([]float64, len(pctls))
	lt := float64(len(latencies))
	for i, p := range pctls {
		ip := (float64(p) / 100.0) * lt
		di := int(ip)

		// since we're dealing with 0th based ranks we need to
		// check if ordinal is a whole number that lands on the percentile
		// if so adjust accordingly
		if ip == float64(di) {
			di = di - 1
		}

		if di < 0 {
			di = 0
		}

		data[i] = latencies[di]
	}

	res := make([]LatencyDistribution, len(pctls))
	for i := 0; i < len(pctls); i++ {
		if data[i] > 0 {
			lat := time.Duration(data[i] * float64(time.Second))
			res[i] = LatencyDistribution{Percentage: pctls[i], Latency: lat}
		}
	}
	return res
}

func histogram(latencies []float64, slowest, fastest float64) []Bucket {
	bc := 10
	buckets := make([]float64, bc+1)
	counts := make([]int, bc+1)
	bs := (slowest - fastest) / float64(bc)
	for i := 0; i < bc; i++ {
		buckets[i] = fastest + bs*float64(i)
	}
	buckets[bc] = slowest
	var bi int
	var max int
	for i := 0; i < len(latencies); {
		if latencies[i] <= buckets[bi] {
			i++
			counts[bi]++
			if max < counts[bi] {
				max = counts[bi]
			}
		} else if bi < len(buckets)-1 {
			bi++
		}
	}
	res := make([]Bucket, len(buckets))
	for i := 0; i < len(buckets); i++ {
		res[i] = Bucket{
			Mark:      buckets[i],
			Count:     counts[i],
			Frequency: float64(counts[i]) / float64(len(latencies)),
		}
	}
	return res
}
