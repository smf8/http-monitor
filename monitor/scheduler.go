package monitor

import (
	"errors"
	"fmt"
	"github.com/gammazero/workerpool"
	"time"
)

type Scheduler struct {
	Mnt  *Monitor
	Quit chan struct{}
}

// NewScheduler creates a new scheduler instance with mnt as monitor
// it also creates a quit signal channel for emergency exits
func NewScheduler(mnt *Monitor) (*Scheduler, error) {
	sch := &Scheduler{Quit: make(chan struct{})}
	if mnt != nil {
		sch.Mnt = mnt
		return sch, nil
	}
	return nil, errors.New("cannot create a scheduler with nil monitor")
}

// DoWithIntervals creates a ticker to the execute mnt.Do() every d duration
// it listens to a quit channel as well for termination signal.
// in order to stop it, Call StopSchedule().
func (sch *Scheduler) DoWithIntervals(d time.Duration) {
	ticker := time.NewTicker(d)
	go func() {
		for {
			select {
			case <-ticker.C:
				sch.Mnt.Do()
			case <-sch.Quit:
				// stopping worker pool from accepting anymore jobs
				err := sch.Mnt.Cancel()

				if err != nil {
					fmt.Println("error canceling monitor on quit signal in DoWithIntervals()")
				}

				// since out mnt's worker pool is useless after cancel we instantiate another one
				sch.Mnt.wp = workerpool.New(sch.Mnt.workerSize)

				ticker.Stop()
				return
			}
		}
	}()
}

// StopSchedule simply closes sch.Quit channel in order to stop it's running schedule
func (sch *Scheduler) StopSchedule() {
	close(sch.Quit)
}
