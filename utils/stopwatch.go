package utils

import (
	"time"
)

type StopWatch struct {
	uuid         string
	StartTime    time.Time
	trackingTime time.Time
}

//func NewStopWatch() *StopWatch {
//	return NewStopWatchUUID(uuid.NewString())
//}

func (sw *StopWatch) Start() time.Time {
	sw.StartTime = time.Now()
	sw.trackingTime = time.Now()

	//fmt.Printf("# [0s] %s >> Start \n", sw.convertDateFormat(sw.StartTime))
	return sw.StartTime
}

func (sw *StopWatch) Track(flags ...string) time.Duration {
	elapsed := time.Since(sw.trackingTime)
	sw.trackingTime = time.Now()

	//fmt.Printf("# [%v] %s >> Tracking \n", elapsed, flags)
	return elapsed
}

func (sw *StopWatch) Stop() (totalDuration time.Duration, endTime time.Time) {
	elapsed_start := time.Since(sw.StartTime)
	//elapsed_track := time.Since(sw.trackingTime)

	//fmt.Printf("# [%v] %s >> Stop | Total laptime : %v \n", elapsed_track, sw.convertDateFormat(time.Now()), elapsed_start)
	return elapsed_start, time.Now()
}

func (sw *StopWatch) convertDateFormat(time time.Time) string {
	format := "2006-01-02 15:04:05.000"
	return time.Format(format)
}

func (sw *StopWatch) GetStart() (time.Time, time.Time) {
	sw.StartTime = time.Now()
	sw.trackingTime = time.Now()

	return sw.StartTime, sw.trackingTime
}

func (sw *StopWatch) GetTrack(flags ...string) time.Duration {
	elapsed := time.Since(sw.trackingTime)
	sw.trackingTime = time.Now()

	return elapsed
}

func (sw *StopWatch) GetCurTotal(flags ...string) time.Duration {
	return time.Since(sw.StartTime)
}

func (sw *StopWatch) GetStop() (time.Duration, time.Duration, time.Time) {
	elapsed_start := time.Since(sw.StartTime)
	elapsed_track := time.Since(sw.trackingTime)

	return elapsed_start, elapsed_track, time.Now()
}

func NewStopWatchUUID(uid string) *StopWatch {
	return &StopWatch{
		uuid:      uid,
		StartTime: time.Now(),
	}
}
