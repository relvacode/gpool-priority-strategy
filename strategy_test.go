package strategy

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/relvacode/gpool"
)

var (
	_ Prioritizer = TestPriorityJob{}
	_ gpool.Job   = TestPriorityJob{}
)

type TestPriorityJob struct {
	id int
	p  int
}

func (job TestPriorityJob) Header() fmt.Stringer {
	return gpool.Header(fmt.Sprint(job.id))
}

func (TestPriorityJob) Abort(error) {}

func (TestPriorityJob) Run(context.Context) error {
	return nil
}

func (job TestPriorityJob) Priority() int {
	return job.p
}

type StrategyTest struct {
	// QueuedOffset is the offset to apply to the job from the start date of the test
	QueuedOffset time.Duration

	// Priority is the priority to set for the job
	Priority int
}

func (tc StrategyTest) JobStatus(i int, from time.Time) *gpool.JobStatus {
	job := TestPriorityJob{
		id: i,
		p:  tc.Priority,
	}
	status := gpool.NewJobStatus(job, context.Background(), gpool.Queued)

	t := from.Add(tc.QueuedOffset)
	status.QueuedOn = &t
	return status
}

type StrategyTestCase struct {
	Name  string
	Queue []StrategyTest
	Pick  string

	AgeFactor      float64
	PriorityFactor float64
}

func (tc StrategyTestCase) Run(t *testing.T) {
	ts := time.Now()
	jobs := make([]*gpool.JobStatus, len(tc.Queue))
	for i := 0; i < len(tc.Queue); i++ {
		jobs[i] = tc.Queue[i].JobStatus(i, ts)
	}
	strategy := Strategy{
		AgeFactor:      tc.AgeFactor,
		PriorityFactor: tc.PriorityFactor,
	}
	picked, ok := strategy.Evaluate(jobs)
	if !ok {
		t.Fatal("Strategy did not return a valid job index")
	}

	job := jobs[picked]
	header := job.Job().Header().String()
	if header != tc.Pick {
		t.Fatalf("Wrong job was picked: Wanted job %s; got %s", tc.Pick, header)
	}
}

func TestStrategy_Evaluate(t *testing.T) {
	cases := []StrategyTestCase{
		{
			Name:           "PickOldest",
			AgeFactor:      1000,
			PriorityFactor: 1000,
			Pick:           "2",
			Queue: []StrategyTest{
				{
					QueuedOffset: -time.Second,
				},
				{
					QueuedOffset: -(time.Second * 2),
				},
				{
					QueuedOffset: -(time.Second * 10),
				},
				{
					QueuedOffset: -time.Second,
				},
			},
		},
		{
			Name:           "PickTopPriority",
			AgeFactor:      1000,
			PriorityFactor: 2000,
			Pick:           "0",
			Queue: []StrategyTest{
				{
					QueuedOffset: -time.Second,
					Priority:     2,
				},
				{
					QueuedOffset: -(time.Second * 2),
				},
				{
					QueuedOffset: -(time.Second * 10),
				},
				{
					QueuedOffset: -time.Second,
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, tc.Run)
	}
}
