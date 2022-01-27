package printer

import (
	"__TestSrc/gRpc_Study/Src1/client/cmd"
	"time"
)

type Report struct {
	Name    string `json:"name,omitempty"`
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

func (r *Report) MakeReport(name string, cmd *cmd.Config, date time.Time) {
	r.Name = name
	r.Address = cmd.Host
	r.Port = cmd.Port
	r.Date = date.Format("2006-01-02")
}
