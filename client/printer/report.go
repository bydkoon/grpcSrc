package printer

import (
	"Src1/client/runner"
	"fmt"
	"time"
)

type Report struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
	Date    string `json:"date"`

	TrackReport TrackReport
}

type TrackReport struct {
	track map[string]time.Duration
}

func (t *TrackReport) TrackCheck(name string, d time.Duration) {
	t.track[name] = d
	//t.Tracks = append(t.Tracks, track)
}

func (t *TrackReport) errorCheck(name string, d time.Duration) {
	t.track[name] = d
	//t.Tracks = append(t.Tracks, track)
}

func (t *TrackReport) Init() {
	t.track = make(map[string]time.Duration)
}

func (r *Report) MakeReport(cfg *runner.Config) {

	r.Address = fmt.Sprintf("%s:%v", cfg.Host, cfg.Port)
	r.Port = cfg.Port
	r.Date = time.Now().Format("2006-01-02")
}
