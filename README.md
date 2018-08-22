# SafeQueue [![Build Status](https://travis-ci.org/valinurovam/safequeue.svg?branch=master)](https://travis-ci.org/valinurovam/safequeue) [![Coverage Status](https://coveralls.io/repos/github/valinurovam/safequeue/badge.svg)](https://coveralls.io/github/valinurovam/safequeue) [![Go Report Card](https://goreportcard.com/badge/github.com/valinurovam/safequeue)](https://goreportcard.com/report/github.com/valinurovam/safequeue)

SafeQueue is low-level, in-memory, thread-safe, simple and fast FIFO queue in pure Go.


# Getting Started

## Installing

```shell
$ go get -u github.com/valinurovam/safequeue
```
## API
- Push(item interface{})
- PushHead(item interface{})
- Pop() (item interface{})
- Length() uint64
- Purge()

## Usage

Populate queue
```go
queue := safequeue.NewSafeQueue(SIZE)
queueLength := SIZE * 8
for item := 0; item < queueLength; item++ {
    queue.Push(item)
}
```

Fetch items
```go
item := queue.Pop()
```

# Benchmarks
```shell
goos: darwin
goarch: amd64
pkg: github.com/valinurovam/safequeue
BenchmarkSafeQueue_Push-8       20000000               119 ns/op
BenchmarkSafeQueue_Pop-8        50000000               25.5 ns/op

```

# License

SafeQueue source code is available under the MIT [License](/LICENSE).