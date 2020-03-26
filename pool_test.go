package workerpool

import (
	"fmt"
	"sync"
	"testing"
)

func TestNew(t *testing.T) {
	pool := New(2)
	defer pool.Stop()

	if pool == nil {
		t.Fail()
	}
}

func TestInvalidWorker(t *testing.T) {
	pool := New(2)
	defer pool.Stop()

	if pool.AddWorker("worker") == nil {
		t.Fail()
	}
}

func TestWorkers(t *testing.T) {
	delegateWorkToWorkers(t, 2, 3, 3) // same workers as jobs
	delegateWorkToWorkers(t, 2, 3, 2) // less workers then jobs
	delegateWorkToWorkers(t, 2, 3, 5) // more workers than jobs
}

func delegateWorkToWorkers(t *testing.T, poolSize int, jobsAmount int, workersAmount int) {
	var wg sync.WaitGroup

	pool := New(poolSize)
	out := make(chan int, jobsAmount)
	worker := func(i int) {
		defer wg.Done()
		out <- i
	}

	for i := 1; i <= workersAmount; i++ {
		if err := pool.AddWorker(worker); err != nil {
			t.Fatal(err)
		}
	}

	wg.Add(jobsAmount)

	for i := 1; i <= jobsAmount; i++ {
		pool.Delegate(i)
	}

	go func() {
		wg.Wait()
		close(out)
		pool.Stop()
	}()

	sum := 0
	for n := range out {
		sum += n
	}

	if sum == 0 {
		fmt.Printf("Delegating job %d to %d workers failed", jobsAmount, workersAmount)
		t.Fail()
	}
}
