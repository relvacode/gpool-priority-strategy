package strategy

import (
	"time"

	"github.com/relvacode/gpool"
)

// Prioritizer is an interface that should be implemented by a Job
// It should return an integer indicating the initial priority of the job.
type Prioritizer interface {
	Priority() int
}

var _ gpool.ScheduleStrategy = (*Strategy)(nil).Evaluate

// Age returns the age of the given job as an integer
func Age(job *gpool.JobStatus) int {
	if job.QueuedOn == nil {
		return 0
	}
	return int(time.Since(*job.QueuedOn).Nanoseconds())
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Strategy is an evaluater for gpool that executes the highest priority job
type Strategy struct {
	AgeFactor      int
	PriorityFactor int
}

// Priority calculates the priority of a particular job
func (sg *Strategy) Priority(job *gpool.JobStatus, maxAge int) (priority int) {
	if pr, ok := job.Job().(Prioritizer); ok {
		priority += max(pr.Priority(), 1) * sg.PriorityFactor
	}
	priority += (Age(job) / maxAge) * sg.AgeFactor
	return
}

// Evaluate returns the index of the next job to be scheduled based on the highest priority
func (sg *Strategy) Evaluate(jobs []*gpool.JobStatus) (int, bool) {
	if len(jobs) == 0 {
		return 0, false
	}

	// Find the oldest job
	var maxAge int
	for i := 0; i < len(jobs); i++ {
		age := Age(jobs[i])
		if age > maxAge {
			maxAge = age
		}
	}
	maxAge = max(maxAge, 1)

	priorities := make([]int, len(jobs))
	for i := 0; i < len(jobs); i++ {
		priorities[i] = sg.Priority(jobs[i], maxAge)
	}
	var pmax int
	var ppos int
	for i := 0; i < len(priorities); i++ {
		if priorities[i] > pmax {
			pmax = priorities[i]
			ppos = i
		}
	}
	return ppos, true
}
