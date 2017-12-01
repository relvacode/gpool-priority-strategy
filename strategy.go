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

// Age returns the age of the given job.
func Age(job *gpool.JobStatus) float64 {
	if job.QueuedOn == nil {
		return 0
	}
	return time.Since(*job.QueuedOn).Seconds()
}

// Priority returns the priority of the given job.
func Priority(job *gpool.JobStatus) float64 {
	if pr, ok := job.Job().(Prioritizer); ok {
		return max(float64(pr.Priority()), 1)
	}
	return 1
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

// Strategy is an evaluater for gpool that executes the highest priority job
type Strategy struct {
	AgeFactor      float64
	PriorityFactor float64
}

// Priority calculates the priority of a particular job
func (sg *Strategy) Priority(job *gpool.JobStatus, maxAge, maxPriority float64) float64 {
	var priority = (Priority(job) / maxPriority) * max(sg.PriorityFactor, 1)
	var age = (Age(job) / maxAge) * max(sg.AgeFactor, 1)
	return age + priority
}

// Evaluate returns the index of the next job to be scheduled based on the highest priority
func (sg *Strategy) Evaluate(jobs []*gpool.JobStatus) (int, bool) {
	if len(jobs) == 0 {
		return 0, false
	}

	// Find the oldest job
	var maxAge float64
	var maxPri float64
	for i := 0; i < len(jobs); i++ {
		age := Age(jobs[i])
		if age > maxAge {
			maxAge = age
		}
		pri := Priority(jobs[i])
		if pri > maxPri {
			maxPri = pri
		}
	}
	maxAge = max(maxAge, 1)
	maxPri = max(maxPri, 1)

	priorities := make([]float64, len(jobs))
	for i := 0; i < len(jobs); i++ {
		priorities[i] = sg.Priority(jobs[i], maxAge, maxPri)
	}
	var pmax float64
	var ppos int
	for i := 0; i < len(priorities); i++ {
		if priorities[i] > pmax {
			pmax = priorities[i]
			ppos = i
		}
	}
	return ppos, true
}
