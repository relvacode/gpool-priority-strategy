# gpool-priority-strategy
Priority based scheduling strategy for gpool

For use with the Go concurrent execution toolkit [gpool](https://github.com/relvacode/gpool)

## Usage

```go
sg := strategy.Strategy{
    AgeFactor: 1000,
    PriorityFactor: 2000,
}
br := gpool.NewSimpleBridge(5, sg)
p := gpool.New(true, br)
```