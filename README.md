# gpool-priority-strategy
Priority based scheduling strategy for gpool

For use with the Go concurrent execution toolkit [gpool](https://github.com/relvacode/gpool)
This library enables gpool to schedule jobs based on

  - The age of the job (time since queue)
  - A custom priority of the job
  - Both

The age and priority are weighted using factor values which modify the importance of each attribute.
A factor value of `0` disables that attribute for consideration in scheduling.

## Usage

```go
sg := strategy.Strategy{
    AgeFactor:      1000,
    PriorityFactor: 2000,
}
br := gpool.NewSimpleBridge(5, sg)
p := gpool.New(true, br)
```

You may have your `Job` interface implement `strategy.Prioritizer`. Doing so enables this library to pick the next job to execute based on some priority value.

```go
type LowPriorityJob struct {

}

func (LowPriorityJob) Priority() int {
    return 1
}

type HighPriorityJob struct {

}

func (HighPriorityJob) Priority() int {
    return 5
}
```