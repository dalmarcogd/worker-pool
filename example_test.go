package workerpool_test

import (
	"fmt"
	"strconv"
	"sync"

	workerpool "github.com/vardius/worker-pool/v2"
)

func Example() {
	var wg sync.WaitGroup

	poolSize := 1
	jobsAmount := 3
	workersAmount := 2

	// create new pool
	pool := workerpool.New(poolSize)
	out := make(chan int, jobsAmount)
	worker := func(i int) {
		defer wg.Done()
		out <- i
	}

	for i := 1; i <= workersAmount; i++ {
		if err := pool.AddWorker(worker); err != nil {
			panic(err)
		}
	}

	wg.Add(jobsAmount)

	for i := 0; i < jobsAmount; i++ {
		if err := pool.Delegate(i); err != nil {
			panic(err)
		}
	}

	go func() {
		// stop all workers after jobs are done
		wg.Wait()
		close(out)
		pool.Stop()
	}()

	sum := 0
	for n := range out {
		sum += n
	}

	fmt.Println(sum)
	// Output:
	// 3
}

func Example_second() {
	poolSize := 2
	jobsAmount := 8
	workersAmount := 3

	ch := make(chan int, jobsAmount)
	defer close(ch)

	// create new pool
	pool := workerpool.New(poolSize)
	defer pool.Stop()

	worker := func(i int, out chan<- int) { out <- i }

	for i := 1; i <= workersAmount; i++ {
		if err := pool.AddWorker(worker); err != nil {
			panic(err)
		}
	}

	go func() {
		for n := 0; n < jobsAmount; n++ {
			if err := pool.Delegate(n, ch); err != nil {
				panic(err)
			}
		}
	}()

	var sum = 0
	for sum < jobsAmount {
		select {
		case <-ch:
			sum++
		}
	}

	fmt.Println(sum)
	// Output:
	// 8
}

func Example_third() {
	poolSize := 2
	jobsAmount := 8
	workersAmount := 3

	var wg sync.WaitGroup
	wg.Add(jobsAmount)

	// allocate queue
	pool := workerpool.New(poolSize)
	worker := func(s string) {
		defer wg.Done()
		defer fmt.Println("job " + s + " is done !")
		fmt.Println("job " + s + " is running ..")
	}

	// mock arg
	argx := make([]string, jobsAmount)
	for j := 0; j < jobsAmount; j++ {
		argx[j] = "_" + strconv.Itoa(j) + "_"
	}

	// start workers
	for i := 1; i <= workersAmount; i++ {
		if err := pool.AddWorker(worker); err != nil {
			panic(err)
		}
	}

	// assign jobs
	for i := 0; i < jobsAmount; i++ {
		go func(i int) {
			if err := pool.Delegate(argx[i]); err != nil {
				panic(err)
			}
		}(i)
	}

	// clean up
	wg.Wait()
	pool.Stop()

	// Unordered output:
	// job _0_ is running ..
	// job _0_ is done !
	// job _1_ is running ..
	// job _1_ is done !
	// job _2_ is running ..
	// job _2_ is done !
	// job _3_ is running ..
	// job _3_ is done !
	// job _4_ is running ..
	// job _4_ is done !
	// job _5_ is running ..
	// job _5_ is done !
	// job _6_ is running ..
	// job _6_ is done !
	// job _7_ is running ..
	// job _7_ is done !
}
