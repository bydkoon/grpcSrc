package printer

import (
	"time"
)

type Duration time.Duration

type Report struct {
	Address       string   `json:"address"`
	Port          int      `json:"port"`
	Date          string   `json:"date"`
	BlockMode     bool     `json:"blockMode"`
	SkipTLSVerify bool     `json:"skipTLS"`
	CertPem       string   `json:"certPem"`
	TimeOut       Duration `json:"timeOut"`
	TotalCount    uint     `json:"TotalCount"`
	TrackReport   TrackReport
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
	t.track = make(
		map[string]time.Duration)
}

//func (r *Re) MakeReport(cfg *runner.Config) {
//
//	r.Address = fmt.Sprintf("%s:%v", cfg.Host, cfg.Port)
//	r.Port = cfg.Port
//	r.Date = time.Now().Format("2006-01-02")
//	r.BlockMode = cfg.BlockMode
//	r.SkipTLSVerify = cfg.SkipTLSVerify
//	r.TimeOut = Duration(cfg.TimeOut)
//	r.CertPem = cfg.Cert
//	r.TotalCount = cfg.TotalRequest

//}
