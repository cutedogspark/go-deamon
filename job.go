package gworker

import "time"

type Job struct {
	Name  string
	Delay time.Duration
}
